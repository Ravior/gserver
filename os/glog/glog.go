package glog

import (
	"github.com/Ravior/gserver/os/glog/internal"
	gconfig2 "github.com/Ravior/gserver/util/gconfig"
	"github.com/lestrrat-go/file-rotatelogs"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"io"
	"os"
	"strings"
	"sync"
	"time"
)

var (
	// DefaultLogger 默认日志对象（控制台输出)
	DefaultLogger *Logger
	loggers       map[string]*Logger
	rw            sync.RWMutex
)

type Logger struct {
	*zap.SugaredLogger
	*gconfig2.LogConfig
}

func init() {
	loggers = make(map[string]*Logger)
}

// Init 日志模块初始化
func Init(serverId string) {
	for name, _ := range gconfig2.Global.Log {
		NewLogger(name, serverId)
	}
	if logger, ok := loggers["default"]; ok {
		DefaultLogger = logger
	} else {
		DefaultLogger = NewConfigLogger(internal.NewConsoleLogConfig(), "")
	}
}

func Get(name ...string) *Logger {
	var loggerName = "default"
	if len(name) > 0 {
		loggerName = name[0]
	}
	if logger, ok := loggers[loggerName]; ok {
		return logger
	}
	return DefaultLogger
}

func NewLogger(name string, serverId string, opts ...zap.Option) *Logger {
	if name == "" {
		return nil
	}

	rw.Lock()
	defer rw.Unlock()

	if logger, ok := loggers[name]; ok {
		return logger
	}

	logConfig, err := internal.GetLogConfig(name)
	if err != nil {
		panic(err)
	}

	logger := NewConfigLogger(logConfig, serverId, opts...)
	loggers[name] = logger

	return logger

}

func NewConfigLogger(config *gconfig2.LogConfig, serverId string, opts ...zap.Option) *Logger {
	encoderConfig := zapcore.EncoderConfig{
		TimeKey:        "ts",
		LevelKey:       "level",
		CallerKey:      "caller",
		MessageKey:     "msg",
		NameKey:        "name",
		StacktraceKey:  "stack",
		FunctionKey:    zapcore.OmitKey,
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.LowercaseLevelEncoder, // 将日志级别字符串转化为小写
		EncodeTime:     zapcore.ISO8601TimeEncoder,    // 时间转化为秒
		EncodeDuration: zapcore.StringDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
		EncodeName:     zapcore.FullNameEncoder,
	}

	if config.PrintCaller {
		opts = append(opts, zap.AddCaller(), zap.AddCallerSkip(1))
	}

	level := internal.GetLevel(config.Level)
	opts = append(opts, zap.AddStacktrace(internal.GetLevel(config.StackLevel)))

	var defaultWriters []zapcore.WriteSyncer
	var errorWriters []zapcore.WriteSyncer
	if config.EnableConsole {
		defaultWriters = append(defaultWriters, zapcore.AddSync(os.Stdout))
	}

	if config.EnableWriteFile && config.FilePath != "" {
		// 默认格式 "%Y-%m-%d.log"
		fileFormat := "%Y-%m-%d.log"
		if config.FileFormat != "" {
			fileFormat = config.FileFormat
		}
		defaultLogFile := strings.Replace(config.FilePath, ".log", "", -1) + "-" + fileFormat
		if serverId != "" {
			defaultLogFile = strings.Replace(config.FilePath, ".log", "", -1) + "-" + serverId + "-" + fileFormat
		}

		defaultHook := getWriter(defaultLogFile, config.MaxAge, config.RotationTime)
		defaultWriters = append(defaultWriters, zapcore.AddSync(defaultHook))

		errorLogFile := strings.Replace(config.FilePath, ".log", "-error", -1) + "-" + fileFormat
		if serverId != "" {
			errorLogFile = strings.Replace(config.FilePath, ".log", "-error", -1) + "-" + serverId + "-" + fileFormat
		}
		errorHook := getWriter(errorLogFile, config.MaxAge, config.RotationTime)
		errorWriters = append(errorWriters, zapcore.AddSync(errorHook))
	}

	logger := &Logger{
		SugaredLogger: internal.NewSugaredLogger(encoderConfig, zapcore.NewMultiWriteSyncer(defaultWriters...), zapcore.NewMultiWriteSyncer(errorWriters...), level, opts...),
		LogConfig:     config,
	}

	return logger
}

func getWriter(filename string, maxAge int, rotationTime int) io.Writer {
	// 默认为24小时切割
	if rotationTime == 0 {
		rotationTime = 24
	}
	writer, _ := rotatelogs.New(
		filename,
		rotatelogs.WithMaxAge(time.Duration(maxAge*24)*time.Hour),
		rotatelogs.WithRotationTime(time.Duration(rotationTime)*time.Hour),
	)
	return writer
}

func Debug(args ...interface{}) {
	if DefaultLogger == nil {
		DefaultLogger = NewConfigLogger(internal.NewConsoleLogConfig(), "")
	}
	DefaultLogger.Debug(args...)
}

func Info(args ...interface{}) {
	if DefaultLogger == nil {
		DefaultLogger = NewConfigLogger(internal.NewConsoleLogConfig(), "")
	}
	DefaultLogger.Info(args...)
}

// Warn uses fmt.Sprint to construct and log a message.
func Warn(args ...interface{}) {
	if DefaultLogger == nil {
		DefaultLogger = NewConfigLogger(internal.NewConsoleLogConfig(), "")
	}
	DefaultLogger.Warn(args...)
}

// Error uses fmt.Sprint to construct and log a message.
func Error(args ...interface{}) {
	if DefaultLogger == nil {
		DefaultLogger = NewConfigLogger(internal.NewConsoleLogConfig(), "")
	}
	DefaultLogger.Error(args...)
}

