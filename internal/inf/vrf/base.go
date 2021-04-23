package vrf

import (
	"fmt"
	"ifman/internal/inf/common"
	"net"
)

type Vrf struct {
	name    string
	tableId uint32
	mac     net.HardwareAddr
}

func (v *Vrf) SetName(s string) error {
	if err := common.ChkName(s); err != nil {
		return err
	}

	v.name = s
	return nil
}

func (v *Vrf) SetTableId(u uint8) {
	v.tableId = uint32(u)
}

func (v *Vrf) SetMac(s string) error {
	mac, err := common.PassMac(s)
	if err != nil {
		return err
	}

	v.mac = mac
	return nil
}

func (v *Vrf) check() error {
	if v.name == "" {
		return fmt.Errorf("invalid parameter")
	}

	return nil
}

func GetAttr() *Vrf {
	return &Vrf{
		name:    "",
		tableId: 0,
		mac:     nil,
	}
}
