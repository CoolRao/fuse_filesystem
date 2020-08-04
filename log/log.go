package log

import (
	"github.com/sirupsen/logrus"
	"io"
	"os"
)

var Logger *logrus.Logger

func InitLogger() {
	var writers []io.Writer
	writers = append(writers, os.Stdout)
	fileAndStdoutWriter := io.MultiWriter(writers...)
	logger := logrus.New()
	logger.Out = fileAndStdoutWriter
	logger.SetLevel(logrus.DebugLevel)
	customFormatter := new(logrus.TextFormatter)
	customFormatter.FullTimestamp = true
	customFormatter.TimestampFormat = "2006-01-02 15:04:05"
	customFormatter.ForceColors = true
	logger.SetFormatter(customFormatter)
	Logger = logger
}
