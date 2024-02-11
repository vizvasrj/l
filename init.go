package l

import (
	"fmt"
	"log"
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

type StorageType int

const (
	Print StorageType = 1 << iota
	FileStorage
	KafkaStorage
)

type FileStorageType struct {
	Filename   string // log file name
	MaxSize    int    // megabytes
	MaxBackups int
	MaxAge     int // days
}

type KafkaStorageType struct {
	Brokers []string
	Topic   string
}

type Logger struct {
	KafkaStorageType
	FileStorageType
}

func (l *Logger) LoggerInit(s StorageType) *zap.Logger {
	tee := []zapcore.Core{l.StdOutZapCore()}

	// if s&Print != 0 {
	// 	fmt.Println("Print")

	// }
	if s&FileStorage != 0 {
		log.Println("FileStorage")
		tee = append(tee, l.FileStorageZapCore(l.FileStorageType))
	}
	if s&KafkaStorage != 0 {
		log.Println("KafkaStorage")
		tee = append(tee, l.KafkaZapCore())

	}
	core := zapcore.NewTee(tee...)
	return zap.New(
		core,
		zap.AddCaller(),
		zap.AddCallerSkip(1),
	)
}

func (l *Logger) StdOutZapCore() zapcore.Core {
	stdout := zapcore.AddSync(os.Stdout)
	level := zap.NewAtomicLevelAt(zap.InfoLevel)

	developmentCfg := zap.NewDevelopmentEncoderConfig()
	developmentCfg.EncodeLevel = zapcore.CapitalColorLevelEncoder

	consoleEncoder := zapcore.NewConsoleEncoder(developmentCfg)
	return zapcore.NewCore(consoleEncoder, stdout, level)

}

func (l *Logger) validateFileStorageType() {
	if l.Filename == "" {
		l.Filename = "logs/app.log"
		log.Println("Filename is not set, using default value logs/app.log")
	}
	if l.MaxSize == 0 {
		l.MaxSize = 10
		log.Panicln("MaxSize is not set, using default value 10 mb")
	}
	if l.MaxBackups == 0 {
		l.MaxBackups = 3
		log.Println("MaxBackups is not set, using default value 3")
	}
	if l.MaxAge == 0 {
		l.MaxAge = 7
		log.Println("MaxAge is not set, using default value 7 days")
	}
}

func (l *Logger) FileStorageZapCore(k FileStorageType) zapcore.Core {

	l.validateFileStorageType()
	file := zapcore.AddSync(&lumberjack.Logger{
		Filename:   k.Filename,
		MaxSize:    k.MaxSize, // megabytes
		MaxBackups: k.MaxBackups,
		MaxAge:     k.MaxAge, // days
	})

	level := zap.NewAtomicLevelAt(zap.InfoLevel)

	productionCfg := zap.NewProductionEncoderConfig()
	productionCfg.TimeKey = "timestamp"
	productionCfg.EncodeTime = zapcore.ISO8601TimeEncoder

	developmentCfg := zap.NewDevelopmentEncoderConfig()
	developmentCfg.EncodeLevel = zapcore.CapitalColorLevelEncoder

	fileEncoder := zapcore.NewJSONEncoder(productionCfg)

	return zapcore.NewCore(fileEncoder, file, level)
}

func (l *Logger) validateKafkaStorageType() error {
	if len(l.Brokers) == 0 {
		log.Panicln("Brokers is not set")
		return fmt.Errorf("brokers is not set")
	}
	if l.Topic == "" {
		log.Panicln("Topic is not set")
		return fmt.Errorf("topic is not set")
	}
	return nil
}

func (l *Logger) KafkaZapCore() zapcore.Core {
	err := l.validateKafkaStorageType()
	if err != nil {
		panic(err)
	}

	kafkaSyncer, err := NewKafkaSyncer(l.Brokers, l.Topic)
	if err != nil {
		panic(err)
	}

	return zapcore.NewCore(
		zapcore.NewJSONEncoder(
			zap.NewProductionEncoderConfig(),
		),
		kafkaSyncer,
		zapcore.DebugLevel,
	)
}

func (l *Logger) Info(msg string, fields ...zap.Field) {
	logger.Info(msg, fields...)
}

func (l *Logger) Warn(msg string, fields ...zap.Field) {
	logger.Warn(msg, fields...)
}

func (l *Logger) Error(msg string, fields ...zap.Field) {
	logger.Error(msg, fields...)
}

func (l *Logger) Fatal(msg string, fields ...zap.Field) {
	logger.Fatal(msg, fields...)
}

func (l *Logger) Debug(msg string, fields ...zap.Field) {
	logger.Debug(msg, fields...)
}

func (l *Logger) Panic(msg string, fields ...zap.Field) {
	logger.Panic(msg, fields...)
}

func (l *Logger) InfoF(format string, args ...interface{}) {
	msg := fmt.Sprintf(format, args...)
	logger.Info(msg)
}

func (l *Logger) WarnF(format string, args ...interface{}) {
	msg := fmt.Sprintf(format, args...)
	logger.Warn(msg)
}

func (l *Logger) ErrorF(format string, args ...interface{}) {
	msg := fmt.Sprintf(format, args...)
	logger.Error(msg)
}

func (l *Logger) FatalF(format string, args ...interface{}) {
	msg := fmt.Sprintf(format, args...)
	logger.Fatal(msg)
}

func (l *Logger) DebugF(format string, args ...interface{}) {
	msg := fmt.Sprintf(format, args...)
	logger.Debug(msg)
}

func (l *Logger) PanicF(format string, args ...interface{}) {
	msg := fmt.Sprintf(format, args...)
	logger.Panic(msg)
}
