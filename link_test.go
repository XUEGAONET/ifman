package main

import (
	"fmt"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/vishvananda/netlink"
	"testing"
)

func TestBridgeNew(t *testing.T) {
	logrus.SetLevel(logrus.TraceLevel)

	link := Bridge{
		BaseLink: BaseLink{
			LinkUp:     true,
			Name:       "br-test1",
			TxQueueLen: 512,
			Mtu:        1500,
			MasterName: "",
			Mac:        "94:94:26:a7:c1:11",
		},
		MulticastSnoopingOn: true,
		VlanFilteringOn:     true,
	}
	err := NewLink(&link)
	if err != nil {
		panic(err)
	}

	nl, err := netlink.LinkByName(link.Name)
	if err != nil {
		panic(err)
	}
	defer func() {
		_ = netlink.LinkDel(nl)
	}()
}

func TestBridgeStaticAttrFix(t *testing.T) {
	logrus.SetLevel(logrus.TraceLevel)

	link := Bridge{
		BaseLink: BaseLink{
			LinkUp:     true,
			Name:       "br-test2",
			TxQueueLen: 512,
			Mtu:        1500,
			MasterName: "",
			Mac:        "94:94:26:a7:c1:11",
		},
		MulticastSnoopingOn: true,
		VlanFilteringOn:     true,
	}
	err := NewLink(&link)
	if err != nil {
		panic(err)
	}

	err = compareBridge(&link)
	if err != nil {
		panic(err)
	}

	link.VlanFilteringOn = false
	err = UpdateLink(&link)
	if err != nil {
		panic(err)
	}

	err = compareBridge(&link)
	if err != nil {
		panic(err)
	}

	link.MulticastSnoopingOn = false
	err = UpdateLink(&link)
	if err != nil {
		panic(err)
	}

	err = compareBridge(&link)
	if err != nil {
		panic(err)
	}

	nl, err := netlink.LinkByName(link.Name)
	if err != nil {
		panic(err)
	}
	_ = netlink.LinkDel(nl)
}

func compareBridge(link Link) error {
	nl, err := netlink.LinkByName(link.GetBaseAttrs().Name)
	if err != nil {
		return errors.WithStack(err)
	}

	nlBr := nl.(*netlink.Bridge)

	if *nlBr.MulticastSnooping != link.(*Bridge).MulticastSnoopingOn {
		return fmt.Errorf("multicast snooping")
	}

	if *nlBr.VlanFiltering != link.(*Bridge).VlanFilteringOn {
		return fmt.Errorf("vlan filter")
	}

	return nil
}

func TestIPTunNew(t *testing.T) {
	link := IPTun{
		BaseLink: BaseLink{
			LinkUp:     true,
			Name:       "ipip-test",
			TxQueueLen: 1500,
			Mtu:        1500,
			MasterName: "",
			Mac:        "",
		},
		Ttl:      16,
		Tos:      7,
		LocalIP:  "127.0.0.1",
		RemoteIP: "127.0.0.1",
	}

	err := NewLink(&link)
	if err != nil {
		panic(err)
	}

	nl, err := netlink.LinkByName(link.Name)
	if err != nil {
		panic(err)
	}

	nlipip := nl.(*netlink.Iptun)
	if nlipip.Ttl != link.Ttl {
		panic("ttl")
	}

	if nlipip.Tos != link.Tos {
		panic("tos")
	}

	if nlipip.Remote.String() != link.RemoteIP {
		panic("remote")
	}

	if nlipip.Local.String() != link.LocalIP {
		panic("local")
	}

	defer func() {
		_ = netlink.LinkDel(nl)
	}()
}

// TODO: add tun test case

