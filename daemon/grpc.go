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
	"context"
	"fmt"
	"github.com/XUEGAONET/ifman/proto"
	"github.com/pkg/errors"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/emptypb"
	"net"
	"syscall"
)

type grpcServer struct {
	proto.UnimplementedIfmanServer
	server *grpc.Server
}

func NewGrpcServer(port uint16) (*grpcServer, error) {
	s := grpcServer{}

	addr, err := net.ResolveTCPAddr("tcp4", fmt.Sprintf("127.0.0.1:%d", port))
	if err != nil {
		return nil, errors.WithStack(err)
	}

	listener, err := net.ListenTCP("tcp4", addr)
	if err != nil {
		if errors.Is(err, syscall.EADDRINUSE) {
			return nil, errors.WithStack(fmt.Errorf("port exists, do not run the second ifman instance on the same time"))
		} else {
			return nil, errors.WithStack(err)
		}
	}

	server := grpc.NewServer()
	proto.RegisterIfmanServer(server, &s)
	err = server.Serve(listener)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	s.server = server

	return &s, nil
}

func (s *grpcServer) Close() {
	s.server.Stop()
}

func (s *grpcServer) ReloadConfig(_ context.Context, _ *emptypb.Empty) (*proto.ReloadResponse, error) {
	resp := proto.ReloadResponse{
		Status:  "success",
		Message: "",
	}

	err := reloadGlobalConfig()
	if err != nil {
		resp.Status = "fail"
		resp.Message = fmt.Sprintf("%v", err)
	}

	return &resp, nil
}

func (s *grpcServer) Recheck(_ context.Context, _ *emptypb.Empty) (*proto.RecheckResponse, error) {
	resp := proto.RecheckResponse{
		Status:  "success",
		Message: "signal has been sent, please go to log for more details",
	}

	if recheckChan != nil {
		recheckChan <- struct{}{}
	} else {
		resp.Status = "fail"
		resp.Message = "recheckChan has not been init"
	}

	return &resp, nil
}
