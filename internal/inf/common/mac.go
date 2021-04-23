package common

import "net"

func PassMac(s string) (net.HardwareAddr, error) {
	hw, err := net.ParseMAC(s)
	if err != nil {
		return nil, err
	}

	return hw, nil
}
