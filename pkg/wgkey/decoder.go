package wgkey

import (
	"encoding/base64"
	"github.com/pkg/errors"
	"regexp"
)

var reg *regexp.Regexp = nil

// DecodeKeyChain will return private key and public key
func DecodeKeyChain(kc string) (string, string, error) {
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
