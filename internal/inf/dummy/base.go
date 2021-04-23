package dummy

import (
	"fmt"
	"ifman/internal/inf/common"
	"net"
)

type Dummy struct {
	name       string
	mtu        int
	txQueueLen int
	mac        net.HardwareAddr
	masterId   int
}

func (d *Dummy) SetName(s string) error {
	if err := common.ChkName(s); err != nil {
		return err
	}

	d.name = s
	return nil
}

func (d *Dummy) SetMtu(u uint16) error {
	if err := common.ChkMtu(u); err != nil {
		return err
	}

	d.mtu = int(u)
	return nil
}

func (d *Dummy) SetTxQueueLen(u uint16) {
	d.txQueueLen = int(u)
}

func (d *Dummy) SetMac(s string) error {
	mac, err := common.PassMac(s)
	if err != nil {
		return err
	}

	d.mac = mac
	return nil
}

func (d *Dummy) SetMaster(s string) error {
	id, err := common.PassMaster(s)
	if err != nil {
		return err
	}

	d.masterId = id
	return nil
}

func (d *Dummy) check() error {
	if d.name == "" {
		return fmt.Errorf("invalid parameter")
	}

	return nil
}

func GetAttr() *Dummy {
	return &Dummy{
		name:       "",
		mtu:        0,
		txQueueLen: 0,
		mac:        nil,
		masterId:   0,
	}
}
