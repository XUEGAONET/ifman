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
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/vishvananda/netlink"
)

type Learning struct {
	Name       string `yaml:"name"`
	LearningOn bool   `yaml:"learning_on"`
}

func UpdateLearning(learning *Learning) error {
	logrus.Debugf("update learning: %#v", learning)

	link, err := netlink.LinkByName(learning.Name)
	if err != nil {
		return errors.WithStack(err)
	}

	err = netlink.LinkSetLearning(link, learning.LearningOn)
	if err != nil {
		return errors.WithStack(err)
	}

	return nil
}
