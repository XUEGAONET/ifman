package main

import (
	"fmt"
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

const (
	LinkTypeLayer2 = iota
	LinkTypeLayer3
)

func NewLink(link Link) error {
	logrus.Debugf("new link: %#v", link)

	var err error = nil
	var nLink netlink.Link = nil

	switch l := link.(type) {
	case *Dummy:
		nLink = &netlink.Dummy{
			LinkAttrs: netlink.LinkAttrs{},
		}
	case *Bridge:
		nLink = &netlink.Bridge{
			LinkAttrs:         netlink.LinkAttrs{},
			MulticastSnooping: &l.MulticastSnoopingOn,
			HelloTime:         nil,
			VlanFiltering:     &l.VlanFilteringOn,
		}
	case *IPTun:
		nLink = &netlink.Iptun{
			LinkAttrs:  netlink.LinkAttrs{},
			Ttl:        l.Ttl,
			Tos:        l.Tos,
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

		if l.LocalIP != "" {
			ip := net.ParseIP(l.LocalIP)
			if ip == nil {
				return errors.WithStack(fmt.Errorf("parse ip address failed"))
			}
			nLink.(*netlink.Iptun).Local = ip
		}

		if l.RemoteIP != "" {
			ip := net.ParseIP(l.RemoteIP)
			if ip == nil {
				return errors.WithStack(fmt.Errorf("parse ip address failed"))
			}
			nLink.(*netlink.Iptun).Remote = ip
		}
	case *Tun:
		nLink = &netlink.Tuntap{
			LinkAttrs:  netlink.LinkAttrs{},
			Mode:       unix.IFF_TUN,
			Flags:      netlink.TUNTAP_DEFAULTS,
			NonPersist: !l.PersistOn,
			Queues:     int(l.Queues),
			Fds:        nil,
			Owner:      0,
			Group:      0,
		}

		if l.MultiQueueOn {
			nLink.(*netlink.Tuntap).Flags = netlink.TUNTAP_MULTI_QUEUE_DEFAULTS
		}
	case *Vlan:
		nLink = &netlink.Vlan{
			LinkAttrs:    netlink.LinkAttrs{},
			VlanId:       int(l.VlanId),
			VlanProtocol: netlink.VLAN_PROTOCOL_8021Q,
		}

		if l.StackingOn {
			nLink.(*netlink.Vlan).VlanProtocol = netlink.VLAN_PROTOCOL_8021AD
		}
	case *Vrf:
		nLink = &netlink.Vrf{
			LinkAttrs: netlink.LinkAttrs{},
			Table:     l.TableId,
		}
	case *VxLAN:
		nLink = &netlink.Vxlan{
			LinkAttrs:      netlink.LinkAttrs{},
			VxlanId:        int(l.Vni),
			VtepDevIndex:   0,
			SrcAddr:        nil,
			Group:          nil,
			TTL:            int(l.Ttl),
			TOS:            int(l.Tos),
			Learning:       l.LearningOn,
			Proxy:          false,
			RSC:            false,
			L2miss:         false,
			L3miss:         false,
			UDPCSum:        false,
			UDP6ZeroCSumTx: false,
			UDP6ZeroCSumRx: false,
			NoAge:          false,
			GBP:            false,
			FlowBased:      false,
			Age:            0,
			Limit:          0,
			Port:           int(l.Port),
			PortLow:        int(l.SrcPortLow),
			PortHigh:       int(l.SrcPortHigh),
		}

		if l.SrcIp != "" {
			ip := net.ParseIP(l.SrcIp)
			if ip == nil {
				return errors.WithStack(errors.New("parse ip address failed"))
			}
			nLink.(*netlink.Vxlan).SrcAddr = ip
		}

		if l.DstIP != "" {
			ip := net.ParseIP(l.DstIP)
			if ip == nil {
				return errors.WithStack(errors.New("parse ip address failed"))
			}
			nLink.(*netlink.Vxlan).Group = ip
		}

		if l.VtepName != "" {
			id, err := netlink.LinkByName(l.VtepName)
			if err != nil {
				return errors.WithStack(err)
			}
			nLink.(*netlink.Vxlan).VtepDevIndex = id.Attrs().Index
		}
	case *WireGuardPtPServer:
		nLink = &WireGuardLink{
			LinkAttrs: netlink.LinkAttrs{},
		}
	case *WireGuardPtPClient:
		nLink = &WireGuardLink{
			LinkAttrs: netlink.LinkAttrs{},
		}
	default:
		return errors.WithStack(fmt.Errorf("unsupport high level link type"))
	}

	err = newBase(link, nLink)
	if err != nil {
		return errors.WithStack(err)
	}

	err = netlink.LinkAdd(nLink)
	if err != nil {
		return errors.WithStack(err)
	}

	return nil
}

func newBase(hLink Link, lLink netlink.Link) error {
	// high level attrs
	high := hLink.GetBaseAttrs()
	// low level attrs
	low := lLink.Attrs()

	switch getLinkType(hLink) {
	case LinkTypeLayer2:
		if high.Mac != "" {
			mac, err := net.ParseMAC(strings.ToLower(high.Mac))
			if err != nil {
				return errors.WithStack(err)
			}
			low.HardwareAddr = mac
		}

		fallthrough
	case LinkTypeLayer3:

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

func UpdateLink(link Link) error {
	logrus.Debugf("update link: %#v", link)

	var err error = nil
	var nLink netlink.Link = nil

	nLink, err = netlink.LinkByName(link.GetBaseAttrs().Name)
	if err != nil {
		return errors.WithStack(err)
	}

	switch l := link.(type) {
	case *Dummy:
	case *Bridge:
	case *IPTun:
	case *Unmanaged:
	case *Tun:
	case *Vlan:
	case *Vrf:
	case *VxLAN:
	case *WireGuardPtPServer:
		pri, pub, err := DecodeWireGuardKeyChain(l.KeyChain)
		if err != nil {
			return errors.WithStack(err)
		}
		l.Private = pri
		l.PeerPublic = pub

		err = checkAndRebuildWireGuard(l)
		if err != nil {
			return errors.WithStack(err)
		}
	case *WireGuardPtPClient:
		pri, pub, err := DecodeWireGuardKeyChain(l.KeyChain)
		if err != nil {
			return errors.WithStack(err)
		}
		l.Private = pri
		l.PeerPublic = pub

		err = checkAndRebuildWireGuard(l)
		if err != nil {
			return errors.WithStack(err)
		}
	default:
		return errors.WithStack(fmt.Errorf("unsupport high level link type"))
	}

	err = updateBase(link, nLink)
	if err != nil {
		return errors.WithStack(err)
	}

	return nil
}

func checkAndRebuildWireGuard(link Link) error {
	cli, err := wgctrl.New()
	if err != nil {
		return errors.WithStack(err)
	}
	defer cli.Close()

	wgCurrent, err := cli.Device(link.GetBaseAttrs().Name)
	if err != nil {
		return errors.WithStack(err)
	}

	var conf *wgtypes.Config
	var needFix = false

	switch l := link.(type) {
	case *WireGuardPtPServer:
		if wgCurrent.PrivateKey.String() != l.Private {
			conf, err = wgServerConf(l.PeerPublic, l.Private, int(l.ListenPort))
			if err != nil {
				return errors.WithStack(err)
			}

			needFix = true
		}
	case *WireGuardPtPClient:
		if wgCurrent.PrivateKey.String() != l.Private {
			conf, err = wgClientConf(l.PeerPublic, l.Private, l.Endpoint, int(l.HeartbeatInterval))
			if err != nil {
				return errors.WithStack(err)
			}

			needFix = true
		}
	default:
		return errors.WithStack(fmt.Errorf("unsupport check and rebuild wireguard type"))
	}

	if needFix {
		err = cli.ConfigureDevice(link.GetBaseAttrs().Name, *conf)
		if err != nil {
			return errors.WithStack(err)
		}
	}

	return nil
}

func wgServerConf(peerPubKey, privateKey string, listenPort int) (*wgtypes.Config, error) {
	_, allow, err := net.ParseCIDR("0.0.0.0/0")
	if err != nil {
		return nil, errors.WithStack(err)
	}

	pub, err := wgtypes.ParseKey(peerPubKey)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	pri, err := wgtypes.ParseKey(privateKey)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	wgPeer := wgtypes.PeerConfig{
		PublicKey:                   pub,
		Remove:                      false,
		UpdateOnly:                  false,
		PresharedKey:                nil,
		Endpoint:                    nil,
		PersistentKeepaliveInterval: nil,
		ReplaceAllowedIPs:           true,
		AllowedIPs:                  []net.IPNet{*allow},
	}

	wgConf := wgtypes.Config{
		PrivateKey:   &pri,
		ListenPort:   &listenPort,
		FirewallMark: nil,
		ReplacePeers: true,
		Peers:        []wgtypes.PeerConfig{wgPeer},
	}

	return &wgConf, nil
}

func wgClientConf(peerPubKey, privateKey string, endpoint string, heartbeatIntervalSec int) (*wgtypes.Config, error) {
	_, allow, err := net.ParseCIDR("0.0.0.0/0")
	if err != nil {
		return nil, errors.WithStack(err)
	}

	pub, err := wgtypes.ParseKey(peerPubKey)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	pri, err := wgtypes.ParseKey(privateKey)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	ep, err := net.ResolveUDPAddr("udp", endpoint)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	heart := time.Duration(heartbeatIntervalSec) * time.Second

	wgPeer := wgtypes.PeerConfig{
		PublicKey:                   pub,
		Remove:                      false,
		UpdateOnly:                  false,
		PresharedKey:                nil,
		Endpoint:                    ep,
		PersistentKeepaliveInterval: &heart,
		ReplaceAllowedIPs:           true,
		AllowedIPs:                  []net.IPNet{*allow},
	}

	wgConf := wgtypes.Config{
		PrivateKey:   &pri,
		ListenPort:   nil,
		FirewallMark: nil,
		ReplacePeers: true,
		Peers:        []wgtypes.PeerConfig{wgPeer},
	}

	return &wgConf, nil
}

func updateBase(hLink Link, lLink netlink.Link) error {
	// high level attrs
	high := hLink.GetBaseAttrs()
	// low level attrs
	low := lLink.Attrs()

	var err error = nil

	switch getLinkType(hLink) {
	case LinkTypeLayer2:
		if high.Mac != "" {
			if strings.ToLower(low.HardwareAddr.String()) != strings.ToLower(high.Mac) {
				newMac, err := net.ParseMAC(strings.ToLower(high.Mac))
				if err != nil {
					return errors.WithStack(err)
				}
				err = netlink.LinkSetHardwareAddr(lLink, newMac)
				if err != nil {
					return errors.WithStack(err)
				}
			}
		}

		fallthrough
	case LinkTypeLayer3:

		fallthrough
	default:
		// link status
		linkUp := low.Flags&net.FlagUp == net.FlagUp
		if linkUp != high.LinkUp {
			if high.LinkUp {
				err = netlink.LinkSetUp(lLink)
			} else {
				err = netlink.LinkSetDown(lLink)
			}

			if err != nil {
				return errors.WithStack(err)
			}
		}

		// link master
		if high.MasterName == "" {
			if low.MasterIndex != 0 {
				err = netlink.LinkSetNoMaster(lLink)
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
				err = netlink.LinkSetMaster(lLink, master)
				if err != nil {
					return errors.WithStack(err)
				}
			}
		}

		// link mtu
		if high.Mtu != 0 {
			if high.Mtu != uint16(low.MTU) {
				err = netlink.LinkSetMTU(lLink, int(high.Mtu))
				if err != nil {
					return errors.WithStack(err)
				}
			}
		}

		// link txQueueLen
		if high.TxQueueLen != 0 {
			if high.TxQueueLen != uint16(low.TxQLen) {
				err = netlink.LinkSetTxQLen(lLink, int(high.TxQueueLen))
				if err != nil {
					return errors.WithStack(err)
				}
			}
		}
	}

	return nil
}
