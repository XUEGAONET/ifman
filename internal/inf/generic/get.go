package generic

import (
	"github.com/vishvananda/netlink"
)

func Get(name string) (*Generic, error) {
	link, err := netlink.LinkByName(name)
	if err != nil {
		return nil, err
	}

	res := GetAttr()
	res.name = link.Attrs().Name
	res.mtu = link.Attrs().MTU
	res.txQueueLen = link.Attrs().TxQLen
	res.mac = link.Attrs().HardwareAddr
	res.masterId = link.Attrs().MasterIndex

	return res, err
}
