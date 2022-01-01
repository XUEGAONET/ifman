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

var Recheck = &cobra.Command{
	Use:     "recheck",
	Short:   "Recheck the runtime config and fix",
	Long:    "",
	Example: "./ifman-ctl recheck",
	Run: func(cmd *cobra.Command, args []string) {
		resp, err := recheck()
		if err != nil {
			log.Fatalln(err)
		}

		fmt.Println(resp.Status)
		fmt.Println(resp.Message)
	},
}

func recheck() (*proto.RecheckResponse, error) {
	conn, err := grpc.Dial(fmt.Sprintf("127.0.0.1:%d", common.GrpcPort), grpc.WithInsecure())
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	c := proto.NewIfmanClient(conn)

	resp, err := c.Recheck(context.Background(), &emptypb.Empty{})
	if err != nil {
		return nil, err
	}

	return resp, nil
}
