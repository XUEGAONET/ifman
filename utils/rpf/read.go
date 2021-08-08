// Copyright 2021 The ifman authors https://github.com/XUEGAONET/ifman
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//    http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

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
