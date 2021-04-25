package tun

func Equal(expect, get *Tun) bool {
	if expect.mtu != 0 {
		if expect.mtu != get.mtu {
			return false
		}
	}

	if expect.txQueueLen != 0 {
		if expect.txQueueLen != get.txQueueLen {
			return false
		}
	}

	if expect.masterId != 0 {
		if expect.masterId != get.masterId {
			return false
		}
	}

	return true
}
