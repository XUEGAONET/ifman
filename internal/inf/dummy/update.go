package dummy

import "github.com/vishvananda/netlink"

func Update(d *Dummy) error {
	err := d.check()
	if err != nil {
		return err
	}

	link, err := netlink.LinkByName(d.name)
	if err != nil {
		return err
	}

	err = netlink.LinkSetMTU(link, d.mtu)
	if err != nil {
		return err
	}

	err = netlink.LinkSetTxQLen(link, d.txQueueLen)
	if err != nil {
		return err
	}

	err = netlink.LinkSetHardwareAddr(link, d.mac)
	if err != nil {
		return err
	}

	if d.masterId == 0 {
		err = netlink.LinkSetNoMaster(link)
	} else {
		err = netlink.LinkSetMasterByIndex(link, d.masterId)
	}

	return err
}
