package main

import (
	"fmt"
	"ifman/internal/utils/ethtool"
	"ifman/internal/utils/rp_filter"
	"strings"
)

func afUtils(c Interface) error {
	offload := ethtool.Offload{
		GRO: c.Gro,
		GSO: c.Gso,
		TSO: c.Tso,
		LRO: c.Lro,
	}
	ethtool.Update(c.Name, &offload)

	mode := 0
	switch strings.ToLower(c.RpFilter) {
	case "off":
		mode = rp_filter.OFF
	case "strict":
		mode = rp_filter.STRICT
	case "loose":
		mode = rp_filter.LOOSE
	default:
		return fmt.Errorf("invalid rp_filter mode")
	}
	return rp_filter.Update(c.Name, mode)
}
