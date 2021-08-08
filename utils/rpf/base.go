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

type RPFType uint8
type IPVersion uint8

const (
	RPF_NONE RPFType = iota
	RPF_STRICT
	RPF_LOOSE
)

// TODO: support ipv6 rp_filter control
const (
	IPV4 IPVersion = iota
	IPV6
)
