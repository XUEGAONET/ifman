package main

import (
	"fmt"
	"golang.zx2c4.com/wireguard/wgctrl"
	"testing"
)

func TestWgGet(t *testing.T) {
	cli, err := wgctrl.New()
	if err != nil {
		panic(err)
	}
	defer cli.Close()

	devs, err := cli.Devices()
	if err != nil {
		panic(err)
	}

	for i, _ := range devs {
		dev := devs[i]

		k := dev.PrivateKey.String()
		fmt.Printf("dump struct: %+v\n", dev)
		fmt.Printf(" * private key: %s\n", k)
	}
}
