package wireguard

import (
	"fmt"
	"golang.zx2c4.com/wireguard/wgctrl/wgtypes"
	"testing"
)

func TestNew(t *testing.T) {
	wg := GetAttr()
	err := wg.SetName("twg001")
	if err != nil {
		panic(err)
	}
	err = wg.SetMtu(1400)
	if err != nil {
		panic(err)
	}
	wg.SetTxQueueLen(4096)
	err = wg.SetListenPort(9999)
	if err != nil {
		panic(err)
	}

	k, err := wgtypes.GenerateKey()
	if err != nil {
		panic(err)
	}

	err = wg.SetKey(k[:])
	if err != nil {
		panic(err)
	}

	err = New(wg)
	if err != nil {
		panic(err)
	}
}

func TestGet(t *testing.T) {
	attr, err := Get("twg001")
	if err != nil {
		panic(err)
	}

	fmt.Printf("%#v\n", attr)
}
