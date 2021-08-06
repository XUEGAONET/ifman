package main

import (
	"fmt"
	rpfp "github.com/XUEGAONET/ifman/utils/rpf"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

type RpFilter struct {
	Name string `yaml:"name"`
	Mode string `yaml:"mode"`
}

func UpdateRpFilter(rpf *RpFilter) error {
	logrus.Debugf("update rp_filter: %#v", rpf)

	var want rpfp.RPFType

	switch rpf.Mode {
	case "off":
		want = rpfp.RPF_NONE
	case "strict":
		want = rpfp.RPF_STRICT
	case "loose":
		want = rpfp.RPF_LOOSE
	default:
		return fmt.Errorf("unsupported rp_filter mode")
	}

	err := rpfp.CheckAndFix(rpf.Name, want)
	if err != nil {
		return errors.WithStack(err)
	}

	return nil
}
