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
	"os"

	"github.com/vishvananda/netlink"
)

func Write(name string, mode RPFType) error {
	link, err := netlink.LinkByName(name)
	if err != nil {
		return err
	}

	f, err := os.OpenFile("/proc/sys/net/ipv4/conf/"+link.Attrs().Name+"/rp_filter", os.O_RDWR, 0644)
	if err != nil {
		return err
	}
	defer f.Close()

	switch mode {
	case RPF_NONE: // 0
		_, err = f.WriteString("0")
	case RPF_STRICT:
		_, err = f.WriteString("1")
	case RPF_LOOSE:
		_, err = f.WriteString("2")
	}

	return err
}
