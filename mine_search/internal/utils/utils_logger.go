package utils

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func InitZapLogger(logType string) *zap.Logger {
	writer := getLogWriter(logType)
	encoder := getLogEncoder()
	core := zapcore.NewCore(encoder, writer, zapcore.DebugLevel)
	return zap.New(core, zap.AddCaller())
}
func MyCaller(caller zapcore.EntryCaller, enc zapcore.PrimitiveArrayEncoder) {
	enc.AppendString(filepath.Base(caller.FullPath()))
}
func getLogEncoder() zapcore.Encoder {
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.EncodeTime = zapcore.RFC3339TimeEncoder
	return zapcore.NewJSONEncoder(encoderConfig)
}
func getLogFileWriterSyncer() zapcore.WriteSyncer {
	loggerPath := initLoggerPath()
	layout := "2006-01-02"
	sep := getPathSeparator()
	fileNameSaved := fmt.Sprintf("%s%slogs%sinfo_debugs%s%s",
		loggerPath,
		sep,
		sep,
		sep,
		time.Now().Format(layout)+".log")
	file, err := os.OpenFile(fileNameSaved, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		return nil
	}
	return zapcore.AddSync(file)
}
func getLogConsoleWriterSyncer() zapcore.WriteSyncer {
	return zapcore.AddSync(os.Stdout)
}
func getLogWriter(logType string) zapcore.WriteSyncer {
	switch logType {
	case "file":
		return getLogFileWriterSyncer()
	case "clg":
		return getLogConsoleWriterSyncer()
	case "*":
		return zap.CombineWriteSyncers(getLogFileWriterSyncer(), getLogConsoleWriterSyncer())
	default:
		return nil
	}
}
func initLoggerPath() string {
	GetAbsPath, err := os.Getwd()
	if err != nil {
		return ""
	}
	return GetAbsPath
}
func getPathSeparator() string {
	switch runtime.GOOS {
	case "linux":
		return "/"
	case "windows":
		return "\\"
	}
	return "/"
}
