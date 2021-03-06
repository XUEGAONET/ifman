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

package common

import (
	"github.com/vishvananda/netlink"
)

const (
	LinkTypeLayer2 = iota
	LinkTypeLayer3
)

type Link interface {
	GetBaseAttrs() *BaseLink
}

type BaseLink struct {
	LinkUp     bool   `yaml:"link_up"`
	Name       string `yaml:"name"`
	TxQueueLen uint16 `yaml:"tx_queue_len"`
	Mtu        uint16 `yaml:"mtu"`
	MasterName string `yaml:"master_name"`
	Mac        string `yaml:"mac"`
	Comment    string `yaml:"comment"`
}

func (receiver *BaseLink) GetBaseAttrs() *BaseLink {
	return receiver
}

type Bridge struct {
	BaseLink `yaml:",inline"`

	MulticastSnoopingOn bool `yaml:"multicast_snooping_on"`
	VlanFilteringOn     bool `yaml:"vlan_filtering_on"`
}

type Dummy struct {
	BaseLink `yaml:",inline"`
}

type IPTun struct {
	BaseLink `yaml:",inline"`

	Ttl      uint8  `yaml:"ttl"`
	Tos      uint8  `yaml:"tos"`
	LocalIP  string `yaml:"local_ip"`
	RemoteIP string `yaml:"remote_ip"`
}

// Generic is a common link type
type Generic struct {
	BaseLink `yaml:",inline"`
}

type Tun struct {
	BaseLink `yaml:",inline"`

	MultiQueueOn bool  `yaml:"multi_queue_on"`
	PersistOn    bool  `yaml:"persist_on"`
	Queues       uint8 `yaml:"queues"`
}

type Vlan struct {
	BaseLink `yaml:",inline"`

	BindLink   string `yaml:"bind_link"`
	VlanId     uint16 `yaml:"vlan_id"`
	StackingOn bool   `yaml:"stacking_on"`
}

type Vrf struct {
	BaseLink `yaml:",inline"`

	TableId uint32 `yaml:"table_id"`
}

type VxLAN struct {
	BaseLink `yaml:",inline"`

	Vni         uint32 `yaml:"vni"`
	SrcIp       string `yaml:"src_ip"`
	DstIP       string `yaml:"dst_ip"`
	Ttl         uint8  `yaml:"ttl"`
	Tos         uint8  `yaml:"tos"`
	Checksum    bool   `yaml:"checksum"`
	LearningOn  bool   `yaml:"learning_on"`
	SrcPortLow  uint16 `yaml:"src_port_low"`
	SrcPortHigh uint16 `yaml:"src_port_high"`
	Port        uint16 `yaml:"port"`
	VtepName    string `yaml:"vtep_name"`
}

type WireGuardPtPServer struct {
	BaseLink `yaml:",inline"`

	ListenPort uint16 `yaml:"listen_port"`
	Private    string `yaml:"private"`
	PeerPublic string `yaml:"peer_public"`

	KeyChain string `yaml:"key_chain"`
}

type WireGuardPtPClient struct {
	BaseLink `yaml:",inline"`

	Endpoint          string `yaml:"endpoint"`
	HeartbeatInterval uint32 `yaml:"heartbeat_interval"`
	Private           string `yaml:"private"`
	PeerPublic        string `yaml:"peer_public"`

	KeyChain string `yaml:"key_chain"`
}

type WireGuardOrigin struct {
	BaseLink `yaml:",inline"`

	ListenPort uint16                `yaml:"listen_port"`
	Private    string                `yaml:"private"`
	Peers      []WireGuardOriginPeer `yaml:"peers"`
}

type WireGuardOriginPeer struct {
	PeerPublic        string   `yaml:"peer_public"`
	AllowedCIDR       []string `yaml:"allowed_cidr"`
	Endpoint          string   `yaml:"endpoint"`
	HeartbeatInterval uint32   `yaml:"heartbeat_interval"`
}

type WireGuardLink struct {
	netlink.LinkAttrs
}

func (w *WireGuardLink) Attrs() *netlink.LinkAttrs {
	return &w.LinkAttrs
}

func (w *WireGuardLink) Type() string {
	return "wireguard"
}

func GetLinkType(link Link) int {
	switch link.(type) {
	case *Bridge:
		return LinkTypeLayer2
	case *Dummy:
		return LinkTypeLayer2
	case *IPTun:
		return LinkTypeLayer3
	case *Tun:
		return LinkTypeLayer3
	case *Vlan:
		return LinkTypeLayer2
	case *Vrf:
		return LinkTypeLayer2
	case *VxLAN:
		return LinkTypeLayer2
	case *WireGuardPtPServer:
		return LinkTypeLayer3
	case *WireGuardPtPClient:
		return LinkTypeLayer3
	case *WireGuardOrigin:
		return LinkTypeLayer3
	default:
		return LinkTypeLayer3
	}
}
