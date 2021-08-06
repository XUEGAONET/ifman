package pid

import (
	"fmt"
	"os"
	"testing"
)

func TestNew(t *testing.T) {
	p := New("/var/run/ifman.pid")
	err := p.Init()
	if err != nil {
		panic(err)
	}
	defer p.Remove()

	pid, err := p.Get()
	if err != nil {
		panic(err)
	}

	fmt.Println("Pid: ", pid)

	if os.Getpid() != pid {
		panic("pid not equal")
	}
}

func TestFile_Get(t *testing.T) {
	p := New("/var/run/ifman.pid")
	err := p.Init()
	if err != nil {
		panic(err)
	}
	defer p.Remove()

	pid, err := p.Get()
	if err != nil {
		panic(err)
	}

	p1 := New("/var/run/ifman.pid")
	pid1, err := p1.Get()
	if err != nil {
		panic(err)
	}

	if pid1 != pid {
		panic("bug!")
	}
}
