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

package addr

import (
	"github.com/pkg/errors"
	"github.com/vishvananda/netlink"
	"net"
)

// New 添加新地址。
// 当地址存在时，添加操作会失败。
func New(name string, addr string, ptp bool, peer string) error {
	link, err := netlink.LinkByName(name)
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
	ip, ipNet, err := net.ParseCIDR(addr)
	if err != nil {
		return errors.WithStack(err)
	}
	ipNet.IP = ip
	nlAddr.IPNet = ipNet

	// ptp mode
	if ptp {
		_, cidr, err := net.ParseCIDR(peer)
		if err != nil {
			return errors.WithStack(err)
		}

		nlAddr.Peer = cidr
	}

	err = netlink.AddrAdd(link, &nlAddr)
	if err != nil {
		return errors.WithStack(err)
	}

	return nil
}

// Update 会检查并自动修复地址。
// 仅仅会检查已经存在的地址的配置是否为预期，对于不存在的地址，会直接忽略。需要自行检查是否存在，并补齐。
func Update(name string, addr string, ptp bool, peer string) error {
	link, err := netlink.LinkByName(name)
	if err != nil {
		return errors.WithStack(err)
	}

	nlAddrs, err := netlink.AddrList(link, netlink.FAMILY_ALL)
	if err != nil {
		return errors.WithStack(err)
	}

	for _, nlAddr := range nlAddrs {
		// 找到指定的地址进行检查，没有的话就结束
		if nlAddr.IPNet.String() == addr {
			if ptp {
				// ptp mode but peer not exist
				if nlAddr.Peer == nil {
					if err = deleteAddr(name, &nlAddr); err != nil {
						return errors.WithStack(err)
					}

					if err = New(name, addr, ptp, peer); err != nil {
						return errors.WithStack(err)
					}
				}

				// ptp mode but peer not equal
				if nlAddr.Peer.String() != peer {
					if err = deleteAddr(name, &nlAddr); err != nil {
						return errors.WithStack(err)
					}

					if err = New(name, addr, ptp, peer); err != nil {
						return errors.WithStack(err)
					}
				}
			}
		}
	}

	return nil
}

func deleteAddr(name string, nlAddr *netlink.Addr) error {
	link, err := netlink.LinkByName(name)
	if err != nil {
		return errors.WithStack(err)
	}

	err = netlink.AddrDel(link, nlAddr)
	if err != nil {
		return errors.WithStack(err)
	}

	return nil
}

func IsAddrExist(name string, addr string) (bool, error) {
	res := false

	link, err := netlink.LinkByName(name)
	if err != nil {
		return res, errors.WithStack(err)
	}

	nlAddrs, err := netlink.AddrList(link, netlink.FAMILY_ALL)
	if err != nil {
		return res, errors.WithStack(err)
	}

	for _, nlAddr := range nlAddrs {
		if nlAddr.IPNet.String() == addr {
			res = true
		}
	}

	return res, nil
}
