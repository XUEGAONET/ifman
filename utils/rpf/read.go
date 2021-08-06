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

package rpf

import (
	"errors"
	"io/ioutil"
	"os"

	"github.com/vishvananda/netlink"
)

func Read(name string) (RPFType, error) {
	link, err := netlink.LinkByName(name)
	if err != nil {
		return 0, err
	}

	f, err := os.Open("/proc/sys/net/ipv4/conf/" + link.Attrs().Name + "/rp_filter")
	if err != nil {
		return 0, err
	}
	defer f.Close()

	content, err := ioutil.ReadAll(f)
	if err != nil {
		return 0, err
	}

	if len(content) == 0 {
		return 0, errors.New("no content in rp_filter proc file")
	}

	value := content[0]
	switch value {
	case 0x30: // 0
		return RPF_NONE, nil
	case 0x31: // 1
		return RPF_STRICT, nil
	case 0x32: // 2
		return RPF_LOOSE, nil
	default:
		return 0, errors.New("invalid rp_filter mode")
	}
}