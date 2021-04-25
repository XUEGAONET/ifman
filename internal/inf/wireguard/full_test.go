package wireguard

import (
	"fmt"
	"testing"
)

func TestNew(t *testing.T) {
	wg := GetAttr()
	err := wg.SetName("twg011")
	if err != nil {
		panic(err)
	}
	err = wg.SetMtu(1400)
	if err != nil {
		panic(err)
	}
	wg.SetTxQueueLen(4096)
	err = wg.SetListenPort(12001)
	if err != nil {
		panic(err)
	}

	//k, err := wgtypes.GeneratePrivateKey()
	//if err != nil {
	//	panic(err)
	//}
	//
	//kp := k.PublicKey()
	//
	//fmt.Printf("private: %s\npublic: %s\n", k.String(), kp.String())

	err = wg.SetPrivate("AAEVzA2C6Fvem5DNnge1T4BAhSgCoO9tJEo8G7ZpsEU=")
	if err != nil {
		panic(err)
	}

	err = wg.SetPeerPublic("DHh7F/FVyOJ6GXkS/pvCnrlc6QKRcyeevm+j4KKC9VE=")
	if err != nil {
		panic(err)
	}

	err = New(wg)
	if err != nil {
		panic(err)
	}
}

func TestGet(t *testing.T) {
	attr, err := Get("twg020")
	if err != nil {
		panic(err)
	}

	fmt.Printf("%#v\n", attr)
}

func TestNew2(t *testing.T) {
	wg := GetAttr()
	err := wg.SetName("twg021")
	if err != nil {
		panic(err)
	}
	err = wg.SetMtu(1400)
	if err != nil {
		panic(err)
	}
	wg.SetTxQueueLen(4096)

	err = wg.SetEndpoint("100.100.100.100:53")
	if err != nil {
		panic(err)
	}

	wg.SetHandshakeIntervalSec(1)

	err = wg.SetPrivate("yOkvmylKT4LN2o601MiVntwFMWB5bSD9GOCotoyo81I=")
	if err != nil {
		panic(err)
	}

	err = wg.SetPeerPublic("V628GeW8NyDsqPS1EonOxeLDhPzRikJ4wbcvRhXErho=")
	if err != nil {
		panic(err)
	}

	err = New(wg)
	if err != nil {
		panic(err)
	}
}
