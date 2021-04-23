package rp_filter

import (
	"github.com/sirupsen/logrus"
	"os"
	"testing"
)

func TestUpdate(t *testing.T) {
	logrus.SetLevel(logrus.TraceLevel)
	logrus.SetOutput(os.Stdout)

	err := Update("tuntest", OFF)
	if err != nil {
		panic(err)
	}
}
