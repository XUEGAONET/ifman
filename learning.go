package main

import (
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/vishvananda/netlink"
)

type Learning struct {
	Name       string `yaml:"name"`
	LearningOn bool   `yaml:"learning_on"`
}

func UpdateLearning(learning *Learning) error {
	logrus.Debugf("update learning: %#v", learning)

	link, err := netlink.LinkByName(learning.Name)
	if err != nil {
		return errors.WithStack(err)
	}

	err = netlink.LinkSetLearning(link, learning.LearningOn)
	if err != nil {
		return errors.WithStack(err)
	}

	return nil
}
