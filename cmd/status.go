package main

import (
	"fmt"
	"ifman/internal/inf/status"
	"strings"
)

func afStatus(c Interface) error {
	current, err := status.Get(c.Name)
	if err != nil {
		return err
	}

	expectation := 0
	switch strings.ToLower(c.Status) {
	case "up":
		expectation = status.UP
	case "down":
		expectation = status.DOWN
	default:
		return fmt.Errorf("invalid expected status parameter")
	}

	if current == expectation {
		return nil
	}

	switch expectation {
	case status.UP:
		return status.Update(c.Name, status.UP)
	case status.DOWN:
		return status.Update(c.Name, status.DOWN)
	default:
		return fmt.Errorf("invalid expected status parameter")
	}
}
