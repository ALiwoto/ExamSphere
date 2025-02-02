package logging

import (
	"fmt"
	"io/fs"
	"log"
	"os"
	"runtime/debug"
	"strings"
	"time"

	"ExamSphere/src/core/utils/timeUtils"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func InitZapLog(debug bool) *zap.Logger {
	var config zap.Config
	if debug {
		config = zap.NewDevelopmentConfig()
	} else {
		config = zap.NewProductionConfig()
	}

	config.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	config.EncoderConfig.TimeKey = "timestamp"
	config.EncoderConfig.EncodeTime = zapcore.RFC3339TimeEncoder
	config.EncoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder
	logger, _ := config.Build(zap.AddCallerSkip(1))
	return logger
}

func LoadLogger(debug bool) func() {
	if AppLogger != nil {
		return nil
	}

	loggerMgr := InitZapLog(debug)
	zap.ReplaceGlobals(loggerMgr)
	AppLogger = loggerMgr.Sugar()

	LogStorages = append(LogStorages, GetMemoryLogStorage())

	return func() {
		_ = loggerMgr.Sync()
	}
}

func GetMemoryLogStorage() LogStorage {
	return &memoryLogStorage{
		logs: make([]*LogEntry, 0),
	}
}

func storeLog(theType LogType, args ...any) {
	if len(LogStorages) == 0 {
		return
	}

	logMessage := fmt.Sprint(args...)
	doStoreLog(theType, logMessage)
}

func storeLogf(theType LogType, template string, args ...any) {
	if len(LogStorages) == 0 {
		return
	}

	logMessage := fmt.Sprintf(template, args...)
	doStoreLog(theType, logMessage)
}

func doStoreLog(theType LogType, logMessage string) {
	logDetails := ""
	if len(logMessage) > MaxLogTitleLen {
		logDetails = logMessage
		logMessage = extractLogTitle(logMessage)
	}

	for _, storage := range LogStorages {
		_ = storage.StoreLog(&LogEntry{
			Type:    theType,
			Date:    time.Now().Format(time.RFC3339),
			Message: logMessage,
			Details: logDetails,
		})
	}
}

func extractLogTitle(logMessage string) string {
	if len(logMessage) < MaxLogTitleLen {
		return logMessage
	}

	logMessage = strings.TrimSpace(strings.Split(logMessage, "\n")[0])
	if len(logMessage) < MaxLogTitleLen {
		return logMessage
	}

	if strings.Contains(logMessage, "]:") {
		logMessage = strings.TrimSpace(strings.Split(logMessage, "]:")[1])
	}
	return logMessage
}

func Warn(args ...interface{}) {
	if AppLogger != nil {
		AppLogger.Warn(args...)
	} else {
		log.Println(args...)
	}

	storeLog(LogTypeWarning, args...)
}

func Error(args ...interface{}) {
	if AppLogger != nil {
		AppLogger.Error(args...)
	} else {
		log.Println(args...)
	}

	storeLog(LogTypeError, args...)
}

// UnexpectedError works like Error function and logs the error details to a
// specified log file (a new log file is used each time).
func UnexpectedError(args ...interface{}) {
	err := fmt.Sprint(args...)
	if AppLogger != nil {
		AppLogger.Error("[UNEXPECTED ERROR]: ", err)
	} else {
		log.Println("[UNEXPECTED ERROR]: ", err)
	}

	errorDetails := fmt.Sprintf(
		"[%s]: %v\nStack Trace: %s",
		time.Now().Format(time.RFC3339),
		err,
		debug.Stack(),
	)
	_ = os.WriteFile(GetErrorLogFilePath(), []byte(errorDetails), fs.ModePerm)

	storeLog(LogTypeError, errorDetails)
}

// UnexpectedPanic works like Error function and logs the error details to a
// specified log file (a new log file is used each time).
func UnexpectedPanic(args ...interface{}) {
	err := fmt.Sprint(args...)
	if AppLogger != nil {
		AppLogger.Error("[UNEXPECTED PANIC]: ", err)
	} else {
		log.Println("[UNEXPECTED PANIC]: ", err)
	}

	panicDetails := fmt.Sprintf(
		"[%s]: %v\nStack Trace: %s",
		time.Now().Format(time.RFC3339),
		err,
		debug.Stack(),
	)
	_ = os.WriteFile(GetPanicLogPath(), []byte(panicDetails), fs.ModePerm)

	storeLog(LogTypeError, panicDetails)
}

func Info(args ...interface{}) {
	if AppLogger != nil {
		AppLogger.Info(args...)
	} else {
		log.Println(args...)
	}

	storeLog(LogTypeInfo, args...)
}

func Infof(template string, args ...interface{}) {
	if AppLogger != nil {
		AppLogger.Infof(template, args...)
	} else {
		log.Printf(template, args...)
	}

	storeLogf(LogTypeInfo, template, args...)
}

func Debug(args ...interface{}) {
	if AppLogger != nil {
		AppLogger.Debug(args...)
	} else {
		log.Println(args...)
	}

	storeLog(LogTypeDebug, args...)
}

func Debugf(template string, args ...interface{}) {
	if AppLogger != nil {
		AppLogger.Debugf(template, args...)
	} else {
		log.Printf(template, args...)
	}

	storeLogf(LogTypeDebug, template, args...)
}

func Fatal(args ...interface{}) {
	if AppLogger != nil {
		AppLogger.Fatal(args...)
	} else {
		log.Fatal(args...)
	}

	storeLog(LogTypeError, args...)
}

func GetAllLogEntries(storageName string) ([]*LogEntry, error) {
	for _, storage := range LogStorages {
		if storageName == "" || storage.GetStorageName() == storageName {
			return storage.GetAllLogEntries()
		}
	}

	return nil, ErrLogStorageNotFound
}

func GetErrorLogFilePath() string {
	p := string(os.PathSeparator)
	return "logs" + p + "errors/" +
		"error_" + timeUtils.GenerateSuitableDateTime() + ".log"
}

func GetPanicLogPath() string {
	p := string(os.PathSeparator)
	return "logs" + p + "panics/" +
		"panic_" + timeUtils.GenerateSuitableDateTime() + ".log"
}
