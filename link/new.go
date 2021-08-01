/**
 * Copyright (c) 2021 Harris <ic0xgkk@gmail.com>
 *
 * Permission is hereby granted, free of charge, to any person obtaining a copy of
 * this software and associated documentation files (the "Software"), to deal in
 * the Software without restriction, including without limitation the rights to
 * use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of
 * the Software, and to permit persons to whom the Software is furnished to do so,
 * subject to the following conditions:
 *
 * The above copyright notice and this permission notice shall be included in all
 * copies or substantial portions of the Software.
 *
 * THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
 * IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS
 * FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR
 * COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER
 * IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN
 * CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.
 */

package link

import (
	"errors"
	"net"
	"strings"
	"time"

	"github.com/vishvananda/netlink"
	"golang.org/x/sys/unix"
	"golang.zx2c4.com/wireguard/wgctrl"
	"golang.zx2c4.com/wireguard/wgctrl/wgtypes"
)

func New(l interface{}) error {
	switch nl := l.(type) {
	case *VxLAN:
		instance := netlink.Vxlan{
			LinkAttrs: netlink.LinkAttrs{},
		}

		if nl.Name == "" {
			return errors.New("blank interface name")
		} else {
			instance.Name = nl.Name
		}

		instance.TxQLen = int(nl.TxQueueLen)
		instance.MTU = int(nl.Mtu)

		if nl.MasterName != "" {
			masterIndex, err := GetLinkIndexByName(nl.MasterName)
			if err != nil {
				return err
			}
			instance.MasterIndex = masterIndex
		}

		if nl.Mac != "" {
			mac, err := MacString2Struct(nl.Mac)
			if err != nil {
				return err
			}
			instance.HardwareAddr = mac
		}

		if nl.Vni >= 2^24 {
			return errors.New("wrong vni range")
		}

		if nl.SrcIp != "" {
			ip, err := IpString2Struct(nl.SrcIp)
			if err != nil {
				return err
			}
			instance.SrcAddr = ip
		}

		if nl.DstIP != "" {
			ip, err := IpString2Struct(nl.DstIP)
			if err != nil {
				return err
			}
			instance.Group = ip
		}

		instance.TTL = int(nl.Ttl)
		instance.TOS = int(nl.Tos)
		instance.Learning = nl.Learning
		instance.PortLow = int(nl.SrcPortLow)
		instance.PortHigh = int(nl.SrcPortHigh)
		instance.Port = int(nl.Port)

		if nl.VtepName != "" {
			index, err := GetLinkIndexByName(nl.VtepName)
			if err != nil {
				return err
			}
			instance.VtepDevIndex = index
		}

		return netlink.LinkAdd(&instance)
	case *Bridge:
		instance := netlink.Bridge{
			LinkAttrs: netlink.LinkAttrs{},
		}

		if nl.Name == "" {
			return errors.New("blank interface name")
		} else {
			instance.Name = nl.Name
		}

		instance.TxQLen = int(nl.TxQueueLen)
		instance.MTU = int(nl.Mtu)

		if nl.MasterName != "" {
			masterIndex, err := GetLinkIndexByName(nl.MasterName)
			if err != nil {
				return err
			}
			instance.MasterIndex = masterIndex
		}

		if nl.Mac != "" {
			mac, err := MacString2Struct(nl.Mac)
			if err != nil {
				return err
			}
			instance.HardwareAddr = mac
		}

		return netlink.LinkAdd(&instance)
	case *Vrf:
		instance := netlink.Vrf{
			LinkAttrs: netlink.LinkAttrs{},
		}

		if nl.Name == "" {
			return errors.New("blank interface name")
		} else {
			instance.Name = nl.Name
		}

		instance.TxQLen = int(nl.TxQueueLen)
		instance.MTU = int(nl.Mtu)

		if nl.Mac != "" {
			mac, err := MacString2Struct(nl.Mac)
			if err != nil {
				return err
			}
			instance.HardwareAddr = mac
		}

		return netlink.LinkAdd(&instance)
	case *Dummy:
		instance := netlink.Dummy{
			LinkAttrs: netlink.LinkAttrs{},
		}

		if nl.Name == "" {
			return errors.New("blank interface name")
		} else {
			instance.Name = nl.Name
		}

		instance.TxQLen = int(nl.TxQueueLen)
		instance.MTU = int(nl.Mtu)

		if nl.MasterName != "" {
			masterIndex, err := GetLinkIndexByName(nl.MasterName)
			if err != nil {
				return err
			}
			instance.MasterIndex = masterIndex
		}

		if nl.Mac != "" {
			mac, err := MacString2Struct(nl.Mac)
			if err != nil {
				return err
			}
			instance.HardwareAddr = mac
		}

		return netlink.LinkAdd(&instance)
	case *WireGaurdPtP:
		instance := WireGuardLink{
			LinkAttrs: netlink.LinkAttrs{},
		}

		if nl.Name == "" {
			return errors.New("blank interface name")
		} else {
			instance.Name = nl.Name
		}

		instance.TxQLen = int(nl.TxQueueLen)
		instance.MTU = int(nl.Mtu)

		if nl.MasterName != "" {
			masterIndex, err := GetLinkIndexByName(nl.MasterName)
			if err != nil {
				return err
			}
			instance.MasterIndex = masterIndex
		}

		err := netlink.LinkAdd(&instance)
		if err != nil {
			return err
		}

		// begin to config wireguard

		wgCtlCli, err := wgctrl.New()
		if err != nil {
			return err
		}
		defer wgCtlCli.Close()

		_, cidr, err := net.ParseCIDR("0.0.0.0/0")
		if err != nil {
			return err
		}

		localKey, err := wgtypes.ParseKey(nl.Private)
		if err != nil {
			return err
		}

		peerKey, err := wgtypes.ParseKey(nl.PeerPublic)
		if err != nil {
			return err
		}

		var (
			confPeer  wgtypes.PeerConfig
			confLocal wgtypes.Config
		)

		// server mode
		if nl.Endpoint == "" && nl.ListenPort != 0 {
			confPeer = wgtypes.PeerConfig{
				PublicKey:                   peerKey,
				Remove:                      false,
				UpdateOnly:                  false,
				PresharedKey:                nil,
				Endpoint:                    nil,
				PersistentKeepaliveInterval: nil,
				ReplaceAllowedIPs:           true,
				AllowedIPs:                  []net.IPNet{*cidr},
			}

			listen := int(nl.ListenPort)
			confLocal = wgtypes.Config{
				PrivateKey:   &localKey,
				ListenPort:   &listen,
				FirewallMark: nil,
				ReplacePeers: true,
				Peers:        []wgtypes.PeerConfig{confPeer},
			}
		} else if nl.Endpoint != "" && nl.ListenPort == 0 { // client mode
			hs := time.Duration(nl.HeartbeatInterval) * time.Second

			endpoint, err := net.ResolveUDPAddr("udp", nl.Endpoint)
			if err != nil {
				return err
			}

			confPeer = wgtypes.PeerConfig{
				PublicKey:                   peerKey,
				Remove:                      false,
				UpdateOnly:                  false,
				PresharedKey:                nil,
				Endpoint:                    endpoint,
				PersistentKeepaliveInterval: &hs,
				ReplaceAllowedIPs:           true,
				AllowedIPs:                  []net.IPNet{*cidr},
			}

			confLocal = wgtypes.Config{
				PrivateKey:   &localKey,
				ListenPort:   nil,
				FirewallMark: nil,
				ReplacePeers: true,
				Peers:        []wgtypes.PeerConfig{confPeer},
			}
		}

		err = wgCtlCli.ConfigureDevice(nl.Name, confLocal)
		if err != nil {
			return err
		}

		return nil
	case *Tun:
		instance := netlink.Tuntap{
			LinkAttrs:  netlink.LinkAttrs{},
			Mode:       unix.IFF_TUN,
			Flags:      netlink.TUNTAP_DEFAULTS,
			NonPersist: false, // always persist
		}

		if nl.Name == "" {
			return errors.New("blank interface name")
		} else {
			instance.Name = nl.Name
		}

		instance.TxQLen = int(nl.TxQueueLen)
		instance.MTU = int(nl.Mtu)

		if nl.MasterName != "" {
			masterIndex, err := GetLinkIndexByName(nl.MasterName)
			if err != nil {
				return err
			}
			instance.MasterIndex = masterIndex
		}

		if nl.MultiQueue {
			instance.Flags = netlink.TUNTAP_MULTI_QUEUE_DEFAULTS
		}

		return netlink.LinkAdd(&instance)
	default:
		return errors.New("wrong interface type")
	}
}

func MacString2Struct(mac string) (net.HardwareAddr, error) {
	return net.ParseMAC(mac)
}

func MacStruct2String(mac net.HardwareAddr) string {
	return strings.ToLower(mac.String())
}

func GetLinkIndexByName(name string) (int, error) {
	link, err := netlink.LinkByName(name)
	if err != nil {
		return 0, err
	}

	return link.Attrs().Index, nil
}

func IpString2Struct(ip string) (net.IP, error) {
	res := net.ParseIP(ip)
	if res == nil {
		return nil, errors.New("parse ip failed")
	}
	return res, nil
}
