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
	"github.com/sirupsen/logrus"
)

func SetLog(level string, w Writer, useSyslog bool) error {
	levelInt, err := LevelStr2Int(level)
	if err != nil {
		return err
	}

	logrusLevel, err := Level2Logrus(levelInt)
	if err != nil {
		return err
	}

	logrus.SetLevel(logrus.Level(logrusLevel))

	//logrus.SetFormatter(&logrus.JSONFormatter{
	//	TimestampFormat:   "2006-01-02 15:04:05",
	//	DisableTimestamp:  false,
	//	DisableHTMLEscape: true,
	//	DataKey:           "",
	//	FieldMap:          nil,
	//	CallerPrettyfier:  nil,
	//	PrettyPrint:       false,
	//})

	logrus.SetFormatter(&logrus.TextFormatter{
		ForceColors:               false,
		DisableColors:             false,
		ForceQuote:                false,
		DisableQuote:              true,
		EnvironmentOverrideColors: false,
		DisableTimestamp:          false,
		FullTimestamp:             false,
		TimestampFormat:           "2006-01-02 15:04:05",
		DisableSorting:            false,
		SortingFunc:               nil,
		DisableLevelTruncation:    false,
		PadLevelText:              false,
		QuoteEmptyFields:          false,
		FieldMap:                  nil,
		CallerPrettyfier:          nil,
	})

	logrus.SetOutput(w.GetWriter())

	if useSyslog {
		err = enableSyslog()
		if err != nil {
			return err
		}
	}

	return nil
}
