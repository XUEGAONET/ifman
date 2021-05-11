package bridge

import (
	"github.com/vishvananda/netlink"
)

func Get(name string) (*Bridge, error) {
	link, err := netlink.LinkByName(name)
	if err != nil {
		return nil, err
	}

	res := GetAttr()
	res.name = link.Attrs().Name
	res.mtu = link.Attrs().MTU
	res.mac = link.Attrs().HardwareAddr
	res.masterId = link.Attrs().MasterIndex
	res.txQueueLen = link.Attrs().TxQLen

	return res, err
}
