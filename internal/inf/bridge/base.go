package bridge

import (
	"fmt"
	"ifman/internal/inf/common"
	"net"
)

type Bridge struct {
	name       string
	mtu        int
	txQueueLen int
	mac        net.HardwareAddr
	masterId   int
}

func (b *Bridge) SetName(s string) error {
	if err := common.ChkName(s); err != nil {
		return err
	}

	b.name = s
	return nil
}

func (b *Bridge) SetMtu(u uint16) error {
	if err := common.ChkMtu(u); err != nil {
		return err
	}

	b.mtu = int(u)
	return nil
}

func (b *Bridge) SetTxQueueLen(u uint16) {
	b.txQueueLen = int(u)
}

func (b *Bridge) SetMac(s string) error {
	mac, err := common.PassMac(s)
	if err != nil {
		return err
	}

	b.mac = mac
	return nil
}

func (b *Bridge) SetMaster(s string) error {
	id, err := common.PassMaster(s)
	if err != nil {
		return err
	}

	b.masterId = id
	return nil
}

func (b *Bridge) check() error {
	if b.name == "" {
		return fmt.Errorf("invalid parameter")
	}

	return nil
}

func GetAttr() *Bridge {
	return &Bridge{
		name:       "",
		mtu:        0,
		txQueueLen: 0,
		mac:        nil,
		masterId:   0,
	}
}
