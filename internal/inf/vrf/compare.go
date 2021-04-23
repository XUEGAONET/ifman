package vrf

func Equal(a, b *Vrf) bool {
	if a.mac.String() != b.mac.String() {
		return false
	}

	if a.tableId != b.tableId {
		return false
	}

	return true
}
