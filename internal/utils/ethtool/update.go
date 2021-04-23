package ethtool

import (
	"context"
	"github.com/sirupsen/logrus"
	"os/exec"
	"time"
)

func Update(name string, o *Offload) {
	setOne(name, "tso", o.TSO)
	setOne(name, "lro", o.LRO)
	setOne(name, "gso", o.GSO)
	setOne(name, "gro", o.GRO)
}

func setOne(name string, typ string, on bool) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	onStr := ""
	if on == true {
		onStr = "on"
	} else {
		onStr = "off"
	}

	c := exec.CommandContext(ctx, "ethtool", "-K", name, typ, onStr)

	output, err := c.CombinedOutput()
	if err != nil {
		logrus.WithField("module", "ethtool").
			Warningf("failed to set %s %s %s: return %s, %v", name, typ, onStr, output, err)
	} else {
		logrus.WithField("module", "ethtool").
			Tracef("set %s %s %s successfully with 0 returned", name, typ, onStr)
	}
}
