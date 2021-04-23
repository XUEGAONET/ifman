package rp_filter

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"io/ioutil"
	"os"
)

func Update(name string, mode int) error {
	v, err := getValue(name)
	if err != nil {
		return err
	}

	logrus.WithField("module", "rp_filter").
		Tracef("interface %s origin rp_filter value is %d", name, v)

	err = setValue(name, mode)
	return err
}

func getValue(name string) (int, error) {
	f, err := os.Open("/proc/sys/net/ipv4/conf/" + name + "/rp_filter")
	if err != nil {
		return 0, err
	}
	defer f.Close()

	content, err := ioutil.ReadAll(f)
	if err != nil {
		return 0, err
	}

	if len(content) == 0 {
		return 0, fmt.Errorf("linux kernel may have a bug")
	}

	value := content[0]
	switch value {
	case 0x30: // 0
		return 0, nil
	case 0x31: // 1
		return 1, nil
	case 0x32: // 2
		return 2, nil
	default:
		return 0, fmt.Errorf("other rp_filter value meet")
	}
}

func setValue(name string, value int) error {
	f, err := os.OpenFile("/proc/sys/net/ipv4/conf/"+name+"/rp_filter", os.O_RDWR, 0644)
	if err != nil {
		return err
	}
	defer f.Close()

	switch value {
	case OFF: // 0
		_, err = f.WriteString("0")
	case STRICT:
		_, err = f.WriteString("1")
	case LOOSE:
		_, err = f.WriteString("2")
	}

	return err
}
