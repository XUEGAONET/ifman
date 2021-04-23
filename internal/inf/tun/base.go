package tun

import (
	"fmt"
	"ifman/internal/inf/common"
)

type Tun struct {
	name       string
	txQueueLen int
	mtu        int
	masterId   int
	multiQueue bool
	persist    bool
}

func (t *Tun) SetName(s string) error {
	if err := common.ChkName(s); err != nil {
		return err
	}

	t.name = s
	return nil
}

func (t *Tun) SetTxQueueLen(u uint16) {
	t.txQueueLen = int(u)
}

func (t *Tun) SetMtu(u uint16) error {
	if err := common.ChkMtu(u); err != nil {
		return err
	}

	t.mtu = int(u)
	return nil
}

func (t *Tun) SetMaster(s string) error {
	id, err := common.PassMaster(s)
	if err != nil {
		return err
	}

	t.masterId = id
	return nil
}

func (t *Tun) SetOneQueue() {
	t.multiQueue = false
}

func (t *Tun) SetMultiQueue() {
	t.multiQueue = true
}

func (t *Tun) SetPersist() {
	t.persist = true
}

func (t *Tun) SetNonPersist() {
	t.persist = false
}

func (t *Tun) check() error {
	if t.name == "" {
		return fmt.Errorf("invalid parameter")
	}

	return nil
}

func GetAttr() *Tun {
	return &Tun{
		name:       "",
		txQueueLen: 0,
		mtu:        0,
		masterId:   0,
		multiQueue: false,
		persist:    false,
	}
}
