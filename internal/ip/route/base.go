package route

import (
	"fmt"
	"net"
)

type Route struct {
	dst     net.IPNet
	gw      net.IP
	tableId int
}

func (r *Route) SetDst(s string) error {
	_, dst, err := net.ParseCIDR(s)
	if err != nil {
		return err
	}

	r.dst = *dst
	return nil
}

func (r *Route) SetGw(s string) error {
	ip := net.ParseIP(s)
	if ip == nil {
		return fmt.Errorf("parse gateway address failed")
	}

	r.gw = ip
	return nil
}

func (r *Route) SetTableId(u uint8) {
	r.tableId = int(u)
}

func (r *Route) check() error {
	if r.dst.IP == nil || r.gw == nil {
		return fmt.Errorf("invalid parameter")
	}

	return nil
}

func GetAttr() *Route {
	return &Route{
		dst:     net.IPNet{},
		gw:      nil,
		tableId: 0,
	}
}
