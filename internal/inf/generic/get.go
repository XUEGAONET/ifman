package generic

import (
	"fmt"
	"github.com/vishvananda/netlink"
)

func Get(name string) (*Generic, error) {
	link, err := netlink.LinkByName(name)
	if err != nil {
		return nil, err
	}

	d, ok := link.(netlink.Link)
	if !ok {
		return nil, fmt.Errorf("LinkByName type asserting failed")
	}

	res := GetAttr()
	res.name = d.Attrs().Name
	res.mtu = d.Attrs().MTU
	res.txQueueLen = d.Attrs().TxQLen
	res.mac = d.Attrs().HardwareAddr
	res.masterId = d.Attrs().MasterIndex

	return res, err
}
