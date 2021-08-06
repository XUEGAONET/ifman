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

package writer

import (
	"errors"
	"io"
	"time"

	rotatelogs "github.com/lestrrat/go-file-rotatelogs"
)

type Rotate struct {
	w io.Writer
}

func (r *Rotate) GetWriter() io.Writer {
	return r.w
}

func NewRotate(dir string, maxAge int, period int) (*Rotate, error) {
	if maxAge <= 0 || period <= 0 {
		return nil, errors.New("wrong max_age or period range")
	}

	f, err := rotatelogs.New(
		dir+"/%Y%m%d%H%M.log",
		rotatelogs.WithMaxAge(time.Duration(maxAge)*time.Second),
		rotatelogs.WithRotationTime(time.Duration(period)*time.Second),
	)
	if err != nil {
		return nil, err
	}

	r := Rotate{
		w: f,
	}
	return &r, nil
}
