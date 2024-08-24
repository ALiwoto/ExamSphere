package logging

import (
	"fmt"
	"io/fs"
	"log"
	"os"
	"runtime/debug"
	"time"

	"ExamSphere/src/core/utils/timeUtils"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var AppLogger *zap.SugaredLogger

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

	return func() {
		_ = loggerMgr.Sync()
	}
}

func Warn(args ...interface{}) {
	if AppLogger != nil {
		AppLogger.Warn(args...)
	} else {
		log.Println(args...)
	}
}

func Error(args ...interface{}) {
	if AppLogger != nil {
		AppLogger.Error(args...)
	} else {
		log.Println(args...)
	}
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
}

func Info(args ...interface{}) {
	if AppLogger != nil {
		AppLogger.Info(args...)
	} else {
		log.Println(args...)
	}
}

func Infof(template string, args ...interface{}) {
	if AppLogger != nil {
		AppLogger.Infof(template, args...)
	} else {
		log.Printf(template, args...)
	}
}

func Debug(args ...interface{}) {
	if AppLogger != nil {
		AppLogger.Debug(args...)
	} else {
		log.Println(args...)
	}
}

func Debugf(template string, args ...interface{}) {
	if AppLogger != nil {
		AppLogger.Debugf(template, args...)
	} else {
		log.Printf(template, args...)
	}
}

func Fatal(args ...interface{}) {
	if AppLogger != nil {
		AppLogger.Fatal(args...)
	} else {
		log.Fatal(args...)
	}
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
