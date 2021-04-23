package route

import "github.com/vishvananda/netlink"

func New(r *Route) error {
	err := r.check()
	if err != nil {
		return err
	}

	return netlink.RouteAdd(&netlink.Route{
		Dst:   &r.dst,
		Gw:    r.gw,
		Table: r.tableId,
	})
}
