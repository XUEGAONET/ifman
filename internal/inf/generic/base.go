package generic

import (
	"fmt"
	"ifman/internal/inf/common"
	"net"
)

type Generic struct {
	name       string
	mtu        int
	txQueueLen int
	mac        net.HardwareAddr
	masterId   int
}

func (g *Generic) SetName(s string) error {
	if err := common.ChkName(s); err != nil {
		return err
	}

	g.name = s
	return nil
}

func (g *Generic) SetMtu(u uint16) error {
	if err := common.ChkMtu(u); err != nil {
		return err
	}

	g.mtu = int(u)
	return nil
}

func (g *Generic) SetTxQueueLen(u uint16) {
	g.txQueueLen = int(u)
}

func (g *Generic) SetMac(s string) error {
	mac, err := common.PassMac(s)
	if err != nil {
		return err
	}

	g.mac = mac
	return nil
}

func (g *Generic) SetMaster(s string) error {
	id, err := common.PassMaster(s)
	if err != nil {
		return err
	}

	g.masterId = id
	return nil
}

func GetAttr() *Generic {
	return &Generic{
		name:       "",
		mtu:        0,
		txQueueLen: 0,
		mac:        nil,
		masterId:   0,
	}
}

func (g *Generic) check() error {
	if g.name == "" {
		return fmt.Errorf("invalid parameter")
	}

	return nil
}
