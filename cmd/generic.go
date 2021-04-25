package main

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"ifman/internal/inf/exist"
	"ifman/internal/inf/generic"
)

func afGeneric(c Interface) error {
	inf := generic.GetAttr()
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
		getInf, err := generic.Get(c.Name)
		if err != nil {
			return err
		}

		if generic.Equal(inf, getInf) {
			logrus.Tracef("generic interface %s check passed", c.Name)
			return nil
		} else {
			logrus.Debugf("generic interface %s check error: expect: %#v, get: %#v", c.Name, inf, getInf)
			err = generic.Update(inf)
			if err != nil {
				return err
			}
		}
	} else { // not exist
		logrus.Infof("generic interface %s not exists", c.Name)
		return fmt.Errorf("generic interface %s not exists", c.Name)
	}

	return nil
}
