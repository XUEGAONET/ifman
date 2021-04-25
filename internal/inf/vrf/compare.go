package vrf

func Equal(expect, get *Vrf) bool {
	if expect.mac != nil {
		if expect.mac.String() != get.mac.String() {
			return false
		}
	}

	if expect.tableId != 0 {
		if expect.tableId != get.tableId {
			return false
		}
	}

	return true
}
