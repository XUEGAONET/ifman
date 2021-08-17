// Copyright 2021 The ifman authors https://github.com/XUEGAONET/ifman
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//    http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

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
