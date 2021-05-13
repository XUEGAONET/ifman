package logger

import (
	"github.com/sirupsen/logrus"
	"testing"
	"time"
)

func TestSetGlobal1(t *testing.T) {
	err := SetGlobal(LevelTrace, ConfigStdout{})
	if err != nil {
		panic(err)
	}

	logrus.Errorln("stdout")
}

func TestSetGlobal2(t *testing.T) {
	c := ConfigSingleFile{}
	c.SetFilePath("/tmp/test.log")
	err := c.SetPermission(0600)
	if err != nil {
		panic(err)
	}

	err = SetGlobal(LevelTrace, c)
	if err != nil {
		panic(err)
	}

	logrus.Errorln("single file")
}

func TestSetGlobal3(t *testing.T) {
	c := ConfigAutoRotate{}
	c.SetCutPeriod(24 * time.Hour)
	c.SetMaxAge(15 * 24 * time.Hour)
	c.SetDirPath("/tmp")

	err := SetGlobal(LevelTrace, c)
	if err != nil {
		panic(err)
	}

	logrus.Errorln("auto rotate 1")
	time.Sleep(2 * time.Second)
	logrus.Errorln("auto rotate 2")
}

func TestSetGlobal4(t *testing.T) {
	c := ConfigSyslog{}

	err := SetGlobal(LevelTrace, c)
	if err != nil {
		panic(err)
	}

	logrus.Infof("teset")
}
