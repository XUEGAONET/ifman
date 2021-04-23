package vrf

import (
	"fmt"
	"testing"
)

func TestNew(t *testing.T) {
	v := GetAttr()
	err := v.SetName("test")
	if err != nil {
		panic(err)
	}
	err = v.SetMac("94:94:26:a7:c8:91")
	if err != nil {
		panic(err)
	}
	v.SetTableId(201)

	err = New(v)
	if err != nil {
		panic(err)
	}
}

func TestGet(t *testing.T) {
	v, err := Get("test")
	if err != nil {
		panic(err)
	}

	fmt.Printf("%#v\n", v)
}
