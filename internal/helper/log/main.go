package log

import (
	"fmt"
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

type OutputTypes struct {
	Console bool
	File    bool
}

type ZapLog struct {
	Output   OutputTypes
	Filename *string
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
	config.EncodeTime = zapcore.ISO8601TimeEncoder

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
