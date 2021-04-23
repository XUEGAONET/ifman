package main

import (
	"fmt"
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

		if generic.Equal(getInf, inf) {
			return nil
		} else {
			err = generic.Update(inf)
			if err != nil {
				return err
			}
		}
	} else { // not exist
		return fmt.Errorf("generic interface %s not exists", c.Name)
	}

	return nil
}
