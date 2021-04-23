package route

import (
	"golang.org/x/sys/unix"
	"testing"
)

func TestNew(t *testing.T) {
	r := GetAttr()
	err := r.SetDst("8.8.8.8/32")
	if err != nil {
		panic(err)
	}
	err = r.SetGw("127.0.0.1")
	if err != nil {
		panic(err)
	}

	err = New(r)
	if err != nil {
		panic(err)
	}

	if err == unix.EEXIST {
		return
	}
}
