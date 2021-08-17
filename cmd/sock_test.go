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
	"github.com/sirupsen/logrus"
	"strings"
	"testing"
	"time"
)

func TestListenCtl(t *testing.T) {
	logrus.SetLevel(logrus.TraceLevel)
	var port uint16 = 11111

	go func() {
		err := ListenCtl(port)
		if err != nil {
			panic(err)
		}
	}()

	time.Sleep(time.Second)

	err := SendReload(port)
	if err != nil {
		panic(err)
	}

	err = ListenCtl(port)
	if err == nil {
		panic("nil pointer! bug!")
	}
	if !strings.Contains(err.Error(), "second") {
		panic("bug")
	}

	time.Sleep(time.Second)
}
