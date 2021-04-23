package addr

import "github.com/vishvananda/netlink"

func New(a *Address) error {
	link, err := netlink.LinkByName(a.name)
	if err != nil {
		return err
	}

	addr := netlink.Addr{
		IPNet: &a.addr,
	}
	err = netlink.AddrAdd(link, &addr)
	if err != nil {
		return err
	}

	return nil
}
