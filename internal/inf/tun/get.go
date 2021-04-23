package tun

import (
	"fmt"
	"github.com/vishvananda/netlink"
	"golang.org/x/sys/unix"
)

func Get(name string) (*Tun, error) {
	link, err := netlink.LinkByName(name)
	if err != nil {
		return nil, err
	}

	tu, ok := link.(*netlink.Tuntap)
	if !ok {
		return nil, fmt.Errorf("LinkByName type asserting failed")
	}

	res := GetAttr()
	res.persist = !tu.NonPersist
	res.name = tu.Name
	res.multiQueue = tu.Flags&unix.IFF_MULTI_QUEUE == unix.IFF_MULTI_QUEUE
	res.masterId = tu.MasterIndex
	res.mtu = tu.MTU
	res.txQueueLen = tu.TxQLen

	return res, err
}
