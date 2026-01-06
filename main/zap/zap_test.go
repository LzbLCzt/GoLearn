package zap

import (
	"os"
	"testing"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

// 初始化日志配置
func initLogger() (*zap.SugaredLogger, error) {
	// 1. 配置日志编码器
	encoderConfig := zapcore.EncoderConfig{
		TimeKey:       "ts",
		LevelKey:      "level",
		CallerKey:     "caller",
		MessageKey:    "msg",
		StacktraceKey: "stacktrace",
		LineEnding:    zapcore.DefaultLineEnding,
		EncodeLevel:   zapcore.LowercaseLevelEncoder, // 小写日志级别
		EncodeTime:    zapcore.ISO8601TimeEncoder,    // ISO8601 时间格式
		EncodeCaller:  zapcore.ShortCallerEncoder,    // 短路径文件格式
	}

	// 2. 配置日志文件切割（使用lumberjack）
	logFile := &lumberjack.Logger{
		Filename:   "./logs/app.log", // 日志文件路径
		MaxSize:    10,               // 单个文件最大大小(MB)
		MaxBackups: 3,                // 保留旧文件的最大个数
		MaxAge:     30,               // 保留旧文件的最大天数
		Compress:   false,            // 是否压缩/归档旧文件
	}

	// 3. 创建日志核心（同时输出到控制台和文件）
	core := zapcore.NewCore(
		zapcore.NewJSONEncoder(encoderConfig), // JSON编码器
		zapcore.NewMultiWriteSyncer( // 多路输出
			zapcore.AddSync(logFile),   // 输出到文件
			zapcore.AddSync(os.Stdout), // 输出到控制台
		),
		zapcore.DebugLevel, // 日志级别（Debug及以上级别都会记录）
	)

	// 4. 创建Logger实例
	logger := zap.New(core, zap.AddCaller())
	return logger.Sugar(), nil
}

func TestZap(t *testing.T) {
	// 初始化日志
	logger, err := initLogger()
	if err != nil {
		panic(err)
	}
	defer logger.Sync() // 确保程序退出前刷新缓存中的日志

	// 记录不同级别的日志
	logger.Debug("这是一条Debug级别的日志", "count", 42)
	logger.Info("服务启动成功", "port", 8080)
	logger.Warn("数据库连接较慢", "duration", 150)
	logger.Error("文件读取失败", "file", "config.yaml", "error", "文件不存在")

	// 使用格式化输出（性能略低，但更灵活）
	logger.Infof("用户 %s 登录成功，ID: %d", "张三", 1001)

	// 结构化日志（推荐方式）
	logger.Infow("请求处理完成",
		"method", "GET",
		"path", "/api/users",
		"status", 200,
		"duration", 45,
	)
}