func TestVlanNew(t *testing.T) {
	logrus.SetLevel(logrus.TraceLevel)

	link := Vlan{
		BaseLink: BaseLink{
			LinkUp:     true,
			Name:       "vlan-test",
			TxQueueLen: 512,
			Mtu:        1500,
			MasterName: "",
			Mac:        "",
		},
		VlanId:     10,
		StackingOn: true,
		BindLink:   "eth0",
	}

	err := NewLink(&link)
	if err != nil {
		panic(err)
	}

	nl, err := netlink.LinkByName(link.Name)
	if err != nil {
		panic(err)
	}

	nlVl := nl.(*netlink.Vlan)
	if nlVl.VlanProtocol != netlink.VLAN_PROTOCOL_8021AD {
		panic("stacking")
	}

	if nlVl.VlanId != 10 {
		panic("vlan id")
	}

	defer func() {
		_ = netlink.LinkDel(nl)
	}()
}

// TODO: add static attrs support
//func TestVlanUpdate(t *testing.T) {
//	logrus.SetLevel(logrus.TraceLevel)
//
//	link := Vlan{
//		BaseLink: BaseLink{
//			LinkUp:     true,
//			Name:       "vlan-test2",
//			TxQueueLen: 1000,
//			Mtu:        1400,
//			MasterName: "",
//			Mac:        "94:94:26:a7:cc:cc",
//		},
//		BindLink:   "eth0",
//		VlanId:     15,
//		StackingOn: false,
//	}
//
//
//}

func TestVrfNew(t *testing.T) {
	logrus.SetLevel(logrus.TraceLevel)

	link := Vrf{
		BaseLink: BaseLink{
			LinkUp:     true,
			Name:       "vrf-test",
			TxQueueLen: 1500,
			Mtu:        1400,
			MasterName: "",
			Mac:        "94:94:26:aa:aa:aa",
		},
		TableId: 199,
	}

	err := NewLink(&link)
	if err != nil {
		panic(err)
	}

	nl, err := netlink.LinkByName(link.Name)
	if err != nil {
		panic(err)
	}

	nlVl := nl.(*netlink.Vrf)
	if nlVl.Table != link.TableId {
		panic("table id")
	}

	defer func() {
		_ = netlink.LinkDel(nl)
	}()
}

func TestVxlanNewRemoteIP(t *testing.T) {
	logrus.SetLevel(logrus.TraceLevel)

	link := VxLAN{
		BaseLink: BaseLink{
			LinkUp:     true,
			Name:       "vxlan-test",
			TxQueueLen: 128,
			Mtu:        1400,
			MasterName: "",
			Mac:        "94:94:26:a7:22:22",
		},
		Vni:         1000,
		SrcIp:       "127.0.0.1",
		DstIP:       "127.0.0.1",
		Ttl:         16,
		Tos:         4,
		LearningOn:  false,
		SrcPortLow:  4789,
		SrcPortHigh: 4789,
		Port:        4789,
		VtepName:    "",
	}

	err := NewLink(&link)
	if err != nil {
		panic(err)
	}

	nl, err := netlink.LinkByName(link.Name)
	if err != nil {
		panic(err)
	}

	logrus.Debugf("%#v", nl)

	defer func() {
		_ = netlink.LinkDel(nl)
	}()
}

func TestVxlanNewVtep(t *testing.T) {
	logrus.SetLevel(logrus.TraceLevel)

	link := VxLAN{
		BaseLink: BaseLink{
			LinkUp:     true,
			Name:       "vxlan-test2",
			TxQueueLen: 128,
			Mtu:        1400,
			MasterName: "",
			Mac:        "94:94:26:a7:22:21",
		},
		Vni:         1010,
		SrcIp:       "",
		DstIP:       "0.0.0.0",
		Ttl:         16,
		Tos:         4,
		LearningOn:  false,
		SrcPortLow:  4789,
		SrcPortHigh: 4789,
		Port:        4789,
		VtepName:    "eth0",
	}

	err := NewLink(&link)
	if err != nil {
		panic(err)
	}

	nl, err := netlink.LinkByName(link.Name)
	if err != nil {
		panic(err)
	}

	logrus.Debugf("%#v", nl)

	defer func() {
		_ = netlink.LinkDel(nl)
	}()
}
