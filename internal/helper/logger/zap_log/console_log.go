package logger

import (
	"fmt"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"os"
	"strings"
)

func newConsoleLogger() *zap.Logger {
	config := zap.NewProductionEncoderConfig()
	// config.EncodeTime = zapcore.TimeEncoderOfLayout(time.RFC822Z)
	config.TimeKey = zapcore.OmitKey // skip time field for dev
	consoleEncoder := zapcore.NewConsoleEncoder(config)
	consoleWriter := zapcore.AddSync(os.Stdout)
	consoleCore := zapcore.NewCore(consoleEncoder, consoleWriter, zapcore.DebugLevel)
	core := zapcore.NewTee(consoleCore)
	return zap.New(core, zap.AddCaller(), zap.AddCallerSkip(2))
}

func Debugf(format string, args ...interface{}) {
	Zlog.Debugf(format, args...)
}
func Infof(format string, args ...interface{}) {
	Zlog.Infof(format, args...)
}
func Warnf(format string, args ...interface{}) {
	Zlog.Warnf(format, args...)
}
func Errorf(format string, args ...interface{}) error {
	Zlog.Errorf(format, args...)
	return fmt.Errorf(format, args...)
}
func Fatalf(format string, args ...interface{}) {
	Zlog.Fatalf(format, args...)
}

func (zl *ZapLog) Debugf(format string, args ...interface{}) {
	Zlog.ConsoleLogger = newConsoleLogger()
	defer Zlog.ConsoleLogger.Sync() // Ensure logs are flushed

	sugar := Zlog.ConsoleLogger.Sugar()

	if zl.Level <= DebugLevel {
		if zl.DebugSymbol != nil {
			fmt.Printf("%s DEBUG %s\n", strings.Repeat(*zl.DebugSymbol, 20), strings.Repeat(*zl.DebugSymbol, 20))
			sugar.Debugf(format, args...)
			fmt.Println(strings.Repeat(*zl.DebugSymbol, 47))
		} else {
			sugar.Debugf(format, args...)
		}
	}

	zl.SetDebugSymbol(cfg.Logging.DebugSymbol) // reset to default symbol
}

func (zl *ZapLog) Infof(format string, args ...interface{}) {
	Zlog.ConsoleLogger = newConsoleLogger()
	defer Zlog.ConsoleLogger.Sync() // Ensure logs are flushed

	sugar := Zlog.ConsoleLogger.Sugar()
	if zl.Level <= InfoLevel {
		sugar.Infof(format, args...)
	}
}

func (zl *ZapLog) Warnf(format string, args ...interface{}) {
	Zlog.ConsoleLogger = newConsoleLogger()
	defer Zlog.ConsoleLogger.Sync() // Ensure logs are flushed

	sugar := Zlog.ConsoleLogger.Sugar()
	if zl.Level <= WarningLevel {
		sugar.Warnf(format, args...)
	}
}

func (zl *ZapLog) Errorf(format string, args ...interface{}) {
	Zlog.ConsoleLogger = newConsoleLogger()
	defer Zlog.ConsoleLogger.Sync() // Ensure logs are flushed

	sugar := Zlog.ConsoleLogger.Sugar()
	if zl.Level <= ErrorLevel {
		sugar.Errorf(format, args...)
	}
}

func (zl *ZapLog) Fatalf(format string, args ...interface{}) {
	Zlog.ConsoleLogger = newConsoleLogger()
	defer Zlog.ConsoleLogger.Sync() // Ensure logs are flushed

	sugar := Zlog.ConsoleLogger.Sugar()
	if zl.Level <= ErrorLevel {
		sugar.Fatalf(format, args...)
	}
}
