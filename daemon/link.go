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
	"fmt"
	"github.com/XUEGAONET/ifman/common"
	"github.com/XUEGAONET/ifman/pkg/wgkey"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/vishvananda/netlink"
	"golang.org/x/sys/unix"
	"golang.zx2c4.com/wireguard/wgctrl"
	"golang.zx2c4.com/wireguard/wgctrl/wgtypes"
	"net"
	"strings"
	"time"
)

func NewLink(customLink common.Link) error {
	logrus.Debugf("new customLink: %#v", customLink)

	var err error = nil
	var nlLink netlink.Link = nil

	switch c := customLink.(type) {
	case *common.Dummy:
		nlLink = &netlink.Dummy{
			LinkAttrs: netlink.LinkAttrs{},
		}
	case *common.Bridge:
		nlLink = &netlink.Bridge{
			LinkAttrs:         netlink.LinkAttrs{},
			MulticastSnooping: &c.MulticastSnoopingOn,
			HelloTime:         nil,
			VlanFiltering:     &c.VlanFilteringOn,
		}
	case *common.IPTun:
		nlLink = &netlink.Iptun{
			LinkAttrs:  netlink.LinkAttrs{},
			Ttl:        c.Ttl,
			Tos:        c.Tos,
			PMtuDisc:   0,
			Link:       0,
			Local:      nil,
			Remote:     nil,
			EncapSport: 0,
			EncapDport: 0,
			EncapType:  0,
			EncapFlags: 0,
			FlowBased:  false,
		}

		if c.LocalIP != "" {
			ip := net.ParseIP(c.LocalIP)
			if ip == nil {
				return errors.WithStack(fmt.Errorf("parse ip address failed"))
			}
			nlLink.(*netlink.Iptun).Local = ip
		}

		if c.RemoteIP != "" {
			ip := net.ParseIP(c.RemoteIP)
			if ip == nil {
				return errors.WithStack(fmt.Errorf("parse ip address failed"))
			}
			nlLink.(*netlink.Iptun).Remote = ip
		}
	case *common.Tun:
		nlLink = &netlink.Tuntap{
			LinkAttrs:  netlink.LinkAttrs{},
			Mode:       unix.IFF_TUN,
			Flags:      netlink.TUNTAP_DEFAULTS,
			NonPersist: !c.PersistOn,
			Queues:     int(c.Queues),
			Fds:        nil,
			Owner:      0,
			Group:      0,
		}

		if c.MultiQueueOn {
			nlLink.(*netlink.Tuntap).Flags = netlink.TUNTAP_MULTI_QUEUE_DEFAULTS
		}
	case *common.Vlan:
		nlLink = &netlink.Vlan{
			LinkAttrs:    netlink.LinkAttrs{},
			VlanId:       int(c.VlanId),
			VlanProtocol: netlink.VLAN_PROTOCOL_8021Q,
		}

		if c.BindLink != "" {
			bind, err := netlink.LinkByName(c.BindLink)
			if err != nil {
				return errors.WithStack(err)
			}
			nlLink.Attrs().ParentIndex = bind.Attrs().Index
		}

		if c.StackingOn {
			nlLink.(*netlink.Vlan).VlanProtocol = netlink.VLAN_PROTOCOL_8021AD
		}
	case *common.Vrf:
		nlLink = &netlink.Vrf{
			LinkAttrs: netlink.LinkAttrs{},
			Table:     c.TableId,
		}
	case *common.VxLAN:
		nlLink = &netlink.Vxlan{
			LinkAttrs:      netlink.LinkAttrs{},
			VxlanId:        int(c.Vni),
			VtepDevIndex:   0,
			SrcAddr:        nil,
			Group:          nil,
			TTL:            int(c.Ttl),
			TOS:            int(c.Tos),
			Learning:       c.LearningOn,
			Proxy:          false,
			RSC:            false,
			L2miss:         false,
			L3miss:         false,
			UDPCSum:        c.Checksum,
			UDP6ZeroCSumTx: false,
			UDP6ZeroCSumRx: false,
			NoAge:          false,
			GBP:            false,
			FlowBased:      false,
			Age:            0,
			Limit:          0,
			Port:           int(c.Port),
			PortLow:        int(c.SrcPortLow),
			PortHigh:       int(c.SrcPortHigh),
		}

		if c.SrcIp != "" {
			ip := net.ParseIP(c.SrcIp)
			if ip == nil {
				return errors.WithStack(errors.New("parse ip address failed"))
			}
			nlLink.(*netlink.Vxlan).SrcAddr = ip
		}

		if c.DstIP != "" {
			ip := net.ParseIP(c.DstIP)
			if ip == nil {
				return errors.WithStack(errors.New("parse ip address failed"))
			}
			nlLink.(*netlink.Vxlan).Group = ip
		}

		if c.VtepName != "" {
			id, err := netlink.LinkByName(c.VtepName)
			if err != nil {
				return errors.WithStack(err)
			}
			nlLink.(*netlink.Vxlan).VtepDevIndex = id.Attrs().Index
		}
	case *common.WireGuardPtPServer:
		nlLink = &common.WireGuardLink{
			LinkAttrs: netlink.LinkAttrs{},
		}
	case *common.WireGuardPtPClient:
		nlLink = &common.WireGuardLink{
			LinkAttrs: netlink.LinkAttrs{},
		}
	case *common.WireGuardOrigin:
		nlLink = &common.WireGuardLink{
			LinkAttrs: netlink.LinkAttrs{},
		}
	default:
		return errors.WithStack(fmt.Errorf("unsupport high level customLink type"))
	}

	err = newBase(customLink, nlLink)
	if err != nil {
		return errors.WithStack(err)
	}

	err = netlink.LinkAdd(nlLink)
	if err != nil {
		return errors.WithStack(err)
	}

	return nil
}

