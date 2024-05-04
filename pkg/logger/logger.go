package logger

import (
	"io"
	"log"
	"os"
	"path/filepath"

	"github.com/jis4nx/go-ecom/helpers/types"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// Logger provides logging functionality with both console and file outputs.
type Logger struct {
	*zap.Logger
	File io.Writer
}

func (l *Logger) InitLogger(s types.ServiceInfo) {
	rootLogger := newLogger(os.Stdout, l.File)
	childLogger := rootLogger.With(
		zap.String("service", s.Name),
		zap.String("host", s.Host),
		zap.String("port", s.Port),
	)
	l.Logger = childLogger
}

func (l *Logger) SetLogFile(baseLogDir, file string) {
	logFile := filepath.Join(baseLogDir, "logs", file)
	if err := os.MkdirAll(filepath.Dir(logFile), 0755); err != nil {
		log.Fatal(err)
	}

	f, err := os.OpenFile(logFile, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0600)
	if err != nil {
		log.Fatal(err)
	}
	l.File = f
}

func newLogger(stdout, file io.Writer) *zap.Logger {
	// Define log levels
	devLevel := zap.NewAtomicLevelAt(zap.DebugLevel)
	prodLevel := zap.NewAtomicLevelAt(zap.InfoLevel)

	// Configure encoder for production mode
	prodCfg := zap.NewProductionEncoderConfig()
	prodCfg.TimeKey = "timestamp"
	prodCfg.EncodeTime = zapcore.ISO8601TimeEncoder

	// Configure encoder for development mode
	devCfg := zap.NewProductionEncoderConfig()
	devCfg.EncodeLevel = zapcore.CapitalColorLevelEncoder

	consoleEncoder := zapcore.NewConsoleEncoder(devCfg)
	fileEncoder := zapcore.NewJSONEncoder(prodCfg)

	// Create a core that writes logs to both console and file
	core := zapcore.NewTee(
		zapcore.NewCore(consoleEncoder, zapcore.AddSync(stdout), devLevel),
		zapcore.NewCore(fileEncoder, zapcore.AddSync(file), prodLevel),
	)

	return zap.New(core)
}
