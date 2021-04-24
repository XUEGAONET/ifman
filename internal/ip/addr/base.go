package addr

import (
	"fmt"
	"ifman/internal/inf/common"
	"net"
)

type Address struct {
	name string
	addr net.IPNet
}

func (a *Address) SetName(s string) error {
	if err := common.ChkName(s); err != nil {
		return err
	}

	a.name = s
	return nil
}

func (a *Address) SetAddr(addr string) error {
	ip, network, err := net.ParseCIDR(addr)
	if err != nil {
		return fmt.Errorf("parse cidr %s failed: %v", addr, err)
	}
	network.IP = ip

	a.addr = *network
	return nil
}

func (a *Address) check() error {
	if a.name == "" {
		return fmt.Errorf("invalid parameter")
	}

	return nil
}

func GetAttr() *Address {
	return &Address{
		name: "",
		addr: net.IPNet{},
	}
}
