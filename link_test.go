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
