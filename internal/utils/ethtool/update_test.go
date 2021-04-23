package ethtool

import (
	"github.com/sirupsen/logrus"
	"os"
	"testing"
)

func TestUpdate(t *testing.T) {
	logrus.SetLevel(logrus.TraceLevel)
	logrus.SetOutput(os.Stdout)

	o := GetAttr()
	o.GRO = false
	o.GRO = false
	o.LRO = false
	o.TSO = false

	Update("tuntest", o)
}
