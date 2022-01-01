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
	"github.com/XUEGAONET/ifman/pkg/log/writer"
	"testing"

	"github.com/sirupsen/logrus"
)

func TestStdoutWriter(t *testing.T) {
	w := writer.NewStdout()
	setWriterAndGenMessage(w)
}

func TestNoneWriter(t *testing.T) {
	w := writer.NewNone()
	setWriterAndGenMessage(w)
}

func TestSingleWriter(t *testing.T) {
	w, err := writer.NewSingle("/tmp/test_ifman.log", 0644)
	if err != nil {
		panic(err)
	}
	setWriterAndGenMessage(w)
}

func TestRotateWriter(t *testing.T) {
	w, err := writer.NewRotate("/tmp", 600, 60)
	if err != nil {
		panic(err)
	}
	setWriterAndGenMessage(w)
}

func setWriterAndGenMessage(w Writer) {
	err := SetLog("trace", w)
	if err != nil {
		panic(err)
	}

	logrus.Error("test log writer")
}
