package main

import (
	"golang.org/x/sys/unix"
	"ifman/internal/ip/route"
)

func afRoute(r Route) error {
	rs := route.GetAttr()
	err := rs.SetDst(r.Dst)
	if err != nil {
		return err
	}
	err = rs.SetGw(r.Gw)
	if err != nil {
		return err
	}
	rs.SetTableId(r.TableId)

	err = route.New(rs)
	if err != nil {
		if err == unix.EEXIST {
			return nil
		} else {
			return err
		}
	}

	return nil
}
