package zap

import (
	"testing"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

func TestZap2(t *testing.T) {
	// 1. 配置编码器（输出格式）
	encoderConfig := zapcore.EncoderConfig{
		TimeKey:       "ts",
		LevelKey:      "level",
		CallerKey:     "caller",
		MessageKey:    "msg",
		StacktraceKey: "stacktrace",
		EncodeLevel:   zapcore.LowercaseLevelEncoder,
		EncodeTime:    zapcore.ISO8601TimeEncoder, // 可读的时间格式
		EncodeCaller:  zapcore.ShortCallerEncoder,
	}
	encoder := zapcore.NewJSONEncoder(encoderConfig) // 输出为JSON格式

	// 2. 配置日志输出目标（使用Lumberjack实现文件轮转）
	logWriter := &lumberjack.Logger{
		Filename:   "./logs/app.log", // 日志文件路径
		MaxSize:    100,              // 单个文件最大100MB
		MaxBackups: 5,                // 保留最多5个旧文件
		MaxAge:     28,               // 保留28天
		Compress:   true,             // 压缩旧文件
	}

	// 3. 设置日志级别
	atomicLevel := zap.NewAtomicLevelAt(zap.InfoLevel)

	// 4. 创建日志核心并构建Logger
	core := zapcore.NewCore(encoder, zapcore.AddSync(logWriter), atomicLevel)
	logger := zap.New(core, zap.AddCaller())
	defer logger.Sync() // 程序退出前刷新缓冲区的日志

	// 使用示例
	logger.Info("服务启动成功",
		zap.String("version", "1.0.0"),
		zap.Int("port", 8080),
	)
}
