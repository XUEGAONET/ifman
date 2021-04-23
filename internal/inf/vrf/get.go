package vrf

import (
	"fmt"
	"github.com/vishvananda/netlink"
)

func Get(name string) (*Vrf, error) {
	link, err := netlink.LinkByName(name)
	if err != nil {
		return nil, err
	}

	vrf, ok := link.(*netlink.Vrf)
	if !ok {
		return nil, fmt.Errorf("LinkByName type asserting failed")
	}

	res := GetAttr()
	res.mac = vrf.HardwareAddr
	res.name = vrf.Name
	res.tableId = vrf.Table

	return res, err
}
