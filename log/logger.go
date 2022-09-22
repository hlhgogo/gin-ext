package log

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strings"
	"time"

	"github.com/go-errors/errors"
	"github.com/hlhgogo/config"
	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	"github.com/sirupsen/logrus"
)

const (
	// Type 日志类型
	Type = "app"
	// DefaultTimestampFormat 时间格式
	DefaultTimestampFormat = "2006-01-02 15:04:05"
)

var (
	log *logrus.Logger
)

// Info info log
func Info(args ...interface{}) {
	log.WithFields(logrus.Fields{
		"type": Type,
	}).Info(args...)
}

// Infof 格式化输出info log
func Infof(format string, args ...interface{}) {
	log.WithFields(logrus.Fields{
		"type": Type,
	}).Info(fmt.Sprintf(format, args...))
}

// InfoFields 格式化输出info log
func InfoFields(fields logrus.Fields, args ...interface{}) {
	fields["type"] = Type
	log.WithFields(fields).Info(args)
}

// Warn warnlog
func Warn(args ...interface{}) {
	log.WithFields(logrus.Fields{
		"type": Type,
	}).Warn(args...)
}

// Warnf 格式化输出warn log
func Warnf(format string, args ...interface{}) {
	log.WithFields(logrus.Fields{
		"type": Type,
	}).Warn(fmt.Sprintf(format, args...))
}

// WarnFields warnlog
func WarnFields(fields logrus.Fields, args ...interface{}) {
	fields["type"] = Type
	log.WithFields(fields).Warn(args)
}

// Error 打印错误对象
func Error(args ...interface{}) {
	err := errors.New(args)
	log.WithFields(logrus.Fields{
		"type":  Type,
		"stack": err.ErrorStack(),
	}).Error(args...)
}

// Errorf 打印错误信息
func Errorf(format string, args ...interface{}) {
	err := errors.New(fmt.Sprintf(format, args...))
	log.WithFields(logrus.Fields{
		"type":  Type,
		"stack": err.ErrorStack(),
	}).Error(args...)
}

// ErrorFields errorlog
func ErrorFields(fields logrus.Fields, args ...interface{}) {
	fields["type"] = Type
	log.WithFields(fields).Error(args)
}

// LineFormatter ...
type LineFormatter struct {
	// TimestampFormat to use for display when a full timestamp is printed
	TimestampFormat string
}

// Format implement the Formatter interface
func (f *LineFormatter) Format(entry *logrus.Entry) ([]byte, error) {
	var b *bytes.Buffer
	if entry.Buffer != nil {
		b = entry.Buffer
	} else {
		b = &bytes.Buffer{}
	}
	timestampFormat := f.TimestampFormat
	if timestampFormat == "" {
		timestampFormat = time.RFC3339
	}

	var field = ""
	if entry.Data != nil {
		if b, err := json.Marshal(entry.Data); err == nil {
			field = string(b)
		}
	}

	// log file and line
	//source := ""
	//fileLine, ok := entry.Data[Source]
	//if ok {
	//	if v, ok := fileLine.(string); ok {
	//		source = v
	//	}
	//}

	b.WriteString(fmt.Sprintf("%s [%s] [%s] - %s %s\n", config.Get().App.Name, strings.ToUpper(entry.Level.String()), entry.Time.Format(timestampFormat), entry.Message, field))

	return b.Bytes(), nil
}

// Setup ...
func Setup() {
	log = logrus.New()

	log.SetReportCaller(true)
	log.SetFormatter(&LineFormatter{TimestampFormat: DefaultTimestampFormat})

	//writer := GetProjectIoWriter()
	//writers := []io.Writer{writer}
	//
	//// Local development output to the console
	//if config.Get().App.ShowTrace {
	//	writers = append(writers, os.Stdout)
	//}

	// log.SetOutput(io.MultiWriter(writers...))

	// output console
	log.SetOutput(os.Stdout)
	log.SetFormatter(&LineFormatter{TimestampFormat: DefaultTimestampFormat})
	//log.SetFormatter(&log.JSONFormatter{TimestampFormat: DefaultTimestampFormat})

	// Set log level
	var level logrus.Level = logrus.TraceLevel
	switch config.Get().Logger.Level {
	case "trace":
		level = logrus.TraceLevel
	case "debug":
		level = logrus.DebugLevel
	case "warn":
		level = logrus.WarnLevel
	case "info":
		level = logrus.InfoLevel
	case "error":
		level = logrus.ErrorLevel
	case "fatal":
		level = logrus.FatalLevel
	}
	log.SetLevel(level)
	log.AddHook(NewContextHook(level))
}

// GetGinLogIoWriter gin日志保存规则ioWriter
func GetGinLogIoWriter() io.Writer {
	writer, err := rotatelogs.New(
		config.Get().Logger.SavePath+"/api-%Y-%m-%d.log",
		rotatelogs.WithMaxAge(time.Duration(config.Get().Logger.SaveDay)*24*time.Hour),       // Maximum file save time
		rotatelogs.WithRotationTime(time.Duration(config.Get().Logger.SaveDay)*24*time.Hour), // Log the cut interval
	)
	if err != nil {
		panic(err)
	}

	return writer
}

// GetProjectIoWriter 业务日志保存规则ioWriter
func GetProjectIoWriter() io.Writer {
	writer, err := rotatelogs.New(
		config.Get().Logger.SavePath+"/gin-%Y-%m-%d.log",
		rotatelogs.WithMaxAge(time.Duration(config.Get().Logger.SaveDay)*24*time.Hour),       // Maximum file save time
		rotatelogs.WithRotationTime(time.Duration(config.Get().Logger.SaveDay)*24*time.Hour), // Log the cut interval
	)
	if err != nil {
		panic(err)
	}

	return writer
}
