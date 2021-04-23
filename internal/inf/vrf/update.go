package vrf

import "github.com/vishvananda/netlink"

func Update(v *Vrf) error {
	err := v.check()
	if err != nil {
		return err
	}

	link, err := netlink.LinkByName(v.name)
	if err != nil {
		return err
	}

	err = netlink.LinkSetHardwareAddr(link, v.mac)
	return err
}
