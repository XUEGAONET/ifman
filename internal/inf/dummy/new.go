package dummy

import (
	"github.com/vishvananda/netlink"
)

func New(d *Dummy) error {
	err := d.check()
	if err != nil {
		return err
	}

	attr := netlink.Dummy{
		LinkAttrs: netlink.LinkAttrs{
			MTU:          d.mtu,
			TxQLen:       d.txQueueLen,
			Name:         d.name,
			HardwareAddr: d.mac,
			MasterIndex:  d.masterId,
		},
	}

	err = netlink.LinkAdd(&attr)

	return err
}
