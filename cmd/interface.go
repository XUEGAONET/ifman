package main

import "fmt"

func afInterface(c Interface) error {
	var err error

	// create or update interface
	switch c.Type {
	case "generic":
		err = afGeneric(c)
	case "dummy":
		err = afDummy(c)
	case "tun":
		err = afTun(c)
	case "vrf":
		err = afVrf(c)
	case "vxlan":
		err = afVxLan(c)
	case "wireguard":
		err = afWireGuard(c)
	default:
		return fmt.Errorf("invalid afInterface type")
	}
	if err != nil {
		return err
	}

	// update interface status
	err = afStatus(c)
	if err != nil {
		return err
	}

	// set utils
	err = afUtils(c)
	if err != nil {
		return err
	}

	// add or update interface address
	for i, _ := range c.Address {
		err = afAddress(c.Name, c.Address[i])
		if err != nil {
			return err
		}
	}

	return nil
}
