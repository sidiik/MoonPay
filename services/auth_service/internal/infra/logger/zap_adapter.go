package logger

import (
	"github.com/sidiik/moonpay/auth_service/internal/domain"
	"go.uber.org/zap"
)

type ZapLogger struct {
	z *zap.Logger
}

func NewZapLogger() domain.Logger {
	return &ZapLogger{z: L()}
}

func (l *ZapLogger) Info(msg string, fields ...any) {
	l.z.Sugar().Infow(msg, fields...)
}

func (l *ZapLogger) Error(msg string, fields ...any) {
	l.z.Sugar().Errorw(msg, fields...)
}

func (l *ZapLogger) Debug(msg string, fields ...any) {
	l.z.Sugar().Debugw(msg, fields...)
}

func (l *ZapLogger) Warn(msg string, fields ...any) {
	l.z.Sugar().Warnw(msg, fields...)
}
