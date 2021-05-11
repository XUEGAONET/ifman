package main

import (
	"github.com/sirupsen/logrus"
	"ifman/internal/inf/bridge"
	"ifman/internal/inf/exist"
)

func afBridge(c Interface) error {
	inf := bridge.GetAttr()
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
	if v, ok := c.Config["mac"]; ok {
		err := inf.SetMac(v.(string))
		if err != nil {
			return err
		}
	}
	if v, ok := c.Config["master"]; ok {
		err := inf.SetMaster(v.(string))
		if err != nil {
			return err
		}
	}

	if exist.IsExisted(c.Name) { // existed
		getInf, err := bridge.Get(c.Name)
		if err != nil {
			return err
		}

		if bridge.Equal(inf, getInf) {
			logrus.Tracef("bridge interface %s check passed", c.Name)
			return nil
		} else {
			logrus.Debugf("bridge interface %s check error: expect: %#v, get: %#v", c.Name, inf, getInf)
			err = bridge.Update(inf)
			if err != nil {
				return err
			}
		}
	} else { // not existed
		logrus.Infof("bridge interface %s not exists", c.Name)

		err := bridge.New(inf)
		if err != nil {
			return err
		}
	}

	return nil
}
