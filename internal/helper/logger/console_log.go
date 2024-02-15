package logger

import (
	"fmt"
	"os"
	"strings"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func newConsoleLogger() *zap.Logger {
	config := zap.NewProductionEncoderConfig()
	config.EncodeTime = zapcore.TimeEncoderOfLayout(time.RFC822Z)
	consoleEncoder := zapcore.NewConsoleEncoder(config)
	consoleWriter := zapcore.AddSync(os.Stdout)
	consoleCore := zapcore.NewCore(consoleEncoder, consoleWriter, zapcore.DebugLevel)
	core := zapcore.NewTee(consoleCore)
	return zap.New(core, zap.AddCaller(), zap.AddCallerSkip(2))
}

func (zl *ZapLog) SetDebugSymbol(symbol string) *ZapLog {
	zl.DebugSymbol = symbol
	return zl
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
func Errorf(format string, args ...interface{}) {
	Zlog.Errorf(format, args...)
}

func (zl *ZapLog) Debugf(format string, args ...interface{}) {
	Zlog.ConsoleLogger = newConsoleLogger()
	defer Zlog.ConsoleLogger.Sync() // Ensure logs are flushed

	sugar := Zlog.ConsoleLogger.Sugar()
	if zl.Level <= DebugLevel {
		fmt.Printf("%s DEBUG %s\n", strings.Repeat(zl.DebugSymbol, 20), strings.Repeat(zl.DebugSymbol, 20))
		sugar.Debugf(format, args...)
		fmt.Println(strings.Repeat(zl.DebugSymbol, 47))
	}

	zl.SetDebugSymbol("*") // reset to default symbol
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
