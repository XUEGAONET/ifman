package status

import (
	"github.com/vishvananda/netlink"
	"net"
)

func Get(name string) (int, error) {
	link, err := netlink.LinkByName(name)
	if err != nil {
		return 0, err
	}

	if link.Attrs().Flags&net.FlagUp == net.FlagUp {
		return UP, nil
	} else {
		return DOWN, nil
	}
}
