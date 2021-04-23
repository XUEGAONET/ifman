package common

import "github.com/vishvananda/netlink"

func PassMaster(s string) (int, error) {
	link, err := netlink.LinkByName(s)
	if err != nil {
		return 0, err
	}

	return link.Attrs().Index, nil
}
