package tun

import (
	"github.com/vishvananda/netlink"
	"golang.org/x/sys/unix"
)

func New(t *Tun) error {
	err := t.check()
	if err != nil {
		return err
	}

	attr := netlink.Tuntap{
		LinkAttrs: netlink.LinkAttrs{
			MTU:         t.mtu,
			TxQLen:      t.txQueueLen,
			Name:        t.name,
			MasterIndex: t.masterId,
		},
		Mode:       unix.IFF_TUN,
		Flags:      netlink.TUNTAP_DEFAULTS,
		NonPersist: !t.persist,
	}
	if t.multiQueue == true {
		attr.Flags = netlink.TUNTAP_MULTI_QUEUE_DEFAULTS
	}

	err = netlink.LinkAdd(&attr)

	return err
}
