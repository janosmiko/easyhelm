package logger

import (
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type Logger struct {
	*zap.SugaredLogger
}

func Init(opts ...zap.Option) *Logger {
	logger := initZapLogger(opts...)

	zap.ReplaceGlobals(logger)

	return &Logger{
		logger.Sugar(),
	}
}

func initZapLogger(opts ...zap.Option) (log *zap.Logger) {
	opts = append(opts, zap.AddCallerSkip(1))

	log, err := initZapConfig().Build(opts...)
	if err != nil {
		panic(err)
	}

	return log
}

func initZapConfig() zap.Config {
	var defaultLevel zapcore.Level

	if viper.GetBool("DEBUG") {
		defaultLevel = zap.DebugLevel
	} else {
		defaultLevel = zap.InfoLevel
	}

	var logLevel zap.AtomicLevel

	switch viper.GetString("LOG_LEVEL") {
	case "debug":
		logLevel = zap.NewAtomicLevelAt(zap.DebugLevel)
	case "info":
		logLevel = zap.NewAtomicLevelAt(zap.InfoLevel)
	case "warn", "warning":
		logLevel = zap.NewAtomicLevelAt(zap.WarnLevel)
	case "err", "error":
		logLevel = zap.NewAtomicLevelAt(zap.ErrorLevel)
	case "fatal":
		logLevel = zap.NewAtomicLevelAt(zap.FatalLevel)
	case "panic":
		logLevel = zap.NewAtomicLevelAt(zap.PanicLevel)
	default:
		logLevel = zap.NewAtomicLevelAt(defaultLevel)
	}

	var zc zap.Config
	if viper.GetBool("PRODUCTION_MODE") {
		zc = zap.NewProductionConfig()
		// zc.EncoderConfig.FunctionKey = "F"
	} else {
		zc = zap.NewDevelopmentConfig()
		zc.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
		zc.EncoderConfig.FunctionKey = "" // F
		zc.EncoderConfig.CallerKey = ""   // L
		zc.EncoderConfig.StacktraceKey = "stacktrace"
	}

	if viper.GetBool("LOG_CALLER") {
		zc.EncoderConfig.FunctionKey = "F" // F
		zc.EncoderConfig.CallerKey = "L"   // L
	}

	zc.OutputPaths = []string{"stdout"}
	zc.ErrorOutputPaths = []string{"stdout"}
	zc.Level = logLevel
	zc.EncoderConfig.TimeKey = "timestamp"
	zc.EncoderConfig.EncodeTime = zapcore.RFC3339TimeEncoder

	zc.DisableStacktrace = !viper.GetBool("STACK_TRACE")

	return zc
}
