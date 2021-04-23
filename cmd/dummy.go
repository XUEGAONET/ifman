package main

import (
	"ifman/internal/inf/dummy"
	"ifman/internal/inf/exist"
)

func afDummy(c Interface) error {
	inf := dummy.GetAttr()
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
		getInf, err := dummy.Get(c.Name)
		if err != nil {
			return err
		}

		if dummy.Equal(getInf, inf) {
			return nil
		} else {
			err = dummy.Update(inf)
			if err != nil {
				return err
			}
		}
	} else { // not existed
		err := dummy.New(inf)
		if err != nil {
			return err
		}
	}

	return nil
}
