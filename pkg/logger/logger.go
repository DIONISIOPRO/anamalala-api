package logger

import (
	"os"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// Logger encapsula o logger Zap
type Logger struct {
	zap *zap.Logger
}

// NewLogger cria uma nova instância de Logger
func NewLogger(environment string) *Logger {
	// Configuração do encoder
	encoderConfig := zapcore.EncoderConfig{
		TimeKey:        "time",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "caller",
		MessageKey:     "msg",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.LowercaseLevelEncoder,
		EncodeTime:     zapcore.ISO8601TimeEncoder,
		EncodeDuration: zapcore.StringDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}

	// Configurar saída
	var core zapcore.Core
	if environment == "production" {
		// Em produção, usar JSON para facilitar a integração com sistemas de log
		core = zapcore.NewCore(
			zapcore.NewJSONEncoder(encoderConfig),
			zapcore.AddSync(os.Stdout),
			zap.NewAtomicLevelAt(zap.InfoLevel),
		)
	} else {
		// Em desenvolvimento, usar console para legibilidade
		core = zapcore.NewCore(
			zapcore.NewConsoleEncoder(encoderConfig),
			zapcore.AddSync(os.Stdout),
			zap.NewAtomicLevelAt(zap.DebugLevel),
		)
	}

	// Criar logger
	logger := zap.New(core, zap.AddCaller(), zap.AddStacktrace(zap.ErrorLevel))

	return &Logger{
		zap: logger,
	}
}

// Info registra uma mensagem no nível INFO
func (l *Logger) Info(msg string, keysAndValues ...interface{}) {
	l.zap.Sugar().Infow(msg, keysAndValues...)
}

// Debug registra uma mensagem no nível DEBUG
func (l *Logger) Debug(msg string, keysAndValues ...interface{}) {
	l.zap.Sugar().Debugw(msg, keysAndValues...)
}

// Warn registra uma mensagem no nível WARN
func (l *Logger) Warn(msg string, keysAndValues ...interface{}) {
	l.zap.Sugar().Warnw(msg, keysAndValues...)
}

// Error registra uma mensagem no nível ERROR
func (l *Logger) Error(msg string, keysAndValues ...interface{}) {
	l.zap.Sugar().Errorw(msg, keysAndValues...)
}

// Fatal registra uma mensagem no nível FATAL e encerra o programa
func (l *Logger) Fatal(msg string, keysAndValues ...interface{}) {
	l.zap.Sugar().Fatalw(msg, keysAndValues...)
}

// WithField adiciona um campo ao logger
func (l *Logger) WithField(key string, value interface{}) *Logger {
	return &Logger{
		zap: l.zap.With(zap.Any(key, value)),
	}
}

// RequestLogger registra informações de requisição HTTP
func (l *Logger) RequestLogger(method, path, ip, userAgent string, status int, duration time.Duration) {
	l.Info("http_request",
		"method", method,
		"path", path,
		"ip", ip,
		"user_agent", userAgent,
		"status", status,
		"duration_ms", duration.Milliseconds(),
	)
}

// Close fecha o logger
func (l *Logger) Close() error {
	return l.zap.Sync()
}
