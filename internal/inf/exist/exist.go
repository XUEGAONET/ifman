package exist

import (
	"github.com/sirupsen/logrus"
	"github.com/vishvananda/netlink"
)

func IsExisted(name string) bool {
	_, err := netlink.LinkByName(name)
	if err == nil {
		return true
	}

	_, ok := err.(netlink.LinkNotFoundError)
	if ok {
		return false
	}
	_, ok = err.(*netlink.LinkNotFoundError)
	if ok {
		return false
	}

	logrus.WithField("module", "exist").
		Errorf("interface existed but occur another err: %v", err)
	return false
}
