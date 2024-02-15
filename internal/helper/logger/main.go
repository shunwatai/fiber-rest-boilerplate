package logger

import (
	"fmt"
	"golang-api-starter/internal/config"
	"os"
	"slices"
	"strings"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

type OutputTypes struct {
	Console bool
	File    bool
}

type ZapLog struct {
	Logger      *zap.Logger
	Output      OutputTypes
	Filename    *string
	Level       string
	DebugSymbol string
}

var cfg = config.Cfg
var Zlog = &ZapLog{}

func NewZlog() {
	Zlog.Output = OutputTypes{
		File:    slices.Contains(cfg.Logging.Zap.Output, "file"),
		Console: slices.Contains(cfg.Logging.Zap.Output, "console"),
	}
	Zlog.Filename = &cfg.Logging.Zap.Filename
	Zlog.Level = "debug"
	Zlog.DebugSymbol = "*"
}

func (zl *ZapLog) GetField(key string, value interface{}, fieldType *string) zap.Field {
	var zapField zap.Field
	switch *fieldType {
	case "int64":
		zapField = zap.Int64(key, value.(int64))
	case "duration":
		zapField = zap.Duration(key, value.(time.Duration))
	case "date":
		zapField = zap.Time(key, value.(time.Time))
	case "bool":
		zapField = zap.Bool(key, value.(bool))
	default:
		zapField = zap.String(key, value.(string))
	}
	return zapField
}

func (zl *ZapLog) SetLevel(lvl string) *ZapLog {
	zl.Level = lvl
	return zl
}

func (zl *ZapLog) SetDebugSymbol(symbol string) *ZapLog {
	zl.DebugSymbol = symbol
	return zl
}

func Printf(format string, args ...interface{}) {
	Zlog.Printf(format, args)
}
func (zl *ZapLog) Printf(format string, args ...interface{}) {
	config := zap.NewProductionEncoderConfig()
	config.EncodeTime = zapcore.TimeEncoderOfLayout(time.RFC822Z)
	consoleEncoder := zapcore.NewConsoleEncoder(config)
	consoleWriter := zapcore.AddSync(os.Stdout)
	consoleCore := zapcore.NewCore(consoleEncoder, consoleWriter, zapcore.DebugLevel)
	core := zapcore.NewTee(consoleCore)
	logger := zap.New(core, zap.AddCaller(), zap.AddStacktrace(zapcore.ErrorLevel))

	defer logger.Sync() // Ensure logs are flushed

	sugar := logger.WithOptions(zap.AddCallerSkip(2)).Sugar()
	switch zl.Level {
	case "info":
		sugar.Infof(format, args...)
	case "warn":
		sugar.Warnf(format, args...)
	case "error":
		sugar.Errorf(format, args...)
	default: // debug
		fmt.Printf("%s DEBUG %s\n", strings.Repeat(zl.DebugSymbol, 20), strings.Repeat(zl.DebugSymbol, 20))
		sugar.Debugf(format, args...)
		fmt.Println(strings.Repeat(zl.DebugSymbol, 47))
	}

	zl.SetLevel("debug")   // reset to default level
	zl.SetDebugSymbol("*") // reset to default symbol
}

func (zl *ZapLog) RequestLog(msg string, keysAndValues ...interface{}) {
	filename := "requests.log"
	if zl.Output.File && zl.Filename != nil {
		filename = *zl.Filename
	}

	logger, err := fileLogger(filename, zl.Output)
	if err != nil {
		fmt.Printf("Failed to initialize logger: %v\n", err)
		return
	}
	defer logger.Sync() // Ensure logs are flushed

	sugar := logger.Sugar()
	sugar.Infow(msg, keysAndValues...)
}

// fileLogger initializes a zap.Logger that writes to both the console and a specified file.
// ref: https://www.golinuxcloud.com/golang-zap-logger/#Setting_Output_in_Zap_Console_Log_File_or_Both
func fileLogger(filename string, outputTypes OutputTypes) (*zap.Logger, error) {
	// Configure the time format
	config := zap.NewProductionEncoderConfig()
	config.EncodeTime = zapcore.RFC3339TimeEncoder

	// Create file and console encoders
	fileEncoder := zapcore.NewJSONEncoder(config)
	consoleEncoder := zapcore.NewConsoleEncoder(config)

	// Open the log file
	// Set up lumberjack as a logger:
	logFile := &lumberjack.Logger{
		Filename:   fmt.Sprintf("./log/%s", filename), // Or any other path
		MaxSize:    500,                               // MB; after this size, a new log file is created
		MaxBackups: 10,                                // Number of backups to keep
		MaxAge:     28,                                // Days
		Compress:   true,                              // Compress the backups using gzip
	}

	// Default logFile without log rotate
	// logFilePath := fmt.Sprintf("./log/%s", filename)
	// logFile, err := os.OpenFile(logFilePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	// if err != nil {
	// 	return nil, fmt.Errorf("failed to open log file: %v", err)
	// }

	// Create writers for file and console
	fileWriter := zapcore.AddSync(logFile)
	consoleWriter := zapcore.AddSync(os.Stdout)

	// Set the log level
	defaultLogLevel := zapcore.DebugLevel

	// Create cores for writing to the file and console
	fileCore := zapcore.NewCore(fileEncoder, fileWriter, defaultLogLevel)
	consoleCore := zapcore.NewCore(consoleEncoder, consoleWriter, defaultLogLevel)

	// Combine cores
	var core zapcore.Core
	if outputTypes.Console && outputTypes.File {
		core = zapcore.NewTee(fileCore, consoleCore)
	} else if outputTypes.Console && !outputTypes.File {
		core = zapcore.NewTee(consoleCore)
	} else if !outputTypes.Console && outputTypes.File {
		core = zapcore.NewTee(fileCore)
	}

	// Create the logger with additional context information (caller, stack trace)
	logger := zap.New(core, zap.AddCaller(), zap.AddStacktrace(zapcore.ErrorLevel))

	return logger, nil
}
