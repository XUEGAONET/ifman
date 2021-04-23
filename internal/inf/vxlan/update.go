package vxlan

import "github.com/vishvananda/netlink"

func Update(v *VxLan) error {
	err := v.check()
	if err != nil {
		return err
	}

	link, err := netlink.LinkByName(v.name)
	if err != nil {
		return err
	}

	err = netlink.LinkSetMTU(link, v.mtu)
	if err != nil {
		return err
	}

	err = netlink.LinkSetTxQLen(link, v.txQueueLen)
	if err != nil {
		return err
	}

	err = netlink.LinkSetHardwareAddr(link, v.mac)
	if err != nil {
		return err
	}

	if v.masterId == 0 {
		err = netlink.LinkSetNoMaster(link)
	} else {
		err = netlink.LinkSetMasterByIndex(link, v.masterId)
	}

	return err
}