func newBase(customLink common.Link, nlLink netlink.Link) error {
	// high level attrs
	high := customLink.GetBaseAttrs()
	// low level attrs
	low := nlLink.Attrs()

	switch common.GetLinkType(customLink) {
	case common.LinkTypeLayer2:
		if high.Mac != "" {
			mac, err := net.ParseMAC(strings.ToLower(high.Mac))
			if err != nil {
				return errors.WithStack(err)
			}
			low.HardwareAddr = mac
		}

		fallthrough
	case common.LinkTypeLayer3:

		fallthrough
	default:
		if high.LinkUp {
			low.Flags = net.FlagUp
		}

		if high.MasterName != "" {
			id, err := netlink.LinkByName(high.MasterName)
			if err != nil {
				return errors.WithStack(err)
			}
			low.MasterIndex = id.Attrs().Index
		}

		low.MTU = int(high.Mtu)

		low.Name = high.Name

		low.TxQLen = int(high.TxQueueLen)
	}

	return nil
}

func UpdateLink(customLink common.Link) error {
	logrus.Debugf("update customLink: %#v", customLink)

	var err error = nil
	var nlLink netlink.Link = nil

	nlLink, err = netlink.LinkByName(customLink.GetBaseAttrs().Name)
	if err != nil {
		return errors.WithStack(err)
	}

	switch c := customLink.(type) {
	case *common.Dummy:
	case *common.Bridge:
		nlBr, ok := nlLink.(*netlink.Bridge)
		if !ok {
			err = rebuildLink(customLink, nlLink)
			if err != nil {
				return errors.WithStack(err)
			}
		}

		if *nlBr.MulticastSnooping != c.MulticastSnoopingOn {
			err = rebuildLink(customLink, nlLink)
			if err != nil {
				return errors.WithStack(err)
			}
		}

		if *nlBr.VlanFiltering != c.VlanFilteringOn {
			err = rebuildLink(customLink, nlLink)
			if err != nil {
				return errors.WithStack(err)
			}
		}
	case *common.IPTun:
	case *common.Generic:
	case *common.Tun:
	case *common.Vlan:
		if c.BindLink != "" {
			bind, err := netlink.LinkByName(c.BindLink)
			if err != nil {
				return errors.WithStack(err)
			}

			if bind.Attrs().Index != nlLink.Attrs().ParentIndex {
				err = rebuildLink(customLink, nlLink)
				if err != nil {
					return errors.WithStack(err)
				}
			}
		}
	case *common.Vrf:
	case *common.VxLAN:
	case *common.WireGuardPtPServer:
		if c.KeyChain != "" {
			pri, pub, err := wgkey.DecodeKeyChain(c.KeyChain)
			if err != nil {
				return errors.WithStack(err)
			}
			c.Private = pri
			c.PeerPublic = pub
		}

		changed, err := isWireGuardChanged(c)
		if err != nil {
			return errors.WithStack(err)
		}

		if changed {
			err = rebuildWireGuard(c)
			if err != nil {
				return errors.WithStack(err)
			}
		}
	case *common.WireGuardPtPClient:
		if c.KeyChain != "" {
			pri, pub, err := wgkey.DecodeKeyChain(c.KeyChain)
			if err != nil {
				return errors.WithStack(err)
			}
			c.Private = pri
			c.PeerPublic = pub
		}

		changed, err := isWireGuardChanged(c)
		if err != nil {
			return errors.WithStack(err)
		}

		if changed {
			err = rebuildWireGuard(c)
			if err != nil {
				return errors.WithStack(err)
			}
		}
	case *common.WireGuardOrigin:

	default:
		return errors.WithStack(fmt.Errorf("unsupport high level customLink type"))
	}

	err = updateBase(customLink, nlLink)
	if err != nil {
		return errors.WithStack(err)
	}

	return nil
}

