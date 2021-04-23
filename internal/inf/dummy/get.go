package dummy

import (
	"fmt"
	"github.com/vishvananda/netlink"
)

func Get(name string) (*Dummy, error) {
	link, err := netlink.LinkByName(name)
	if err != nil {
		return nil, err
	}

	d, ok := link.(*netlink.Dummy)
	if !ok {
		return nil, fmt.Errorf("LinkByName type asserting failed")
	}

	res := GetAttr()
	res.name = d.Name
	res.mtu = d.MTU
	res.mac = d.HardwareAddr
	res.masterId = d.MasterIndex
	res.txQueueLen = d.TxQLen

	return res, err
}
