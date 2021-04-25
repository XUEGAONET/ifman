package wireguard

import "golang.zx2c4.com/wireguard/wgctrl/wgtypes"

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

		if !keyEqual(a.peerPublic, b.peerPublic) {
			return false
		}
	}

	if !keyEqual(a.private, b.private) {
		return false
	}

	return true
}

func keyEqual(a, b *wgtypes.Key) bool {
	if (a != nil && b == nil) || (a == nil || b != nil) {
		return false
	}

	if a != nil { // both are not nil key
		if !compareKey(*a, *b) {
			return false
		}
	}

	return true
}

func compareKey(a, b wgtypes.Key) bool {
	for i := 0; i < 32; i++ {
		if a[i] != b[i] {
			return false
		}
	}

	return true
}
