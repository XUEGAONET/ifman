package generic

func Equal(a, b *Generic) bool {
	if a.mtu != b.mtu {
		return false
	}

	if a.mac.String() != b.mac.String() {
		return false
	}

	// TODO:
	//   not all the interface can be set
	if a.txQueueLen != b.txQueueLen {
		return false
	}

	if a.masterId != b.masterId {
		return false
	}

	return true
}
