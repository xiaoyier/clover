package log

import (
	"flag"
	"os"
	"path"
	"time"

	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	"github.com/rifflock/lfshook"
	"github.com/sirupsen/logrus"
)

var (
	logPath   = flag.String("log_path", "", "the path save the log file")
	logFile   = flag.String("log_file", "latest.log", "log file name")
	logLevel  = flag.String("log_level", "info", "log level")
	logMaxAge = flag.Int("log_max_age", 24*7, "the log file max age(By Hour)")
)

var logger *logrus.Logger

func InitLogger() {

	logger = logrus.New()
	logrus.SetOutput(os.Stdout)
	logrus.SetFormatter(&logrus.JSONFormatter{
		TimestampFormat: "2006-01-02 15:04:05.000",
	})

	switch *logLevel {
	case "debug":
		logrus.SetLevel(logrus.DebugLevel)
	case "info":
		logrus.SetLevel(logrus.InfoLevel)
	case "warn":
		logrus.SetLevel(logrus.WarnLevel)
	case "error":
		logrus.SetLevel(logrus.ErrorLevel)
	default:
		logrus.SetLevel(logrus.InfoLevel)
	}

	if len(*logPath) > 0 {
		// Add file output hook
		logrus.AddHook(NewHook(*logPath, *logFile, *logMaxAge))
	}
}

func NewHook(logPath, fileName string, maxAge int) logrus.Hook {

	writer, _ := rotatelogs.New(
		path.Join(logPath, "clover-%Y-%m-%d-%H.log"),
		rotatelogs.WithLinkName(fileName),
		rotatelogs.WithMaxAge(time.Hour*time.Duration(maxAge)),
		rotatelogs.WithRotationTime(time.Hour),
	)

	return lfshook.NewHook(writer, &logrus.JSONFormatter{
		TimestampFormat: "2006-01-02 15:03:04.000",
	})
}

func WithCategory(category string) *logrus.Entry {
	return logger.WithField("category", category)
}

func WithTraceID(trace string) *logrus.Entry {
	return logger.WithField("trace_id", trace)
}
