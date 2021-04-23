package vrf

import (
	"github.com/vishvananda/netlink"
)

func New(v *Vrf) error {
	err := v.check()
	if err != nil {
		return err
	}

	attr := netlink.Vrf{
		LinkAttrs: netlink.LinkAttrs{
			Name:         v.name,
			HardwareAddr: v.mac,
		},
		Table: v.tableId,
	}

	err = netlink.LinkAdd(&attr)

	return err
}
