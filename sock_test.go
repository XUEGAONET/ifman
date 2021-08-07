package main

import (
	"github.com/sirupsen/logrus"
	"strings"
	"testing"
	"time"
)

func TestListenCtl(t *testing.T) {
	logrus.SetLevel(logrus.TraceLevel)
	var port uint16 = 11111

	go func() {
		err := ListenCtl(port)
		if err != nil {
			panic(err)
		}
	}()

	time.Sleep(time.Second)

	err := SendReload(port)
	if err != nil {
		panic(err)
	}

	err = ListenCtl(port)
	if err == nil {
		panic("nil pointer! bug!")
	}
	if !strings.Contains(err.Error(), "second") {
		panic("bug")
	}

	time.Sleep(time.Second)
}
