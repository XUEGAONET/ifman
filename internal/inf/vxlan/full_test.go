package vxlan

import (
	"fmt"
	"testing"
)

func TestNew(t *testing.T) {
	v := GetAttr()
	err := v.SetName("vxlantest")
	if err != nil {
		panic(err)
	}
	v.SetTxQueueLen(2048)
	err = v.SetMtu(1500)
	if err != nil {
		panic(err)
	}
	err = v.SetMac("94:94:26:a7:c2:11")
	if err != nil {
		panic(err)
	}
	v.SetTtl(16)
	v.SetTos(7)
	v.SetPort(4789)
	v.SetLearning()
	err = v.SetVni(1111)
	if err != nil {
		panic(err)
	}
	err = v.SetDst("1.1.1.1")
	if err != nil {
		panic(err)
	}
	err = v.SetSrc("172.30.12.5")
	if err != nil {
		panic(err)
	}
	v.SetSrcPortLow(4789)
	v.SetSrcPortHigh(10240)

	err = New(v)
	if err != nil {
		panic(err)
	}
}

func TestGet(t *testing.T) {
	v, err := Get("vxlantest")
	if err != nil {
		panic(err)
	}

	fmt.Printf("%#v\n", v)
}
