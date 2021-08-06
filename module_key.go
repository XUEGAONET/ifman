package main

import (
	"encoding/base64"
	"fmt"
	"github.com/pkg/errors"
	"golang.zx2c4.com/wireguard/wgctrl/wgtypes"
)

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
