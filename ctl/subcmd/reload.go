package subcmd

import (
	"context"
	"fmt"
	"github.com/XUEGAONET/ifman/common"
	"github.com/XUEGAONET/ifman/proto"
	"github.com/spf13/cobra"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/emptypb"
	"log"
)

var Reload = &cobra.Command{
	Use:     "reload",
	Short:   "Reload the runtime config",
	Long:    "",
	Example: "./ifman-ctl reload",
	Run: func(cmd *cobra.Command, args []string) {
		resp, err := reload()
		if err != nil {
			log.Fatalln(err)
		}

		fmt.Println(resp.Status)
		fmt.Println(resp.Message)
	},
}

func reload() (*proto.ReloadResponse, error) {
	conn, err := grpc.Dial(fmt.Sprintf("127.0.0.1:%d", common.GrpcPort), grpc.WithInsecure())
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	c := proto.NewIfmanClient(conn)

	resp, err := c.ReloadConfig(context.Background(), &emptypb.Empty{})
	if err != nil {
		return nil, err
	}

	return resp, nil
}
