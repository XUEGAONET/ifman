package dummy

func Equal(expect, get *Dummy) bool {
	if expect.mac != nil {
		if expect.mac.String() != get.mac.String() {
			return false
		}
	}

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
