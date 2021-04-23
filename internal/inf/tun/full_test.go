package tun

import (
	"fmt"
	"testing"
)

func TestNew(t *testing.T) {
	tun := GetAttr()
	err := tun.SetName("tuntest")
	if err != nil {
		panic(err)
	}
	err = tun.SetMtu(1500)
	if err != nil {
		panic(err)
	}
	tun.SetTxQueueLen(2048)
	tun.SetMultiQueue()
	tun.SetPersist()

	err = New(tun)
	if err != nil {
		panic(err)
	}
}

func TestGet(t *testing.T) {
	tun, err := Get("tuntest")
	if err != nil {
		panic(err)
	}
	fmt.Printf("%#v\n", tun)
}

func TestUpdate(t *testing.T) {
	tun := GetAttr()
	err := tun.SetName("tuntest")
	if err != nil {
		panic(err)
	}
	err = tun.SetMtu(1400)
	if err != nil {
		panic(err)
	}
	tun.SetTxQueueLen(2048)
	tun.SetMultiQueue()
	tun.SetPersist()

	err = Update(tun)
	if err != nil {
		panic(err)
	}
}
