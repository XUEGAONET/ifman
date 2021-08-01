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
	"testing"

	xWriter "github.com/XUEGAONET/ifman/log/writer"
	"github.com/sirupsen/logrus"
)

func TestStdoutWriter(t *testing.T) {
	w := xWriter.NewStdout()
	setWriterAndGenMessage(w)
}

func TestNoneWriter(t *testing.T) {
	w := xWriter.NewNone()
	setWriterAndGenMessage(w)
}

func TestSingleWriter(t *testing.T) {
	w, err := xWriter.NewSingle("/tmp/test_ifman.log", 0644)
	if err != nil {
		panic(err)
	}
	setWriterAndGenMessage(w)
}

func TestRotateWriter(t *testing.T) {
	w, err := xWriter.NewRotate("/tmp", 600, 60)
	if err != nil {
		panic(err)
	}
	setWriterAndGenMessage(w)
}

func setWriterAndGenMessage(w Writer) {
	err := SetLog("trace", w, false)
	if err != nil {
		panic(err)
	}

	logrus.Error("test log writer")
}
