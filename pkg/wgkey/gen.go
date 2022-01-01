package wgkey

import (
	"github.com/pkg/errors"
	"golang.zx2c4.com/wireguard/wgctrl/wgtypes"
)

// GenerateKeyPair return public key and private key in string
func GenerateKeyPair() (string, string, error) {
	pri, err := wgtypes.GeneratePrivateKey()
	if err != nil {
		return "", "", errors.Wrap(err, "generate private key failed")
	}
	pub := pri.PublicKey()

	return pub.String(), pri.String(), nil
}

func GeneratePrivate() (string, error) {
	pri, err := wgtypes.GeneratePrivateKey()
	if err != nil {
		return "", errors.Wrap(err, "generate private key failed")
	}

	return pri.String(), nil
}

func GeneratePublic(private string, number int) ([]string, error) {
	key, err := wgtypes.ParseKey(private)
	if err != nil {
		return nil, errors.Wrap(err, "parse private private key failed")
	}

	var res []string
	for i := 0; i < number; i++ {
		res = append(res, key.PublicKey().String())
	}

	return res, nil
}
