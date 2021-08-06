package main

import (
	"fmt"
	"github.com/XUEGAONET/ifman/pkg/log"
	"github.com/XUEGAONET/ifman/pkg/log/writer"
	"github.com/pkg/errors"
	"strings"
)

func initLogger(logger *Logger) error {
	var mode log.Writer = nil
	var err error = nil

	switch strings.ToLower(logger.Mode) {
	case "none":
		mode = writer.NewNone()
	case "rotate":
		mode, err = writer.NewRotate(logger.Rotate.Dir, int(logger.Rotate.MaxAgeSec), int(logger.Rotate.PeriodSec))
		if err != nil {
			return errors.WithStack(err)
		}
	case "single":
		mode, err = writer.NewSingle(logger.Single.Path, logger.Single.Permit)
		if err != nil {
			return errors.WithStack(err)
		}
	case "stdout":
		mode = writer.NewStdout()
	default:
		return fmt.Errorf("unsupported logger mode")
	}

	err = log.SetLog(logger.Level, mode, logger.SyslogOn)
	if err != nil {
		return errors.WithStack(err)
	}

	return nil
}
