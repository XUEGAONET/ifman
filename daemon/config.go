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
	"github.com/XUEGAONET/ifman/common"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"gopkg.in/yaml.v3"
	"io/ioutil"
	"os"
	"sync"
)

func parseLocalConfig(path string) (*common.Config, error) {
	b, err := readFile(path)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	c := common.Config{}
	err = yaml.Unmarshal(b, &c)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return &c, nil
}

func readFile(path string) ([]byte, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	defer f.Close()

	return ioutil.ReadAll(f)
}

type dynamicConfig struct {
	conf *common.Config
	path string
	lock sync.RWMutex
}

var _globalConfig dynamicConfig

// getGlobalConfig 会返回现有配置。
// 当配置未被初始化时，会返回空指针。为防止异常，请在使用该函数前先进行初始化操作。
func getGlobalConfig() *common.Config {
	_globalConfig.lock.RLock()
	defer _globalConfig.lock.RUnlock()

	return _globalConfig.conf
}

func reloadGlobalConfig() error {
	_globalConfig.lock.Lock()
	defer _globalConfig.lock.Unlock()

	c, err := parseLocalConfig(_globalConfig.path)
	if err != nil {
		return errors.WithStack(err)
	}

	_globalConfig.conf = c

	logrus.Infof("global config reloaded")

	return nil
}

func initGlobalConfig(path string) error {
	_globalConfig = dynamicConfig{
		conf: nil,
		path: path,
		lock: sync.RWMutex{},
	}

	c, err := parseLocalConfig(path)
	if err != nil {
		return errors.WithStack(err)
	}

	_globalConfig.lock.Lock()
	defer _globalConfig.lock.Unlock()
	_globalConfig.conf = c

	return nil
}
