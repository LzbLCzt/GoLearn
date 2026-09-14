package main

import (
	"context"
	"database/sql"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

// 各类超时参数统一放宽，避免大 SQL / 大事务被中断
const (
	dialTimeout     = 60 * time.Second   // 建立 TCP 连接超时
	readTimeout     = 1800 * time.Second // 读超时 30 分钟
	writeTimeout    = 1800 * time.Second // 写超时 30 分钟
	pingTimeout     = 60 * time.Second
	precheckTimeout = 600 * time.Second  // 执行前 COUNT 校验 10 分钟
	execTimeout     = 1800 * time.Second // 单条 UPDATE 30 分钟
	connMaxLifetime = 30 * time.Minute
	connMaxIdleTime = 10 * time.Minute
	maxAllowedPkt   = 64 << 20 // 64MB，保证大 SQL 完整发送
)

type config struct {
	host     string
	port     int
	user     string
	password string
	db       string
	sqlFile  string
	dryRun   bool
	logFile  string
}

var whereRegexp = regexp.MustCompile(`(?is)^update\s+(\S+)\s+set\s+(.*?)\s+where\s+(.*)$`)

func main() {
	cfg := parseFlags()

	logWriter, closeLog, err := buildLogger(cfg.logFile)
	if err != nil {
		fmt.Fprintf(os.Stderr, "初始化日志失败: %v\n", err)
		os.Exit(1)
	}
	defer closeLog()
	logf := log.New(logWriter, "", log.LstdFlags|log.Lmicroseconds).Printf

	logf("========== 任务开始 ==========")
	logf("数据库地址   : %s:%d", cfg.host, cfg.port)
	logf("数据库用户   : %s", cfg.user)
	logf("默认库       : %q", cfg.db)
	logf("SQL 文件     : %s", cfg.sqlFile)
	logf("执行模式     : %s", modeName(cfg.dryRun))
	logf("超时参数     : dial=%s read=%s write=%s exec=%s precheck=%s",
		dialTimeout, readTimeout, writeTimeout, execTimeout, precheckTimeout)

	raw, err := os.ReadFile(cfg.sqlFile)
	if err != nil {
		logf("读取 SQL 文件失败: %v", err)
		os.Exit(1)
	}
	logf("SQL 文件大小 : %d 字节 (%.2f KB)", len(raw), float64(len(raw))/1024)

	statements := splitStatements(string(raw))
	logf("解析 SQL 语句: %d 条", len(statements))
	for i, s := range statements {
		logf("  语句 #%d: 长度=%d 字节, 预览=%s", i+1, len(s), preview(s, 160))
	}
	if len(statements) == 0 {
		logf("没有可执行的 SQL 语句，退出")
		os.Exit(1)
	}

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=true&loc=Local"+
		"&collation=utf8mb4_general_ci&timeout=%s&readTimeout=%s&writeTimeout=%s&maxAllowedPacket=%d",
		cfg.user, cfg.password, cfg.host, cfg.port, cfg.db,
		dialTimeout, readTimeout, writeTimeout, maxAllowedPkt)
	logf("DSN(脱敏)    : %s:%s@tcp(%s:%d)/%s", cfg.user, "******", cfg.host, cfg.port, cfg.db)

	db, err := sql.Open("mysql", dsn)
	if err != nil {
		logf("创建连接池失败: %v", err)
		os.Exit(1)
	}
	defer db.Close()

	db.SetMaxOpenConns(4)
	db.SetMaxIdleConns(2)
	db.SetConnMaxLifetime(connMaxLifetime)
	db.SetConnMaxIdleTime(connMaxIdleTime)
	logf("连接池配置   : maxOpen=4 maxIdle=2 maxLifetime=%s maxIdleTime=%s",
		connMaxLifetime, connMaxIdleTime)

	if err := pingWithRetry(db, logf); err != nil {
		logf("数据库连接失败: %v", err)
		os.Exit(1)
	}
	logf("数据库连接成功")

	tuneSession(db, logf)
	logServerInfo(db, logf)

	start := time.Now()
	totalAffected := int64(0)
	success, failed := 0, 0

	for i, stmt := range statements {
		idx := i + 1
		logf("---------- 执行语句 #%d/%d ----------", idx, len(statements))
		logf("语句预览: %s", preview(stmt, 200))

		precheck(db, logf, stmt)

		stmtStart := time.Now()
		execCtx, cancel := context.WithTimeout(context.Background(), execTimeout)
		var res sql.Result
		if cfg.dryRun {
			res, err = execInRolledBackTx(execCtx, db, stmt)
		} else {
			res, err = db.ExecContext(execCtx, stmt)
		}
		cancel()

		elapsed := time.Since(stmtStart)
		if err != nil {
			failed++
			logf("执行失败: 耗时=%s, error=%v", elapsed, err)
			continue
		}
		success++

		affected, affErr := res.RowsAffected()
		if affErr != nil {
			logf("获取影响行数失败: %v", affErr)
		} else {
			totalAffected += affected
			logf("执行成功: 耗时=%s, 影响行数=%d", elapsed, affected)
		}
		if lastID, e := res.LastInsertId(); e == nil {
			logf("LastInsertId = %d (UPDATE 场景通常为 0)", lastID)
		}
		logWarnings(db, logf)
	}

	logf("========== 执行汇总 ==========")
	logf("总语句数     : %d", len(statements))
	logf("执行成功     : %d", success)
	logf("执行失败     : %d", failed)
	logf("累计影响行数 : %d", totalAffected)
	logf("总耗时       : %s", time.Since(start))
	logf("执行模式     : %s", modeName(cfg.dryRun))
	if cfg.dryRun {
		logf("注意: DRY-RUN 模式事务已回滚，数据库未发生任何变更")
	}
	logPoolStats(db, logf)
	logf("========== 任务结束 ==========")

	if failed > 0 {
		os.Exit(2)
	}
}

