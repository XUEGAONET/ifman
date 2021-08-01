/**
 * Copyright (c) 2021 Harris <ic0xgkk@gmail.com>
 *
 * Permission is hereby granted, free of charge, to any person obtaining a copy of
 * this software and associated documentation files (the "Software"), to deal in
 * the Software without restriction, including without limitation the rights to
 * use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of
 * the Software, and to permit persons to whom the Software is furnished to do so,
 * subject to the following conditions:
 *
 * The above copyright notice and this permission notice shall be included in all
 * copies or substantial portions of the Software.
 *
 * THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
 * IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS
 * FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR
 * COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER
 * IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN
 * CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.
 */

package link

import (
	"fmt"
	"testing"

	"github.com/vishvananda/netlink"
)

func TestNewDummy(t *testing.T) {
	i := Dummy{
		Name:       "dummytest",
		TxQueueLen: 8,
		Mtu:        1480,
		MasterName: "",
		Mac:        "",
	}

	err := New(&i)
	if err != nil {
		panic(err)
	}

	link, err := netlink.LinkByName(i.Name)
	if err != nil {
		panic(err)
	}

	l := link.(*netlink.Dummy)
	fmt.Printf("%+v\n", l)

	err = netlink.LinkDel(link)
	if err != nil {
		panic(err)
	}
}

func TestNewTun(t *testing.T) {
	i := Tun{
		Name:       "tuntest",
		TxQueueLen: 8,
		Mtu:        1480,
		MasterName: "",
		MultiQueue: true,
	}

	err := New(&i)
	if err != nil {
		panic(err)
	}

	link, err := netlink.LinkByName(i.Name)
	if err != nil {
		panic(err)
	}

	l := link.(*netlink.Tuntap)
	fmt.Printf("%+v\n", l)

	err = netlink.LinkDel(link)
	if err != nil {
		panic(err)
	}
}
