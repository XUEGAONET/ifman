// Copyright 2021 The ifman authors https://github.com/XUEGAONET/ifman
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//    http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"fmt"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"net"
	"syscall"
)

func ListenCtl(port uint16) error {
	addr, err := net.ResolveTCPAddr("tcp4", fmt.Sprintf("127.0.0.1:%d", port))
	if err != nil {
		return errors.WithStack(err)
	}

	c, err := net.ListenTCP("tcp4", addr)
	if err != nil {
		if errors.Is(err, syscall.EADDRINUSE) {
			return errors.WithStack(fmt.Errorf("port exist, do not run the second ifman on the same time"))
		} else {
			return errors.WithStack(err)
		}
	}

	go func() {
		for {
			conn, err := c.Accept()
			if err != nil {
				logrus.Warnf("accept conn from socket failed: %+v", errors.WithStack(err))
				continue
			}

			go func(conn net.Conn) {
				defer conn.Close()

				buf := make([]byte, 1024)

				n, err := conn.Read(buf)
				if err != nil {
					logrus.Warnf("read from conn failed: %+v", errors.WithStack(err))
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

func SendReload(port uint16) error {
	c, err := net.Dial("tcp4", fmt.Sprintf("127.0.0.1:%d", port))
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
