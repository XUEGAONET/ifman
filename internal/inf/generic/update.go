package generic

import "github.com/vishvananda/netlink"

func Update(g *Generic) error {
	err := g.check()
	if err != nil {
		return err
	}

	link, err := netlink.LinkByName(g.name)
	if err != nil {
		return err
	}

	err = netlink.LinkSetMTU(link, g.mtu)
	if err != nil {
		return err
	}

	err = netlink.LinkSetTxQLen(link, g.txQueueLen)
	if err != nil {
		return err
	}

	if g.mac != nil {
		err = netlink.LinkSetHardwareAddr(link, g.mac)
		if err != nil {
			return err
		}
	}

	if g.masterId == 0 {
		err = netlink.LinkSetNoMaster(link)
	} else {
		err = netlink.LinkSetMasterByIndex(link, g.masterId)
	}

	return err
}
