package tun

func Equal(a, b *Tun) bool {
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
