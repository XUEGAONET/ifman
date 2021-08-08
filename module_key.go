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

package main

import (
	"encoding/base64"
	"fmt"
	"github.com/pkg/errors"
	"golang.zx2c4.com/wireguard/wgctrl/wgtypes"
	"regexp"
)

var reg *regexp.Regexp = nil

// generateWireGuardKeyPair return public key and private key in string
func generateWireGuardKeyPair() (string, string, error) {
	pri, err := wgtypes.GeneratePrivateKey()
	if err != nil {
		return "", "", errors.Wrap(err, "generate private key failed")
	}
	pub := pri.PublicKey()

	return pub.String(), pri.String(), nil
}

func generateWireGuardKeyChain() {
	pub1, pri1, err := generateWireGuardKeyPair()
	if err != nil {
		panic(err)
	}

	pub2, pri2, err := generateWireGuardKeyPair()
	if err != nil {
		panic(err)
	}

	chain1 := fmt.Sprintf("%s||%s", pri1, pub2)
	chain2 := fmt.Sprintf("%s||%s", pri2, pub1)

	encoded1 := base64.StdEncoding.EncodeToString([]byte(chain1))
	encoded2 := base64.StdEncoding.EncodeToString([]byte(chain2))

	fmt.Printf("WireGuard key chain do not contain [ and ] \n")
	fmt.Printf("* Chain 1: [%s]\n", encoded1)
	fmt.Printf("* Chain 2: [%s]\n", encoded2)
}

// DecodeWireGuardKeyChain will return private key and public key
func DecodeWireGuardKeyChain(kc string) (string, string, error) {
	decodeByte, err := base64.StdEncoding.DecodeString(kc)
	if err != nil {
		return "", "", errors.WithStack(err)
	}
	decode := string(decodeByte)

	if reg == nil {
		reg, err = regexp.Compile("^(?P<Private>.*?)\\|\\|(?P<Public>.*?)$")
		if err != nil {
			return "", "", errors.WithStack(err)
		}
	}

	matches := reg.FindStringSubmatch(decode)
	priIdx := reg.SubexpIndex("Private")
	pubIdx := reg.SubexpIndex("Public")

	return matches[priIdx], matches[pubIdx], nil
}
