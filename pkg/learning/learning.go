package learning

import (
	"github.com/pkg/errors"
	"github.com/vishvananda/netlink"
)

func SetLearningOn(name string) error {
	link, err := netlink.LinkByName(name)
	if err != nil {
		return errors.WithStack(err)
	}

	err = netlink.LinkSetLearning(link, true)
	if err != nil {
		return errors.WithStack(err)
	}

	return nil
}

func SetLearningOff(name string) error {
	link, err := netlink.LinkByName(name)
	if err != nil {
		return errors.WithStack(err)
	}

	err = netlink.LinkSetLearning(link, false)
	if err != nil {
		return errors.WithStack(err)
	}

	return nil
}
