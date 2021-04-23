package tun

import "github.com/vishvananda/netlink"

func Update(t *Tun) error {
	err := t.check()
	if err != nil {
		return err
	}

	link, err := netlink.LinkByName(t.name)
	if err != nil {
		return err
	}

	err = netlink.LinkSetMTU(link, t.mtu)
	if err != nil {
		return err
	}

	err = netlink.LinkSetTxQLen(link, t.txQueueLen)
	if err != nil {
		return err
	}

	if t.masterId == 0 {
		err = netlink.LinkSetNoMaster(link)
	} else {
		err = netlink.LinkSetMasterByIndex(link, t.masterId)
	}

	return err
}
