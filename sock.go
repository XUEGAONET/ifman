package main

import (
	"fmt"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"net"
	"syscall"
)

func ListenSocket(path string) error {
	c, err := net.Listen("unix", path)
	if err != nil {
		if errors.Is(err, syscall.EADDRINUSE) {
			return errors.WithStack(fmt.Errorf("socket file exist, do not run the second ifman on the same time"))
		} else {
			return errors.WithStack(err)
		}
	}

	go func() {
		for {
			conn, err := c.Accept()
			if err != nil {
				logrus.Warnf("accept conn from unix socket failed: %+v", errors.WithStack(err))
				continue
			}

			go func(conn net.Conn) {
				defer conn.Close()

				buf := make([]byte, 1024)

				n, err := conn.Read(buf)
				if err != nil {
					logrus.Warnf("read from unix conn failed: %+v", errors.WithStack(err))
					return
				}

				switch string(buf[:n]) {
				case "reload":
					err = refreshCoreConfig()
					if err != nil {
						logrus.Errorf("refresh core config failed: %+v", errors.WithStack(err))
					}
				}
			}(conn)
		}
	}()

	return nil
}

func SendReload(path string) error {
	c, err := net.Dial("unix", path)
	if err != nil {
		return errors.WithStack(err)
	}
	defer c.Close()

	_, err = c.Write([]byte("reload"))
	if err != nil {
		return errors.WithStack(err)
	}

	return nil
}
