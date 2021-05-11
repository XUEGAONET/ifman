package bridge

import (
	"github.com/vishvananda/netlink"
)

func New(d *Bridge) error {
	err := d.check()
	if err != nil {
		return err
	}

	attr := netlink.Bridge{
		LinkAttrs: netlink.LinkAttrs{
			MTU:          d.mtu,
			TxQLen:       d.txQueueLen,
			Name:         d.name,
			HardwareAddr: d.mac,
			MasterIndex:  d.masterId,
		},
		MulticastSnooping: nil,
		HelloTime:         nil,
		VlanFiltering:     nil,
	}

	err = netlink.LinkAdd(&attr)

	return err
}
