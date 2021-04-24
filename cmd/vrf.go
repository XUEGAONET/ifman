package main

import (
	"github.com/sirupsen/logrus"
	"ifman/internal/inf/exist"
	"ifman/internal/inf/vrf"
)

func afVrf(c Interface) error {
	inf := vrf.GetAttr()
	err := inf.SetName(c.Name)
	if err != nil {
		return err
	}
	if v, ok := c.Config["table_id"]; ok {
		inf.SetTableId(uint8(v.(int)))
	}
	if v, ok := c.Config["mac"]; ok {
		err := inf.SetMac(v.(string))
		if err != nil {
			return err
		}
	}

	if exist.IsExisted(c.Name) { // existed
		getInf, err := vrf.Get(c.Name)
		if err != nil {
			return err
		}

		if vrf.Equal(getInf, inf) {
			logrus.Tracef("vrf interface %s check passed", c.Name)
			return nil
		} else {
			logrus.Debugf("vrf interface %s check error: current: %#v, want: %#v", c.Name, getInf, inf)
			err = vrf.Update(inf)
			if err != nil {
				return err
			}
		}
	} else { // not existed
		logrus.Infof("vrf interface %s not exists", c.Name)

		err := vrf.New(inf)
		if err != nil {
			return err
		}
	}

	return nil
}
