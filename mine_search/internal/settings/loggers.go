package settings

import (
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func NewLogger() *zap.Logger {
	loggerCore := zapcore.NewCore(logEndcoder(), zapcore.AddSync(os.Stdout), zap.DebugLevel)
	return zap.New(loggerCore, zap.AddCaller())
}
func logEndcoder() zapcore.Encoder {
	encodeConfig := zap.NewProductionEncoderConfig()
	encodeConfig.EncodeTime = zapcore.RFC3339TimeEncoder
	return zapcore.NewJSONEncoder(encodeConfig)
}
