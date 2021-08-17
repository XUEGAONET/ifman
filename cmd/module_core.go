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
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	netlink "github.com/vishvananda/netlink"
	"gopkg.in/yaml.v3"
	"strings"
	"time"
)

func startCoreService() error {
	ticker := make(chan struct{}, 1)
	go func() {
		c := getCoreConfig()
		period := time.Duration(c.Common.CheckPeriodSec) * time.Second

		for {
			ticker <- struct{}{}
			time.Sleep(period)
		}
	}()

	for {
		select {
		case <-ticker:
			conf := getCoreConfig()

			processAllLink(conf)
			processAllRpFilter(conf)
			processAllLearning(conf)
			processAllAddr(conf)
		}
	}
}

func processAllLearning(conf *Config) {
	logrus.Infoln("start to process all learning")

	for i, _ := range conf.Learning {
		c := conf.Learning[i]

		err := UpdateLearning(&c)
		if err != nil {
			logrus.Errorf("update learning mode failed: %+v", err)
			continue
		}
	}
}

func processAllRpFilter(conf *Config) {
	logrus.Infoln("start to process all rp_filter")

	for i, _ := range conf.RpFilter {
		c := conf.RpFilter[i]

		err := UpdateRpFilter(&c)
		if err != nil {
			logrus.Errorf("update rp_filter failed: %+v", err)
			continue
		}
	}
}

func processAllAddr(conf *Config) {
	logrus.Infoln("start to process all addr")

	for i, _ := range conf.Addr {
		c := conf.Addr[i]

		exist, err := IsAddrExist(&c)
		if err != nil {
			logrus.Errorf("get addr exist status failed: %+v", err)
			continue
		}

		if exist {
			err = UpdateAddr(&c)
		} else {
			err = NewAddr(&c)
		}
		if err != nil {
			logrus.Errorf("update or new addr failed: %+v", err)
			continue
		}
	}
}

func processAllLink(conf *Config) {
	logrus.Infoln("start to process all link")

	for i, _ := range conf.Interface {
		c := conf.Interface[i]

		l, ok := c.(map[string]interface{})
		if !ok {
			logrus.Errorf("assert link config failed")
			continue
		}

		name, typ, err := getInfoFromLink(l)
		if err != nil {
			logrus.Errorf("get info from link failed: %+v", err)
			continue
		}

		logrus.Infof("process %s:%s interface", typ, name)

		b, err := yaml.Marshal(c)
		if err != nil {
			logrus.Errorf("marshal one link config failed: %+v", err)
			continue
		}

		newLink, err := getLinkFromYaml(typ, b)
		if err != nil {
			logrus.Errorf("get link from yaml failed: %+v", err)
			continue
		}

		err = processOneLink(newLink)
		if err != nil {
			logrus.Errorf("process one link failed: %+v", err)
			continue
		}
	}
}

func processOneLink(link Link) error {
	_, err := netlink.LinkByName(link.GetBaseAttrs().Name)
	if err != nil {
		// link not exist
		// netlink的包太垃圾了，迫使我只能这样写了
		_, ok1 := err.(*netlink.LinkNotFoundError)
		_, ok2 := err.(netlink.LinkNotFoundError)
		if ok1 || ok2 {
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

func getLinkFromYaml(typ string, b []byte) (Link, error) {
	var err error = nil
	var link Link = nil

	switch strings.ToLower(typ) {
	case "bridge":
		link = &Bridge{}
	case "dummy":
		link = &Dummy{}
	case "iptun":
		link = &IPTun{}
	case "unmanaged":
		link = &Unmanaged{}
	case "tun":
		link = &Tun{}
	case "vlan":
		link = &Vlan{}
	case "vrf":
		link = &Vrf{}
	case "vxlan":
		link = &VxLAN{}
	case "wireguard_ptp_server":
		link = &WireGuardPtPServer{}
	case "wireguard_ptp_client":
		link = &WireGuardPtPClient{}
	default:
		return nil, errors.WithStack(fmt.Errorf("unsopprted link type from yaml"))
	}

	err = yaml.Unmarshal(b, link)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return link, nil
}
