package main

import (
	"flag"
	"fmt"
	"ifman/internal/logger"
	"io/ioutil"
	"os"
	"strings"
	"time"

	"github.com/sirupsen/logrus"
	"gopkg.in/yaml.v3"
)

func main() {
	conf := RootConfig{}

	arg := flag.String("config", "config.yaml", "config path")
	flag.Parse()

	confArr, err := openConf(*arg)
	if err != nil {
		panic(err)
	}

	err = yaml.Unmarshal(confArr, &conf)
	if err != nil {
		panic(err)
	}

	err = initLogger(conf.Logger)
	if err != nil {
		panic(err)
	}

	logrus.Debugf("config dump: %#v", conf)

	for {
		logrus.Infof("interface and route operating begin")

		logrus.Debugf("interface verifying begin")
		suInterface(conf.Interface)
		logrus.Debugf("interface verifying end")

		logrus.Debugf("route verifying begin")
		suRoute(conf.Route)
		logrus.Debugf("route verifying end")

		logrus.Infof("interface and route operating finished")
		time.Sleep(time.Duration(conf.Verify.CheckPeriod) * time.Second)
	}
}

func suInterface(cs []Interface) {
	for i, _ := range cs {
		c := cs[i]

		logrus.WithField("signature", "suInterface").Debugf("%s interface operating begin", c.Name)
		err := afInterface(c)
		if err != nil {
			logrus.WithField("signature", "suInterface").
				WithField("interface", c.Name).
				Errorf("afInterface failed: %v", err)
		}
	}
}

func suRoute(rs []Route) {
	for i, _ := range rs {
		logrus.WithField("signature", "suRoute").Debugf("%s route operating begin", rs[i].Dst)
		err := afRoute(rs[i])
		if err != nil {
			logrus.WithField("signature", "suRoute").
				WithField("route", rs[i].Dst+"->"+rs[i].Gw).
				Errorf("afRoute failed: %v", err)
		}
	}
}

func initLogger(l Logger) error {
	var cc interface{} = nil
	switch l.Mode {
	case "stdout":
		cc = logger.ConfigStdout{}
	case "single":
		t := logger.ConfigSingleFile{}
		t.SetFilePath(l.Single.Path)
		err := t.SetPermission(uint32(l.Single.Perm))
		if err != nil {
			return err
		}

		cc = t
	case "rotate":
		t := logger.ConfigAutoRotate{}
		t.SetDirPath(l.Rotate.Dir)
		t.SetCutPeriod(time.Duration(l.Rotate.CutPeriod) * time.Second)
		t.SetMaxAge(time.Duration(l.Rotate.Age) * time.Second)

		cc = t
	case "syslog":
		cc = logger.ConfigSyslog{}
	default:
		return fmt.Errorf("invalid logger mode")
	}

	var level int
	switch strings.ToLower(l.Level) {
	case "trace":
		level = logger.LevelTrace
	case "debug":
		level = logger.LevelDebug
	case "info":
		level = logger.LevelInfo
	case "warning", "warn":
		level = logger.LevelWarning
	case "error":
		level = logger.LevelError
	default:
		return fmt.Errorf("invalid logger level")
	}

	return logger.SetGlobal(level, cc)
}

func openConf(p string) ([]byte, error) {
	f, err := os.Open(p)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	arr, err := ioutil.ReadAll(f)
	if err != nil {
		return nil, err
	}

	return arr, nil
}