func rebuildLink(customLink common.Link, nlLink netlink.Link) error {
	logrus.Warnf("current customLink %s is not as expected, go to rebuild", customLink.GetBaseAttrs().Name)
	logrus.Debugf("current nLink: %#v, want customLink: %#v", nlLink, customLink)

	err := netlink.LinkDel(nlLink)
	if err != nil {
		return errors.WithStack(err)
	}

	err = NewLink(customLink)
	if err != nil {
		return errors.WithStack(err)
	}

	return nil
}

func isWireGuardChanged(wantCustomLink common.Link) (bool, error) {
	cli, err := wgctrl.New()
	if err != nil {
		return false, errors.WithStack(err)
	}
	defer cli.Close()

	wgCurrent, err := cli.Device(wantCustomLink.GetBaseAttrs().Name)
	if err != nil {
		return false, errors.WithStack(err)
	}

	switch w := wantCustomLink.(type) {
	case *common.WireGuardPtPServer:
		if wgCurrent.ListenPort != int(w.ListenPort) ||
			wgCurrent.PrivateKey.String() != w.Private {
			return true, nil
		}
	case *common.WireGuardPtPClient:
		if wgCurrent.PrivateKey.String() != w.Private {
			return true, nil
		}
	case *common.WireGuardOrigin:
		if (w.ListenPort != 0 && wgCurrent.ListenPort != int(w.ListenPort)) ||
			wgCurrent.PrivateKey.String() != w.Private {
			return true, nil
		}

		if len(wgCurrent.Peers) != len(w.Peers) {
			return true, nil
		}

		pubKeys := make(map[string]struct{}, len(w.Peers))
		for i, _ := range w.Peers {
			pubKeys[w.Peers[i].PeerPublic] = struct{}{}
		}

		for i, _ := range wgCurrent.Peers {
			_, ok := pubKeys[wgCurrent.Peers[i].PublicKey.String()]
			if !ok {
				return true, nil
			}
		}
	default:
		return false, errors.WithStack(fmt.Errorf("unsupport check and rebuild wireguard type"))
	}

	return false, nil
}

