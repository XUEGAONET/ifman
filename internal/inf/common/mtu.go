package common

import "fmt"

func ChkMtu(u uint16) error {
	if u < 64 {
		return fmt.Errorf("invalid mtu range set")
	}

	return nil
}
