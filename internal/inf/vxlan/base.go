package vxlan

import (
	"fmt"
	"ifman/internal/inf/common"
	"net"
)

type VxLan struct {
	name        string
	txQueueLen  int
	mtu         int
	masterId    int
	mac         net.HardwareAddr
	vni         int
	src         net.IP
	dst         net.IP
	ttl         int
	tos         int
	learning    bool
	srcPortLow  int
	srcPortHigh int
	port        int
	vtepId      int
}

func (v *VxLan) SetVtep(s string) error {
	id, err := common.PassMaster(s)
	if err != nil {
		return err
	}

	v.vtepId = id
	return nil
}

func (v *VxLan) SetPort(port uint16) {
	v.port = int(port)
}

func (v *VxLan) SetSrcPortHigh(port uint16) {
	v.srcPortHigh = int(port)
}

func (v *VxLan) SetSrcPortLow(port uint16) {
	v.srcPortLow = int(port)
}

func (v *VxLan) SetLearning() {
	v.learning = true
}

func (v *VxLan) SetNoLearning() {
	v.learning = false
}

func (v *VxLan) SetTos(tos uint8) {
	v.tos = int(tos)
}

func (v *VxLan) SetTtl(ttl uint8) {
	v.ttl = int(ttl)
}

func (v *VxLan) SetDst(s string) error {
	ip := net.ParseIP(s)
	if ip == nil {
		return fmt.Errorf("parse ip address failed")
	}

	v.dst = ip
	return nil
}

func (v *VxLan) SetSrc(s string) error {
	ip := net.ParseIP(s)
	if ip == nil {
		return fmt.Errorf("parse ip address failed")
	}

	v.src = ip
	return nil
}

func (v *VxLan) SetVni(vni uint32) error {
	if vni >= 16777216 {
		return fmt.Errorf("invalid vxlan vni range set")
	}

	v.vni = int(vni)
	return nil
}

func (v *VxLan) SetMac(s string) error {
	mac, err := common.PassMac(s)
	if err != nil {
		return err
	}

	v.mac = mac
	return nil
}

func (v *VxLan) SetMaster(s string) error {
	id, err := common.PassMaster(s)
	if err != nil {
		return err
	}

	v.masterId = id
	return nil
}

func (v *VxLan) SetMtu(u uint16) error {
	if err := common.ChkMtu(u); err != nil {
		return err
	}

	v.mtu = int(u)
	return nil
}

func (v *VxLan) SetName(s string) error {
	if err := common.ChkName(s); err != nil {
		return err
	}

	v.name = s
	return nil
}

func (v *VxLan) SetTxQueueLen(u uint16) {
	v.txQueueLen = int(u)
}

func (v *VxLan) check() error {
	if v.name == "" || v.vni == 0 {
		return fmt.Errorf("invalid parameter")
	}

	return nil
}

func GetAttr() *VxLan {
	return &VxLan{
		name:        "",
		txQueueLen:  0,
		mtu:         0,
		masterId:    0,
		mac:         nil,
		vni:         0,
		src:         nil,
		dst:         nil,
		ttl:         0,
		tos:         0,
		learning:    false,
		srcPortLow:  0,
		srcPortHigh: 0,
		port:        0,
	}
}
