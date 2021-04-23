package status

import (
	"fmt"
	"github.com/vishvananda/netlink"
)

func Update(name string, status int) error {
	link, err := netlink.LinkByName(name)
	if err != nil {
		return err
	}

	switch status {
	case UP:
		err = netlink.LinkSetUp(link)
	case DOWN:
		err = netlink.LinkSetDown(link)
	default:
		return fmt.Errorf("invalid status parameter")
	}

	return err
}
