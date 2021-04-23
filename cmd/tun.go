package main

import (
	"ifman/internal/inf/exist"
	"ifman/internal/inf/tun"
)

func afTun(c Interface) error {
	inf := tun.GetAttr()
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
	if v, ok := c.Config["multi_queue"]; ok {
		if v.(bool) {
			inf.SetMultiQueue()
		} else {
			inf.SetOneQueue()
		}
	}
	if v, ok := c.Config["persist"]; ok {
		if v.(bool) {
			inf.SetPersist()
		} else {
			inf.SetNonPersist()
		}
	}

	if exist.IsExisted(c.Name) { // existed
		getInf, err := tun.Get(c.Name)
		if err != nil {
			return err
		}

		if tun.Equal(getInf, inf) {
			return nil
		} else {
			err = tun.Update(inf)
			if err != nil {
				return err
			}
		}
	} else { // not existed
		err := tun.New(inf)
		if err != nil {
			return err
		}
	}

	return nil
}