// DPanic uses fmt.Sprint to construct and log a message. In development, the
// glog then panics. (See DPanicLevel for details.)
func DPanic(args ...interface{}) {
	if DefaultLogger == nil {
		DefaultLogger = NewConfigLogger(internal.NewConsoleLogConfig(), "")
	}
	DefaultLogger.DPanic(args...)
}

// Panic uses fmt.Sprint to construct and log a message, then panics.
func Panic(args ...interface{}) {
	DefaultLogger.Panic(args...)
}

// Fatal uses fmt.Sprint to construct and log a message, then calls os.Exit.
func Fatal(args ...interface{}) {
	if DefaultLogger == nil {
		DefaultLogger = NewConfigLogger(internal.NewConsoleLogConfig(), "")
	}
	DefaultLogger.Fatal(args...)
}

// Debugf uses fmt.Sprintf to log a templated message.
func Debugf(template string, args ...interface{}) {
	if DefaultLogger == nil {
		DefaultLogger = NewConfigLogger(internal.NewConsoleLogConfig(), "")
	}
	DefaultLogger.Debugf(template, args...)
}

// Infof uses fmt.Sprintf to log a templated message.
func Infof(template string, args ...interface{}) {
	if DefaultLogger == nil {
		DefaultLogger = NewConfigLogger(internal.NewConsoleLogConfig(), "")
	}
	DefaultLogger.Infof(template, args...)
}

// Warnf uses fmt.Sprintf to log a templated message.
func Warnf(template string, args ...interface{}) {
	if DefaultLogger == nil {
		DefaultLogger = NewConfigLogger(internal.NewConsoleLogConfig(), "")
	}
	DefaultLogger.Warnf(template, args...)
}

// Errorf uses fmt.Sprintf to log a templated message.
func Errorf(template string, args ...interface{}) {
	if DefaultLogger == nil {
		DefaultLogger = NewConfigLogger(internal.NewConsoleLogConfig(), "")
	}
	DefaultLogger.Errorf(template, args...)
}

// DPanicf uses fmt.Sprintf to log a templated message. In development, the
// glog then panics. (See DPanicLevel for details.)
func DPanicf(template string, args ...interface{}) {
	if DefaultLogger == nil {
		DefaultLogger = NewConfigLogger(internal.NewConsoleLogConfig(), "")
	}
	DefaultLogger.DPanicf(template, args...)
}

// Panicf uses fmt.Sprintf to log a templated message, then panics.
func Panicf(template string, args ...interface{}) {
	if DefaultLogger == nil {
		DefaultLogger = NewConfigLogger(internal.NewConsoleLogConfig(), "")
	}
	DefaultLogger.Panicf(template, args...)
}

// Fatalf uses fmt.Sprintf to log a templated message, then calls os.Exit.
func Fatalf(template string, args ...interface{}) {
	if DefaultLogger == nil {
		DefaultLogger = NewConfigLogger(internal.NewConsoleLogConfig(), "")
	}
	DefaultLogger.Fatalf(template, args...)
}

// Debugw logs a message with some additional context. The variadic key-value
// pairs are treated as they are in With.
//
// When debug-level logging is disabled, this is much faster than
//  s.With(keysAndValues).Debug(proto)
func Debugw(msg string, keysAndValues ...interface{}) {
	if DefaultLogger == nil {
		DefaultLogger = NewConfigLogger(internal.NewConsoleLogConfig(), "")
	}
	DefaultLogger.Debugw(msg, keysAndValues...)
}

// Infow logs a message with some additional context. The variadic key-value
// pairs are treated as they are in With.
func Infow(msg string, keysAndValues ...interface{}) {
	if DefaultLogger == nil {
		DefaultLogger = NewConfigLogger(internal.NewConsoleLogConfig(), "")
	}
	DefaultLogger.Infow(msg, keysAndValues...)
}

// Warnw logs a message with some additional context. The variadic key-value
// pairs are treated as they are in With.
func Warnw(msg string, keysAndValues ...interface{}) {
	if DefaultLogger == nil {
		DefaultLogger = NewConfigLogger(internal.NewConsoleLogConfig(), "")
	}
	DefaultLogger.Warnw(msg, keysAndValues...)
}

// Errorw logs a message with some additional context. The variadic key-value
// pairs are treated as they are in With.
func Errorw(msg string, keysAndValues ...interface{}) {
	if DefaultLogger == nil {
		DefaultLogger = NewConfigLogger(internal.NewConsoleLogConfig(), "")
	}
	DefaultLogger.Errorw(msg, keysAndValues...)
}

// DPanicw logs a message with some additional context. In development, the
// glog then panics. (See DPanicLevel for details.) The variadic key-value
// pairs are treated as they are in With.
func DPanicw(msg string, keysAndValues ...interface{}) {
	if DefaultLogger == nil {
		DefaultLogger = NewConfigLogger(internal.NewConsoleLogConfig(), "")
	}
	DefaultLogger.DPanicw(msg, keysAndValues...)
}

// Panicw logs a message with some additional context, then panics. The
// variadic key-value pairs are treated as they are in With.
func Panicw(msg string, keysAndValues ...interface{}) {
	if DefaultLogger == nil {
		DefaultLogger = NewConfigLogger(internal.NewConsoleLogConfig(), "")
	}
	DefaultLogger.Panicw(msg, keysAndValues...)
}

// Fatalw logs a message with some additional context, then calls os.Exit. The
// variadic key-value pairs are treated as they are in With.
func Fatalw(msg string, keysAndValues ...interface{}) {
	if DefaultLogger == nil {
		DefaultLogger = NewConfigLogger(internal.NewConsoleLogConfig(), "")
	}
	DefaultLogger.Fatalw(msg, keysAndValues...)
}
