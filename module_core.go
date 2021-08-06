package main

import (
	"fmt"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	netlink "github.com/vishvananda/netlink"
	"gopkg.in/yaml.v3"
	"os"
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

	select {
	case <-ticker:
		conf := getCoreConfig()

		processAllLink(conf)
		processAllRpFilter(conf)
		processAllLearning(conf)
		processAllAddr(conf)
	case <-refreshEvent:
		err := refreshCoreConfig()
		if err != nil {
			return errors.WithStack(err)
		}
	}

	return nil
}

func processAllLearning(conf *Config) {
	logrus.Infoln("start to process all learning")

	for i, _ := range conf.Learning {
		c := conf.Learning[i]

		err := UpdateLearning(&c)
		if err != nil {
			logrus.Errorf("update learning mode failed: %s", err.Error())
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
			logrus.Errorf("update rp_filter failed: %s", err.Error())
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
			logrus.Errorf("get addr exist status failed: %s", err.Error())
			continue
		}

		if exist {
			err = UpdateAddr(&c)
		} else {
			err = NewAddr(&c)
		}
		if err != nil {
			logrus.Errorf("update or new addr failed: %s", err.Error())
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
			logrus.Errorf("get info from link failed: %s", err.Error())
			continue
		}

		logrus.Infof("process %s:%s interface", typ, name)

		b, err := yaml.Marshal(c)
		if err != nil {
			logrus.Errorf("marshal one link config failed: %s", err.Error())
			continue
		}

		newLink, err := getLinkFromYaml(typ, b)
		if err != nil {
			logrus.Errorf("get link from yaml failed: %s", err.Error())
			continue
		}

		err = processOneLink(newLink)
		if err != nil {
			logrus.Errorf("process one link failed: %s", err.Error())
			continue
		}
	}
}

func processOneLink(link Link) error {
	_, err := netlink.LinkByName(link.GetBaseAttrs().Name)
	if err != nil {
		// link not exist
		if errors.Is(err, os.ErrNotExist) {
			err = NewLink(link)
			if err != nil {
				return errors.WithStack(err)
			}

			err = UpdateLink(link)
			if err != nil {
				return errors.WithStack(err)
			}
		}

		// unexpected error
		return errors.WithStack(err)
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
		return "", "", fmt.Errorf("get name from link map failed")
	}

	name, ok := tName.(string)
	if !ok {
		return "", "", fmt.Errorf("convert name from interface to string failed")
	}

	tTyp, ok := link["type"]
	if !ok {
		return "", "", fmt.Errorf("get type from link map failed")
	}

	typ, ok := tTyp.(string)
	if !ok {
		return "", "", fmt.Errorf("convert type from interface to string failed")
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
		return nil, fmt.Errorf("unsopprted link type from yaml")
	}

	err = yaml.Unmarshal(b, &link)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return link, nil
}
