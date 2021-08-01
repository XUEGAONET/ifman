/**
 * Copyright (c) 2021 Harris <ic0xgkk@gmail.com>
 *
 * Permission is hereby granted, free of charge, to any person obtaining a copy of
 * this software and associated documentation files (the "Software"), to deal in
 * the Software without restriction, including without limitation the rights to
 * use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of
 * the Software, and to permit persons to whom the Software is furnished to do so,
 * subject to the following conditions:
 *
 * The above copyright notice and this permission notice shall be included in all
 * copies or substantial portions of the Software.
 *
 * THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
 * IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS
 * FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR
 * COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER
 * IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN
 * CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.
 */

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
