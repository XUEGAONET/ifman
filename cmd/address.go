package main

import (
	"golang.org/x/sys/unix"
	"ifman/internal/ip/addr"
)

func afAddress(name, ip string) error {
	a := addr.GetAttr()
	err := a.SetName(name)
	if err != nil {
		return err
	}
	err = a.SetAddr(ip)
	if err != nil {
		return err
	}

	err = addr.New(a)
	if err != nil {
		if err == unix.EEXIST {
			return nil
		} else {
			return err
		}
	}
	return nil
}
