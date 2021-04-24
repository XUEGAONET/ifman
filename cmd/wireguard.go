package main

import (
	"github.com/sirupsen/logrus"
	"ifman/internal/inf/exist"
	"ifman/internal/inf/wireguard"
)

func afWireGuard(c Interface) error {
	inf := wireguard.GetAttr()
	err := inf.SetName(c.Name)
	if err != nil {
		return err
	}
	if v, ok := c.Config["mtu"]; ok {
		err := inf.SetMtu(uint16(v.(int)))
		if err != nil {
			return err
		}
	}
	if v, ok := c.Config["tx_queue_len"]; ok {
		inf.SetTxQueueLen(uint16(v.(int)))
	}
	if v, ok := c.Config["master"]; ok {
		err := inf.SetMaster(v.(string))
		if err != nil {
			return err
		}
	}
	if v, ok := c.Config["key"]; ok {
		err := inf.SetKey(keyHash(v.(string)))
		if err != nil {
			return err
		}
	}
	if v, ok := c.Config["listen_port"]; ok {
		err := inf.SetListenPort(uint16(v.(int)))
		if err != nil {
			return err
		}
	}
	if v, ok := c.Config["endpoint"]; ok {
		err := inf.SetEndpoint(v.(string))
		if err != nil {
			return err
		}
	}
	if v, ok := c.Config["handshake"]; ok {
		inf.SetHandshakeIntervalSec(uint16(v.(int)))
	}

	if exist.IsExisted(c.Name) { // existed
		getInf, err := wireguard.Get(c.Name)
		if err != nil {
			return err
		}

		if wireguard.Equal(getInf, inf) {
			logrus.Tracef("wireguard interface %s check passed", c.Name)
			return nil
		} else {
			logrus.Debugf("wireguard interface %s check error: current: %#v, want: %#v", c.Name, getInf, inf)
			err = wireguard.Update(inf)
			if err != nil {
				return err
			}
		}
	} else { // not existed
		logrus.Infof("wireguard interface %s not exists", c.Name)

		err := wireguard.New(inf)
		if err != nil {
			return err
		}
	}

	return nil
}
