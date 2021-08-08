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
	rpfp "github.com/XUEGAONET/ifman/utils/rpf"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

type RpFilter struct {
	Name string `yaml:"name"`
	Mode string `yaml:"mode"`
}

func UpdateRpFilter(rpf *RpFilter) error {
	logrus.Debugf("update rp_filter: %#v", rpf)

	var want rpfp.RPFType

	switch rpf.Mode {
	case "off":
		want = rpfp.RPF_NONE
	case "strict":
		want = rpfp.RPF_STRICT
	case "loose":
		want = rpfp.RPF_LOOSE
	default:
		return errors.WithStack(fmt.Errorf("unsupported rp_filter mode"))
	}

	err := rpfp.CheckAndFix(rpf.Name, want)
	if err != nil {
		return errors.WithStack(err)
	}

	return nil
}
