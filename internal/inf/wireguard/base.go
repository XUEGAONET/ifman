package wireguard

import (
	"fmt"
	"github.com/vishvananda/netlink"
	"ifman/internal/inf/common"
	"net"
	"time"
)

type typWg struct {
	netlink.LinkAttrs
}

func (w *typWg) Attrs() *netlink.LinkAttrs {
	return &w.LinkAttrs
}

func (w *typWg) Type() string {
	return "wireguard"
}

type WireGuard struct {
	name       string
	mtu        int
	txQueueLen int
	masterId   int
	listenPort int
	key        []byte
	endpoint   *net.UDPAddr
	hsInterval *time.Duration
}

func (w *WireGuard) SetName(s string) error {
	if err := common.ChkName(s); err != nil {
		return err
	}

	w.name = s
	return nil
}

func (w *WireGuard) SetMtu(u uint16) error {
	if err := common.ChkMtu(u); err != nil {
		return err
	}

	w.mtu = int(u)
	return nil
}

func (w *WireGuard) SetTxQueueLen(u uint16) {
	w.txQueueLen = int(u)
}

func (w *WireGuard) SetMaster(s string) error {
	id, err := common.PassMaster(s)
	if err != nil {
		return err
	}

	w.masterId = id
	return nil
}

func (w *WireGuard) SetListenPort(u uint16) error {
	if u == 0 {
		return fmt.Errorf("invalid port range")
	}

	w.listenPort = int(u)
	return nil
}

func (w *WireGuard) SetKey(k []byte) error {
	if len(k) != 32 {
		return fmt.Errorf("invalid key")
	}

	t := make([]byte, 32)
	copy(t, k)
	w.key = t
	return nil
}

func (w *WireGuard) SetEndpoint(s string) error {
	addr, err := net.ResolveUDPAddr("udp", s)
	if err != nil {
		return err
	}

	w.endpoint = addr
	return nil
}

func (w *WireGuard) SetHandshakeIntervalSec(u uint16) {
	t := time.Duration(u) * time.Second
	w.hsInterval = &t
}

func (w *WireGuard) check() error {
	if w.name == "" || w.key == nil ||
		(w.hsInterval != nil && w.listenPort != 0 && w.endpoint != nil) {
		return fmt.Errorf("invalid parameter")
	}
	return nil
}

func GetAttr() *WireGuard {
	return &WireGuard{
		name:       "",
		mtu:        0,
		txQueueLen: 0,
		masterId:   0,
		listenPort: 0,
		key:        nil,
		endpoint:   nil,
		hsInterval: nil,
	}
}
