package dummy

func Equal(a, b *Dummy) bool {
	if a.mac.String() != b.mac.String() {
		return false
	}

	if a.mtu != b.mtu {
		return false
	}

	if a.txQueueLen != b.txQueueLen {
		return false
	}

	if a.masterId != b.masterId {
		return false
	}

	return true
}