func rebuildWireGuard(wantCustomLink common.Link) error {
	cli, err := wgctrl.New()
	if err != nil {
		return errors.WithStack(err)
	}
	defer cli.Close()

	conf := &wgtypes.Config{
		PrivateKey:   nil,
		ListenPort:   nil,
		FirewallMark: nil,
		ReplacePeers: true,
		Peers:        nil,
	}

	switch w := wantCustomLink.(type) {
	case *common.WireGuardPtPServer:
		_, allow, _ := net.ParseCIDR("0.0.0.0/0")
		pub, err := wgtypes.ParseKey(w.PeerPublic)
		if err != nil {
			return errors.WithStack(err)
		}
		pri, err := wgtypes.ParseKey(w.Private)
		if err != nil {
			return errors.WithStack(err)
		}

		conf.PrivateKey = &pri
		port := int(w.ListenPort)
		conf.ListenPort = &port

		peer := wgtypes.PeerConfig{
			PublicKey:                   pub,
			Remove:                      false,
			UpdateOnly:                  false,
			PresharedKey:                nil,
			Endpoint:                    nil,
			PersistentKeepaliveInterval: nil,
			ReplaceAllowedIPs:           true,
			AllowedIPs:                  []net.IPNet{*allow},
		}
		conf.Peers = append(conf.Peers, peer)
	case *common.WireGuardPtPClient:
		_, allow, _ := net.ParseCIDR("0.0.0.0/0")
		pub, err := wgtypes.ParseKey(w.PeerPublic)
		if err != nil {
			return errors.WithStack(err)
		}
		pri, err := wgtypes.ParseKey(w.Private)
		if err != nil {
			return errors.WithStack(err)
		}
		ep, err := net.ResolveUDPAddr("udp", w.Endpoint)
		if err != nil {
			return errors.WithStack(err)
		}
		heart := time.Duration(w.HeartbeatInterval) * time.Second

		conf.PrivateKey = &pri

		peer := wgtypes.PeerConfig{
			PublicKey:                   pub,
			Remove:                      false,
			UpdateOnly:                  false,
			PresharedKey:                nil,
			Endpoint:                    ep,
			PersistentKeepaliveInterval: &heart,
			ReplaceAllowedIPs:           true,
			AllowedIPs:                  []net.IPNet{*allow},
		}
		conf.Peers = append(conf.Peers, peer)
	case *common.WireGuardOrigin:
		pri, err := wgtypes.ParseKey(w.Private)
		if err != nil {
			return errors.WithStack(err)
		}
		conf.PrivateKey = &pri

		if w.ListenPort != 0 {
			port := int(w.ListenPort)
			conf.ListenPort = &port
		}

		for _, customPeer := range w.Peers {
			var peer wgtypes.PeerConfig

			pub, err := wgtypes.ParseKey(customPeer.PeerPublic)
			if err != nil {
				return errors.WithStack(err)
			}
			peer.PublicKey = pub

			if customPeer.HeartbeatInterval != 0 {
				heart := time.Duration(customPeer.HeartbeatInterval) * time.Second
				peer.PersistentKeepaliveInterval = &heart
			}

			if customPeer.Endpoint != "" {
				ep, err := net.ResolveUDPAddr("udp", customPeer.Endpoint)
				if err != nil {
					return errors.WithStack(err)
				}
				peer.Endpoint = ep
			}

			for _, ip := range customPeer.AllowedCIDR {
				_, allow, err := net.ParseCIDR(ip)
				if err != nil {
					return fmt.Errorf("parse wireguard %s allowip %s failed: %v", w.Name, ip, err)
				}
				peer.AllowedIPs = append(peer.AllowedIPs, *allow)
			}

			conf.Peers = append(conf.Peers, peer)
		}
	default:
		return errors.WithStack(fmt.Errorf("unsupport check and rebuild wireguard type"))
	}

	err = cli.ConfigureDevice(wantCustomLink.GetBaseAttrs().Name, *conf)
	if err != nil {
		return errors.WithStack(err)
	}
	return nil
}

func updateBase(customLink common.Link, nlLink netlink.Link) error {
	// high level attrs
	high := customLink.GetBaseAttrs()
	// low level attrs
	low := nlLink.Attrs()

	var err error = nil

	switch common.GetLinkType(customLink) {
	case common.LinkTypeLayer2:
		if high.Mac != "" {
			if strings.ToLower(low.HardwareAddr.String()) != strings.ToLower(high.Mac) {
				newMac, err := net.ParseMAC(strings.ToLower(high.Mac))
				if err != nil {
					return errors.WithStack(err)
				}
				err = netlink.LinkSetHardwareAddr(nlLink, newMac)
				if err != nil {
					return errors.WithStack(err)
				}
			}
		}

		fallthrough
	case common.LinkTypeLayer3:

		fallthrough
	default:
		// link status
		linkUp := low.Flags&net.FlagUp == net.FlagUp
		if linkUp != high.LinkUp {
			if high.LinkUp {
				err = netlink.LinkSetUp(nlLink)
			} else {
				err = netlink.LinkSetDown(nlLink)
			}

			if err != nil {
				return errors.WithStack(err)
			}
		}

		// link master
		if high.MasterName == "" {
			if low.MasterIndex != 0 {
				err = netlink.LinkSetNoMaster(nlLink)
				if err != nil {
					return errors.WithStack(err)
				}
			}
		} else {
			master, err := netlink.LinkByName(high.MasterName)
			if err != nil {
				return errors.WithStack(err)
			}
			if master.Attrs().Index != low.MasterIndex {
				err = netlink.LinkSetMaster(nlLink, master)
				if err != nil {
					return errors.WithStack(err)
				}
			}
		}

		// link mtu
		if high.Mtu != 0 {
			if high.Mtu != uint16(low.MTU) {
				err = netlink.LinkSetMTU(nlLink, int(high.Mtu))
				if err != nil {
					return errors.WithStack(err)
				}
			}
		}

		// link txQueueLen
		if high.TxQueueLen != 0 {
			if high.TxQueueLen != uint16(low.TxQLen) {
				err = netlink.LinkSetTxQLen(nlLink, int(high.TxQueueLen))
				if err != nil {
					return errors.WithStack(err)
				}
			}
		}
	}

	return nil
}