func parseFlags() config {
	var cfg config
	flag.StringVar(&cfg.host, "host", "21.233.61.181", "数据库地址")
	flag.IntVar(&cfg.port, "port", 3306, "数据库端口")
	flag.StringVar(&cfg.user, "user", "router_1yjcjcvf", "数据库用户名")
	flag.StringVar(&cfg.password, "password", "My8JJ5X7ywMnmIcv", "数据库密码")
	flag.StringVar(&cfg.db, "db", "polaris_router", "默认库名，SQL 中表未带库前缀时生效")
	flag.StringVar(&cfg.sqlFile, "sql", "", "待执行的 SQL 文件路径")
	flag.BoolVar(&cfg.dryRun, "dry-run", true, "true 时在事务中执行后回滚，不产生真实变更")
	flag.StringVar(&cfg.logFile, "log", "", "日志文件路径，默认在 SQL 文件同目录按时间戳生成")
	flag.Parse()

	if cfg.sqlFile == "" {
		if flag.NArg() > 0 {
			cfg.sqlFile = flag.Arg(0)
		} else {
			cfg.sqlFile = "update_service_route_rule_full.sql"
		}
	}
	if !filepath.IsAbs(cfg.sqlFile) {
		if abs, err := filepath.Abs(cfg.sqlFile); err == nil {
			cfg.sqlFile = abs
		}
	}
	if cfg.logFile == "" {
		cfg.logFile = filepath.Join(filepath.Dir(cfg.sqlFile),
			fmt.Sprintf("run_%s.log", time.Now().Format("20060102_150405")))
	}
	return cfg
}

func buildLogger(path string) (io.Writer, func(), error) {
	f, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0o644)
	if err != nil {
		return nil, func() {}, err
	}
	return io.MultiWriter(os.Stdout, f), func() { _ = f.Sync(); _ = f.Close() }, nil
}

func pingWithRetry(db *sql.DB, logf func(string, ...interface{})) error {
	var lastErr error
	for i := 1; i <= 3; i++ {
		ctx, cancel := context.WithTimeout(context.Background(), pingTimeout)
		start := time.Now()
		lastErr = db.PingContext(ctx)
		cancel()
		if lastErr == nil {
			logf("Ping 成功: 第 %d 次尝试, 耗时=%s", i, time.Since(start))
			return nil
		}
		logf("Ping 失败: 第 %d 次尝试, 耗时=%s, error=%v", i, time.Since(start), lastErr)
		time.Sleep(time.Duration(i) * 2 * time.Second)
	}
	return lastErr
}

// tuneSession 放宽会话级超时与锁等待，避免大批量更新被服务端中断
func tuneSession(db *sql.DB, logf func(string, ...interface{})) {
	sessions := []string{
		"SET SESSION innodb_lock_wait_timeout = 600",
		"SET SESSION lock_wait_timeout = 600",
		"SET SESSION net_read_timeout = 1800",
		"SET SESSION net_write_timeout = 1800",
	}
	for _, s := range sessions {
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		_, err := db.ExecContext(ctx, s)
		cancel()
		if err != nil {
			logf("会话设置失败: %s, error=%v", s, err)
			continue
		}
		logf("会话设置成功: %s", s)
	}
}

func logServerInfo(db *sql.DB, logf func(string, ...interface{})) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	var v string
	if err := db.QueryRowContext(ctx, "SELECT VERSION()").Scan(&v); err == nil {
		logf("MySQL 版本   : %s", v)
	} else {
		logf("获取 MySQL 版本失败: %v", err)
	}
	if err := db.QueryRowContext(ctx, "SELECT DATABASE()").Scan(&v); err == nil {
		logf("当前库       : %s", v)
	}
	if err := db.QueryRowContext(ctx, "SELECT CURRENT_USER()").Scan(&v); err == nil {
		logf("会话用户     : %s", v)
	}

	rows, err := db.QueryContext(ctx, "SHOW VARIABLES WHERE Variable_name IN "+
		"('max_allowed_packet','innodb_lock_wait_timeout','net_read_timeout','net_write_timeout','wait_timeout')")
	if err == nil {
		defer rows.Close()
		for rows.Next() {
			var name, val string
			if err := rows.Scan(&name, &val); err == nil {
				logf("服务端变量   : %s = %s", name, val)
			}
		}
	}
}

