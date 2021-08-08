// Copyright 2021 The ifman authors https://github.com/XUEGAONET/ifman
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//    http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

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
