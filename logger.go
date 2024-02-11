package l

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"log"
	"os"
	"time"

	"github.com/segmentio/kafka-go"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type KafkaSyncer struct {
	writer *kafka.Writer
	topic  string
}

func NewKafkaSyncer(brokers []string, topic string) (*KafkaSyncer, error) {
	keypair, err := tls.LoadX509KeyPair("service.cert", "service.key")
	if err != nil {
		log.Fatalf("Failed to load access key and/or access certificate: %s", err)
	}

	caCert, err := os.ReadFile("ca.pem")
	if err != nil {
		log.Fatalf("Failed to read CA certificate file: %s", err)
	}

	caCertPool := x509.NewCertPool()
	ok := caCertPool.AppendCertsFromPEM(caCert)
	if !ok {
		log.Fatalf("Failed to parse CA certificate file: %s", err)
	}

	dialer := &kafka.Dialer{
		Timeout:   10 * time.Second,
		DualStack: true,
		TLS: &tls.Config{
			Certificates: []tls.Certificate{keypair},
			RootCAs:      caCertPool,
		},
	}

	writer := kafka.NewWriter(kafka.WriterConfig{
		Brokers: brokers,
		Topic:   topic,
		// Balancer: &kafka.LeastBytes{},
		Dialer: dialer,
	})

	return &KafkaSyncer{writer: writer, topic: topic}, nil
}

func (k *KafkaSyncer) Write(p []byte) (n int, err error) {
	err = k.writer.WriteMessages(context.Background(),
		kafka.Message{
			Value: p,
		},
	)
	if err != nil {
		return 0, err
	}
	return len(p), nil
}

func (k *KafkaSyncer) Sync() error {
	return k.writer.Close()
}

// Init initializes the logger.
func CreateLogger() *zap.Logger {
	// log.Println("Creating logger")
	stdout := zapcore.AddSync(os.Stdout)

	// file := zapcore.AddSync(&lumberjack.Logger{
	// 	Filename:   "logs/app.log",
	// 	MaxSize:    10, // megabytes
	// 	MaxBackups: 3,
	// 	MaxAge:     7, // days
	// })

	level := zap.NewAtomicLevelAt(zap.InfoLevel)

	// productionCfg := zap.NewProductionEncoderConfig()
	// productionCfg.TimeKey = "timestamp"
	// productionCfg.EncodeTime = zapcore.ISO8601TimeEncoder

	developmentCfg := zap.NewDevelopmentEncoderConfig()
	developmentCfg.EncodeLevel = zapcore.CapitalColorLevelEncoder

	consoleEncoder := zapcore.NewConsoleEncoder(developmentCfg)
	// fileEncoder := zapcore.NewJSONEncoder(productionCfg)

	core := zapcore.NewTee(
		zapcore.NewCore(consoleEncoder, stdout, level),
		// zapcore.NewCore(fileEncoder, file, level),
	)

	return zap.New(
		core,
		zap.AddCaller(),
		zap.AddCallerSkip(1),
	)
}
func KafkaLogger() *zap.Logger {
	kafkaSyncer, err := NewKafkaSyncer([]string{"kafka-257bfc54-lakjos-f2b6.a.aivencloud.com:19190"}, "logs")
	if err != nil {
		panic(err)
	}

	level := zap.NewAtomicLevelAt(zap.InfoLevel)
	stdout := zapcore.AddSync(os.Stdout)
	developmentCfg := zap.NewDevelopmentEncoderConfig()
	consoleEncoder := zapcore.NewConsoleEncoder(developmentCfg)

	logger := zap.New(
		zapcore.NewTee(
			zapcore.NewCore(
				zapcore.NewJSONEncoder(
					zap.NewProductionEncoderConfig(),
				),
				kafkaSyncer,
				zapcore.DebugLevel,
			),
			zapcore.NewCore(consoleEncoder, stdout, level),
		),
		zap.AddCaller(),
		zap.AddCallerSkip(1),
	)

	return logger
	// logger.Info("This log will be sent to Kafka")
}

var logger = CreateLogger()

func Info(msg string, fields ...zap.Field) {
	logger.Info(msg, fields...)
}

// func Warn(msg string, fields ...zap.Field) {
// 	logger.Warn(msg, fields...)
// }

// func Error(msg string, fields ...zap.Field) {
// 	logger.Error(msg, fields...)
// }

// func Fatal(msg string, fields ...zap.Field) {
// 	logger.Fatal(msg, fields...)
// }

// func Debug(msg string, fields ...zap.Field) {
// 	logger.Debug(msg, fields...)
// }

// func Panic(msg string, fields ...zap.Field) {
// 	logger.Panic(msg, fields...)
// }

// func InfoF(format string, args ...interface{}) {
// 	msg := fmt.Sprintf(format, args...)
// 	logger.Info(msg)
// }

// func WarnF(format string, args ...interface{}) {
// 	msg := fmt.Sprintf(format, args...)
// 	logger.Warn(msg)
// }

// func ErrorF(format string, args ...interface{}) {
// 	msg := fmt.Sprintf(format, args...)
// 	logger.Error(msg)
// }

// func FatalF(format string, args ...interface{}) {
// 	msg := fmt.Sprintf(format, args...)
// 	logger.Fatal(msg)
// }

// func DebugF(format string, args ...interface{}) {
// 	msg := fmt.Sprintf(format, args...)
// 	logger.Debug(msg)
// }

// func PanicF(format string, args ...interface{}) {
// 	msg := fmt.Sprintf(format, args...)
// 	logger.Panic(msg)
// }
