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

package log

import (
	"errors"
	"log/syslog"
	"strings"

	"github.com/sirupsen/logrus"
)

const (
	LevelTrace int = iota
	LevelDebug
	LevelInfo
	LevelWarning
	LevelError
	LevelFatal
)

func LevelStr2Int(level string) (int, error) {
	switch strings.ToLower(level) {
	case "trace":
		return LevelTrace, nil
	case "debug":
		return LevelDebug, nil
	case "info":
		return LevelInfo, nil
	case "warning":
		return LevelWarning, nil
	case "warn":
		return LevelWarning, nil
	case "error":
		return LevelError, nil
	case "fatal":
		return LevelFatal, nil
	default:
		return 0, errors.New("wrong level string")
	}
}

func LevelInt2Str(level int) (string, error) {
	switch level {
	case LevelFatal:
		return "fatal", nil
	case LevelError:
		return "error", nil
	case LevelWarning:
		return "warn", nil
	case LevelInfo:
		return "info", nil
	case LevelDebug:
		return "debug", nil
	case LevelTrace:
		return "trace", nil
	default:
		return "", errors.New("wrong int level range")
	}
}

func Level2Syslog(level int) (syslog.Priority, error) {
	switch level {
	case LevelFatal:
		return syslog.LOG_EMERG, nil
	case LevelError:
		return syslog.LOG_ERR, nil
	case LevelWarning:
		return syslog.LOG_WARNING, nil
	case LevelInfo:
		return syslog.LOG_NOTICE, nil
	case LevelDebug:
		return syslog.LOG_INFO, nil
	case LevelTrace:
		return syslog.LOG_DEBUG, nil
	default:
		return 0, errors.New("wrong int level range")
	}
}

func Logrus2Level(level int) (int, error) {
	switch level {
	case int(logrus.FatalLevel):
		return LevelFatal, nil
	case int(logrus.ErrorLevel):
		return LevelError, nil
	case int(logrus.WarnLevel):
		return LevelWarning, nil
	case int(logrus.InfoLevel):
		return LevelInfo, nil
	case int(logrus.DebugLevel):
		return LevelDebug, nil
	case int(logrus.TraceLevel):
		return LevelTrace, nil
	default:
		return 0, errors.New("wrong logrus level range")
	}
}

func Level2Logrus(level int) (int, error) {
	switch level {
	case LevelFatal:
		return int(logrus.FatalLevel), nil
	case LevelError:
		return int(logrus.ErrorLevel), nil
	case LevelWarning:
		return int(logrus.WarnLevel), nil
	case LevelInfo:
		return int(logrus.InfoLevel), nil
	case LevelDebug:
		return int(logrus.DebugLevel), nil
	case LevelTrace:
		return int(logrus.TraceLevel), nil
	default:
		return 0, errors.New("wrong int level range")
	}
}
