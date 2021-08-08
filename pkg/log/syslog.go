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
	lSyslog "github.com/sirupsen/logrus/hooks/syslog"
)

// Calling must be after logrus initing
func enableSyslog() error {
	levelLevel, err := Logrus2Level(int(logrus.GetLevel()))
	if err != nil {
		return nil
	}

	syslogLevel, err := Level2Syslog(levelLevel)
	if err != nil {
		return nil
	}

	hook, err := lSyslog.NewSyslogHook("", "", syslogLevel, "ifman")
	if err != nil {
		return err
	}

	logrus.AddHook(hook)
	return nil
}
