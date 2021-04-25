package main

import (
	"github.com/sirupsen/logrus"
	"ifman/internal/inf/exist"
	"ifman/internal/inf/vxlan"
)

func afVxLan(c Interface) error {
	inf := vxlan.GetAttr()
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
	if v, ok := c.Config["vni"]; ok {
		err := inf.SetVni(uint32(v.(int)))
		if err != nil {
			return err
		}
	}
	if v, ok := c.Config["src"]; ok {
		err := inf.SetSrc(v.(string))
		if err != nil {
			return err
		}
	}
	if v, ok := c.Config["dst"]; ok {
		err := inf.SetDst(v.(string))
		if err != nil {
			return err
		}
	}
	if v, ok := c.Config["ttl"]; ok {
		inf.SetTtl(uint8(v.(int)))
	}
	if v, ok := c.Config["tos"]; ok {
		inf.SetTos(uint8(v.(int)))
	}
	if v, ok := c.Config["learning"]; ok {
		if v.(bool) {
			inf.SetLearning()
		} else {
			inf.SetNoLearning()
		}
	}
	if v, ok := c.Config["src_port_low"]; ok {
		inf.SetSrcPortLow(uint16(v.(int)))
	}
	if v, ok := c.Config["src_port_high"]; ok {
		inf.SetSrcPortHigh(uint16(v.(int)))
	}
	if v, ok := c.Config["port"]; ok {
		inf.SetPort(uint16(v.(int)))
	}

	if exist.IsExisted(c.Name) { // existed
		getInf, err := vxlan.Get(c.Name)
		if err != nil {
			return err
		}

		if vxlan.Equal(inf, getInf) {
			logrus.Tracef("vxlan interface %s check passed", c.Name)
			return nil
		} else {
			logrus.Debugf("vxlan interface %s check error: expect: %#v, get: %#v", c.Name, inf, getInf)
			err = vxlan.Update(inf)
			if err != nil {
				return err
			}
		}
	} else { // not existed
		logrus.Infof("vxlan interface %s not exists", c.Name)

		err := vxlan.New(inf)
		if err != nil {
			return err
		}
	}

	return nil
}
