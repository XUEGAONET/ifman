package vxlan

import (
	"fmt"
	"github.com/vishvananda/netlink"
)

func Get(name string) (*VxLan, error) {
	link, err := netlink.LinkByName(name)
	if err != nil {
		return nil, err
	}

	v, ok := link.(*netlink.Vxlan)
	if !ok {
		return nil, fmt.Errorf("LinkByName type asserting failed")
	}

	res := GetAttr()
	res.name = v.Name
	res.mtu = v.MTU
	res.txQueueLen = v.TxQLen
	res.masterId = v.MasterIndex
	res.vni = v.VxlanId
	res.tos = v.TOS
	res.ttl = v.TTL
	res.src = v.SrcAddr
	res.dst = v.Group
	res.port = v.Port
	res.srcPortLow = v.PortLow
	res.srcPortHigh = v.PortHigh
	res.learning = v.Learning
	res.mac = v.HardwareAddr

	return res, err
}
