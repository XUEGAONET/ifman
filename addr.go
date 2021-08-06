package main

import (
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/vishvananda/netlink"
	"net"
)

type Addr struct {
	Name       string `yaml:"name"`
	Address    string `yaml:"address"`
	PtpMode    bool   `yaml:"ptp_mode"`
	PeerPrefix string `yaml:"peer_prefix"`
}

func NewAddr(addr *Addr) error {
	logrus.Debugf("new addr: %#v", addr)

	link, err := netlink.LinkByName(addr.Name)
	if err != nil {
		return errors.WithStack(err)
	}

	nlAddr := netlink.Addr{
		IPNet:       nil,
		Label:       "",
		Flags:       0,
		Scope:       0,
		Peer:        nil,
		Broadcast:   nil,
		PreferedLft: 0,
		ValidLft:    0,
	}

	// ip address
	ip, ipNet, err := net.ParseCIDR(addr.Address)
	if err != nil {
		return errors.WithStack(err)
	}
	ipNet.IP = ip
	nlAddr.IPNet = ipNet

	// ptp mode
	if addr.PtpMode {
		_, cidr, err := net.ParseCIDR(addr.PeerPrefix)
		if err != nil {
			return errors.WithStack(err)
		}

		nlAddr.Peer = cidr
	}

	err = netlink.AddrAdd(link, &nlAddr)
	if err != nil {
		return errors.WithStack(err)
	}

	logrus.Debugf("new addr succeed")
	return nil
}

// UpdateAddr will check and auto fix addr config
// please check addr whether exist before use this call
func UpdateAddr(addr *Addr) error {
	logrus.Debugf("update addr: %#v", addr)

	link, err := netlink.LinkByName(addr.Name)
	if err != nil {
		return errors.WithStack(err)
	}

	nlAddrs, err := netlink.AddrList(link, netlink.FAMILY_ALL)
	if err != nil {
		return errors.WithStack(err)
	}

	for i, _ := range nlAddrs {
		nlAddr := nlAddrs[i]

		if nlAddr.IPNet.String() == addr.Address {
			if addr.PtpMode {
				// ptp mode but peer not exist
				if nlAddr.Peer == nil {
					err = fixAddr(addr, &nlAddr)
					if err != nil {
						return errors.WithStack(err)
					}

				}

				// ptp mode but peer not equal
				if nlAddr.Peer.String() != addr.PeerPrefix {
					err = fixAddr(addr, &nlAddr)
					if err != nil {
						return errors.WithStack(err)
					}
				}
			}
		}
	}

	logrus.Debugf("update addr succeed")
	return nil
}

func fixAddr(addr *Addr, nlAddr *netlink.Addr) error {
	link, err := netlink.LinkByName(addr.Name)
	if err != nil {
		return errors.WithStack(err)
	}

	err = netlink.AddrDel(link, nlAddr)
	if err != nil {
		return errors.WithStack(err)
	}

	err = NewAddr(addr)
	if err != nil {
		return errors.WithStack(err)
	}

	return nil
}

func IsAddrExist(addr *Addr) (bool, error) {
	res := false

	link, err := netlink.LinkByName(addr.Name)
	if err != nil {
		return res, errors.WithStack(err)
	}

	nlAddrs, err := netlink.AddrList(link, netlink.FAMILY_ALL)
	if err != nil {
		return res, errors.WithStack(err)
	}

	for i, _ := range nlAddrs {
		nlAddr := nlAddrs[i]

		if nlAddr.IPNet.String() == addr.Address {
			res = true
		}
	}

	return res, nil
}
