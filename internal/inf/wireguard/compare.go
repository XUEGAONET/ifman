package wireguard

func Equal(a, b *WireGuard) bool {
	if a.mtu != b.mtu {
		return false
	}

	if a.txQueueLen != b.txQueueLen {
		return false
	}

	if a.masterId != b.masterId {
		return false
	}

	// endpoint not equal
	if (a.endpoint != nil && b.endpoint == nil) || (a.endpoint == nil && b.endpoint != nil) {
		return false
	}

	// endpoint must equal
	if a.endpoint == nil { // server mode
		if a.listenPort != b.listenPort {
			return false
		}
	} else { // client mode
		if a.endpoint.String() != b.endpoint.String() {
			return false
		}

		if *a.hsInterval != *b.hsInterval {
			return false
		}
	}

	if string(a.key) != string(b.key) {
		return false
	}

	return true
}
