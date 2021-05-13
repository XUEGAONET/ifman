package logger

import (
	"fmt"
	"time"
)

type ConfigSingleFile struct {
	filePath   string
	permission uint32
}

func (receiver *ConfigSingleFile) SetFilePath(s string) {
	receiver.filePath = s
}

func (receiver *ConfigSingleFile) SetPermission(p uint32) error {
	if p > 1000 {
		return fmt.Errorf("invalid permission mask")
	}

	receiver.permission = p
	return nil
}

type ConfigAutoRotate struct {
	dirPath      string
	rotatePeriod time.Duration
	rotateMaxAge time.Duration
}

func (receiver *ConfigAutoRotate) SetDirPath(s string) {
	receiver.dirPath = s
}

func (receiver *ConfigAutoRotate) SetCutPeriod(t time.Duration) {
	receiver.rotatePeriod = t
}

func (receiver *ConfigAutoRotate) SetMaxAge(t time.Duration) {
	receiver.rotateMaxAge = t
}

type ConfigStdout struct{}

type ConfigSyslog struct{}
