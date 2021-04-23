package common

import (
	"fmt"
	"strings"
)

func ChkName(s string) error {
	if strings.Contains(s, "/") || strings.Contains(s, " ") {
		return fmt.Errorf("invalid character in interface name")
	}

	if len(s) >= 16 {
		return fmt.Errorf("interface name out of length")
	}

	return nil
}
