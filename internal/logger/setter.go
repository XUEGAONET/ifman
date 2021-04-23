package logger

import (
	"fmt"
	rotatelogs "github.com/lestrrat/go-file-rotatelogs"
	"github.com/sirupsen/logrus"
	"io"
	"os"
)

const (
	LevelTrace   = 6
	LevelDebug   = 5
	LevelInfo    = 4
	LevelWarning = 3
	LevelError   = 2
	LevelFatal   = 1
)

func SetGlobal(level int, cc interface{}) error {
	if cc == nil {
		return fmt.Errorf("nil config error")
	}

	setLevel(level)
	setFormatter()

	writer, err := getWriter(cc)
	if err != nil {
		return err
	}
	setWriter(writer)

	return nil
}

func getWriter(c interface{}) (io.Writer, error) {
	var result io.Writer
	var err error
	switch c := c.(type) {
	case ConfigStdout:
		result = os.Stdout
	case ConfigSingleFile:
		result, err = os.OpenFile(c.filePath, os.O_RDWR|os.O_CREATE|os.O_APPEND, os.FileMode(c.permission))
	case ConfigAutoRotate:
		result, err = rotatelogs.New(
			c.dirPath+"/%Y%m%d%H%M.log",
			rotatelogs.WithMaxAge(c.rotateMaxAge),
			rotatelogs.WithRotationTime(c.rotatePeriod),
		)
	}

	if err != nil {
		return nil, err
	}

	return result, err
}

func setWriter(w io.Writer) {
	logrus.SetOutput(w)
}

func setLevel(l int) {
	logrus.SetLevel(logrus.Level(l))
}

func setFormatter() {
	logrus.SetFormatter(&logrus.JSONFormatter{
		TimestampFormat:   "2006-01-02 15:04:05",
		DisableTimestamp:  false,
		DisableHTMLEscape: true,
		DataKey:           "",
		FieldMap:          nil,
		CallerPrettyfier:  nil,
		PrettyPrint:       false,
	})
}
