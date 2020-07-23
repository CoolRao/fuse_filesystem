package log

import (
	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	"github.com/rifflock/lfshook"
	"github.com/sirupsen/logrus"
	"io"
	"os"
	"path"
	"time"
)


var Logger *logrus.Logger

func InitLogger(logPath string, logFileName string) error {
	fileName := path.Join(logPath, logFileName)
	var writers []io.Writer
	// todo
	if true {
		writers = append(writers, os.Stdout)
	} else {
		src, err := os.OpenFile(os.DevNull, os.O_APPEND|os.O_WRONLY, os.ModeAppend)
		if err != nil {
			panic(err)
		}
		writers = append(writers, src)
	}
	fileAndStdoutWriter := io.MultiWriter(writers...)
	logger := logrus.New()
	logger.Out = fileAndStdoutWriter
	logger.SetLevel(logrus.DebugLevel)
	customFormatter := new(logrus.TextFormatter)
	customFormatter.FullTimestamp = true
	customFormatter.TimestampFormat = "2006-01-02 15:04:05"
	customFormatter.ForceColors = true
	logger.SetFormatter(customFormatter)
	logWriter, err := rotatelogs.New(
		fileName+".%Y%m%d.log",
		//rotatelogs.WithLinkName(fileName),
		rotatelogs.WithMaxAge(7*24*time.Hour),
		rotatelogs.WithRotationTime(24*time.Hour),
	)
	if err != nil {
		return err
	}
	writeMap := lfshook.WriterMap{
		logrus.InfoLevel:  logWriter,
		logrus.FatalLevel: logWriter,
		logrus.DebugLevel: logWriter,
		logrus.WarnLevel:  logWriter,
		logrus.ErrorLevel: logWriter,
		logrus.PanicLevel: logWriter,
	}
	lfHook := lfshook.NewHook(writeMap, &logrus.JSONFormatter{
		TimestampFormat: "2006-01-02 15:04:05",
	})
	logger.AddHook(lfHook)
	Logger = logger
	return nil
}
