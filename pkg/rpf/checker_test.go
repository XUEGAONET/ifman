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

package rpf

import (
	"testing"
)

const (
	name string = "eth0"
)

func TestRead(t *testing.T) {
	mode, err := Read(name)
	if err != nil {
		panic(err)
	}

	err = Write(name, RPF_LOOSE)
	if err != nil {
		panic(err)
	}

	loose, err := Read(name)
	if err != nil {
		panic(err)
	}

	if loose != RPF_LOOSE {
		panic("loose")
	}

	err = CheckAndFix(name, RPF_STRICT)
	if err != nil {
		panic(err)
	}

	strict, err := Read(name)
	if err != nil {
		panic(err)
	}

	if strict != RPF_STRICT {
		panic("strict")
	}

	err = Write(name, mode)
	if err != nil {
		panic("write")
	}
}
