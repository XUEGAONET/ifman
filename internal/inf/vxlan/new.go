package vxlan

import "github.com/vishvananda/netlink"

func New(v *VxLan) error {
	err := v.check()
	if err != nil {
		return err
	}

	attr := netlink.Vxlan{
		LinkAttrs: netlink.LinkAttrs{
			MTU:          v.mtu,
			TxQLen:       v.txQueueLen,
			Name:         v.name,
			HardwareAddr: v.mac,
			MasterIndex:  v.masterId,
		},
		VxlanId:  v.vni,
		SrcAddr:  v.src,
		Group:    v.dst,
		TTL:      v.ttl,
		TOS:      v.tos,
		Learning: v.learning,
		Port:     v.dstPort,
		PortLow:  v.srcPortLow,
		PortHigh: v.srcPortHigh,
	}
	err = netlink.LinkAdd(&attr)
	return err
}