// precheck 用相同 WHERE 条件先做一次 COUNT，便于核对预期影响范围
func precheck(db *sql.DB, logf func(string, ...interface{}), stmt string) {
	m := whereRegexp.FindStringSubmatch(stmt)
	if len(m) != 4 {
		logf("跳过执行前校验: 未能解析出表名与 WHERE 条件")
		return
	}
	table := m[1]
	cond := strings.TrimSpace(m[3])
	countSQL := fmt.Sprintf("SELECT COUNT(*) FROM %s WHERE %s", table, cond)
	logf("执行前校验   : SELECT COUNT(*) FROM %s WHERE ... (条件同 UPDATE)", table)

	ctx, cancel := context.WithTimeout(context.Background(), precheckTimeout)
	defer cancel()
	start := time.Now()
	var cnt int64
	if err := db.QueryRowContext(ctx, countSQL).Scan(&cnt); err != nil {
		logf("执行前校验失败: 耗时=%s, error=%v", time.Since(start), err)
		return
	}
	logf("执行前匹配行数: %d (耗时=%s)", cnt, time.Since(start))
}

func execInRolledBackTx(ctx context.Context, db *sql.DB, stmt string) (sql.Result, error) {
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	res, err := tx.ExecContext(ctx, stmt)
	if err != nil {
		_ = tx.Rollback()
		return nil, err
	}
	if err := tx.Rollback(); err != nil {
		return res, fmt.Errorf("回滚失败: %w", err)
	}
	return res, nil
}

func logWarnings(db *sql.DB, logf func(string, ...interface{})) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	rows, err := db.QueryContext(ctx, "SHOW WARNINGS LIMIT 20")
	if err != nil {
		return
	}
	defer rows.Close()

	cols, _ := rows.Columns()
	n := 0
	for rows.Next() {
		vals := make([]sql.NullString, len(cols))
		ptrs := make([]interface{}, len(cols))
		for i := range vals {
			ptrs[i] = &vals[i]
		}
		if err := rows.Scan(ptrs...); err != nil {
			logf("读取 warning 失败: %v", err)
			return
		}
		parts := make([]string, 0, len(cols))
		for i, v := range vals {
			parts = append(parts, fmt.Sprintf("%s=%s", cols[i], v.String))
		}
		logf("WARNING: %s", strings.Join(parts, " | "))
		n++
	}
	if n == 0 {
		logf("无 warning")
	}
}

func logPoolStats(db *sql.DB, logf func(string, ...interface{})) {
	s := db.Stats()
	logf("连接池状态   : open=%d inUse=%d idle=%d waitCount=%d waitDuration=%s",
		s.OpenConnections, s.InUse, s.Idle, s.WaitCount, s.WaitDuration)
}

// splitStatements 按分号切分 SQL，跳过字符串字面量与注释中的分号
func splitStatements(input string) []string {
	var out []string
	var buf strings.Builder
	runes := []rune(input)
	i := 0
	for i < len(runes) {
		c := runes[i]
		switch {
		case c == '\'' || c == '"' || c == '`':
			quote := c
			buf.WriteRune(c)
			i++
			for i < len(runes) {
				if runes[i] == '\\' && i+1 < len(runes) {
					buf.WriteRune(runes[i])
					buf.WriteRune(runes[i+1])
					i += 2
					continue
				}
				buf.WriteRune(runes[i])
				if runes[i] == quote {
					if i+1 < len(runes) && runes[i+1] == quote {
						buf.WriteRune(runes[i+1])
						i += 2
						continue
					}
					i++
					break
				}
				i++
			}
		case c == '-' && i+2 < len(runes) && runes[i+1] == '-' &&
			(runes[i+2] == ' ' || runes[i+2] == '\n' || runes[i+2] == '\t'):
			for i < len(runes) && runes[i] != '\n' {
				i++
			}
		case c == '#':
			for i < len(runes) && runes[i] != '\n' {
				i++
			}
		case c == '/' && i+1 < len(runes) && runes[i+1] == '*':
			i += 2
			for i+1 < len(runes) && !(runes[i] == '*' && runes[i+1] == '/') {
				i++
			}
			i += 2
		case c == ';':
			if s := strings.TrimSpace(buf.String()); s != "" {
				out = append(out, s)
			}
			buf.Reset()
			i++
		default:
			buf.WriteRune(c)
			i++
		}
	}
	if s := strings.TrimSpace(buf.String()); s != "" {
		out = append(out, s)
	}
	return out
}

func preview(s string, n int) string {
	oneLine := strings.Join(strings.Fields(s), " ")
	runes := []rune(oneLine)
	if len(runes) <= n {
		return oneLine
	}
	return string(runes[:n]) + " ..."
}

func modeName(dryRun bool) string {
	if dryRun {
		return "DRY-RUN(事务内执行后回滚，不落库)"
	}
	return "REAL(真实写入)"
}
