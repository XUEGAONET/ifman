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
	"github.com/XUEGAONET/ifman/common"
	pkgaddr "github.com/XUEGAONET/ifman/pkg/addr"
	"github.com/XUEGAONET/ifman/pkg/learning"
	"github.com/XUEGAONET/ifman/pkg/rpf"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/vishvananda/netlink"
	"gopkg.in/yaml.v3"
	"strings"
	"time"
)

var recheckChan chan struct{}

func startService() error {
	recheckChan = make(chan struct{}, 1)
	recheckChan <- struct{}{}

	ticker := time.NewTicker(time.Duration(getGlobalConfig().Common.CheckPeriodSec) * time.Second)

	for {
		select {
		case <-ticker.C:
		case <-recheckChan:
			logrus.Infoln("recheck signal received")
		}

		conf := getGlobalConfig()
		logrus.Tracef("get global config: %+v", conf)

		processAllLink(conf.Interface)
		processAllRpFilter(conf.RpFilter)
		processAllLearning(conf.Learning)
		processAllAddr(conf.Addr)
	}
}

func processAllLearning(conf []common.Learning) {
	var err error

	logrus.Traceln("start to process all learning")

	for _, c := range conf {
		if c.LearningOn {
			err = learning.SetLearningOn(c.Name)
		} else {
			err = learning.SetLearningOff(c.Name)
		}

		if err != nil {
			logrus.Errorf("update learning failed: %v", err)
		}
	}
}

func processAllRpFilter(conf []common.RpFilter) {
	var err error

	logrus.Traceln("start to process all rp_filter")

	for _, c := range conf {
		err = rpf.CheckAndFix(c.Name, rpf.RPFType(c.Mode))
		if err != nil {
			logrus.Errorf("update rp_filter failed: %v", err)
		}
	}
}

func processAllAddr(conf []common.Addr) {
	logrus.Traceln("start to process all addr")

	for _, c := range conf {
		exist, err := pkgaddr.IsAddrExist(c.Name, c.Address)
		if err != nil {
			logrus.Errorf("get addr exist status failed: %v", err)
			continue
		}

		if exist {
			err = pkgaddr.Update(c.Name, c.Address, c.PtpMode, c.PeerPrefix)
		} else {
			err = pkgaddr.New(c.Name, c.Address, c.PtpMode, c.PeerPrefix)
		}
		if err != nil {
			logrus.Errorf("update or new addr failed: %v", err)
			continue
		}
	}
}

func processAllLink(conf []interface{}) {
	logrus.Traceln("start to process all link")

	for _, c := range conf {
		l, ok := c.(map[string]interface{})
		if !ok {
			logrus.Errorf("assert custom link config failed")
			continue
		}

		name, typ, err := getInfoFromLink(l)
		if err != nil {
			logrus.Errorf("get info from link failed: %v", err)
			continue
		}

		logrus.Debugf("start process %s:%s interface", typ, name)

		b, err := yaml.Marshal(c)
		if err != nil {
			logrus.Errorf("marshal one link config failed: %v", err)
			continue
		}

		newLink, err := getLinkFromYaml(typ, b)
		if err != nil {
			logrus.Errorf("get link from yaml failed: %v", err)
			continue
		}

		err = processOneLink(newLink)
		if err != nil {
			logrus.Errorf("process one link failed: %v", err)
			continue
		}
	}
}

func processOneLink(link common.Link) error {
	_, err := netlink.LinkByName(link.GetBaseAttrs().Name)
	if err != nil {
		// link not exist
		// netlink的包太垃圾了，迫使我只能这样写了
		_, ok1 := err.(*netlink.LinkNotFoundError)
		_, ok2 := err.(netlink.LinkNotFoundError)
		if ok1 || ok2 {
			logrus.Debugf("link not exist, go to create")

			err = NewLink(link)
			if err != nil {
				return errors.WithStack(err)
			}

			err = UpdateLink(link)
			if err != nil {
				return errors.WithStack(err)
			}
		} else {
			// unexpected error
			return errors.Wrap(err, "process one link with unexpected error")
		}
	} else { // no error
		logrus.Debugf("update link")

		err = UpdateLink(link)
		if err != nil {
			return errors.WithStack(err)
		}
	}

	return nil
}

func getInfoFromLink(link map[string]interface{}) (string, string, error) {
	tName, ok := link["name"]
	if !ok {
		return "", "", errors.WithStack(fmt.Errorf("get name from link map failed"))
	}

	name, ok := tName.(string)
	if !ok {
		return "", "", errors.WithStack(fmt.Errorf("convert name from interface to string failed"))
	}

	tTyp, ok := link["type"]
	if !ok {
		return "", "", errors.WithStack(fmt.Errorf("get type from link map failed"))
	}

	typ, ok := tTyp.(string)
	if !ok {
		return "", "", errors.WithStack(fmt.Errorf("convert type from interface to string failed"))
	}

	return name, typ, nil
}

// getLinkFromYaml 返回自定义的Link。
// 传入自定义的类型、yaml的局部的内容。
// 增加新类型时需要在此处添加。
func getLinkFromYaml(typ string, b []byte) (common.Link, error) {
	var err error = nil
	var link common.Link = nil

	switch strings.ToLower(typ) {
	case "bridge":
		link = &common.Bridge{}
	case "dummy":
		link = &common.Dummy{}
	case "iptun":
		link = &common.IPTun{}
	case "generic":
		link = &common.Generic{}
	case "tun":
		link = &common.Tun{}
	case "vlan":
		link = &common.Vlan{}
	case "vrf":
		link = &common.Vrf{}
	case "vxlan":
		link = &common.VxLAN{}
	case "wireguard_ptp_server":
		link = &common.WireGuardPtPServer{}
	case "wireguard_ptp_client":
		link = &common.WireGuardPtPClient{}
	case "wireguard_origin":
		link = &common.WireGuardOrigin{}
	default:
		return nil, errors.WithStack(fmt.Errorf("unsopprted link type from yaml"))
	}

	err = yaml.Unmarshal(b, link)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return link, nil
}
