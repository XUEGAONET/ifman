package subcmd

import (
	"encoding/base64"
	"fmt"
	"github.com/XUEGAONET/ifman/pkg/wgkey"
	"github.com/spf13/cobra"
	"log"
)

var Key = &cobra.Command{
	Use:     "key",
	Short:   "Generate wireguard key pair and chain",
	Example: "./ifman-ctl key",
	Run: func(cmd *cobra.Command, args []string) {
		pub1, pri1, err := wgkey.GenerateKeyPair()
		if err != nil {
			log.Fatalln(err)
		}

		pub2, pri2, err := wgkey.GenerateKeyPair()
		if err != nil {
			log.Fatalln(err)
		}

		chain1 := fmt.Sprintf("%s||%s", pri1, pub2)
		chain2 := fmt.Sprintf("%s||%s", pri2, pub1)

		encoded1 := base64.StdEncoding.EncodeToString([]byte(chain1))
		encoded2 := base64.StdEncoding.EncodeToString([]byte(chain2))

		fmt.Printf("WireGuard key chain do not contain '[' and ']' \n")
		fmt.Printf("* 1 Private: \t[%s]\n", pri1)
		fmt.Printf("* 1 Public: \t[%s]\n", pub1)
		fmt.Printf("* 2 Private: \t[%s]\n", pri2)
		fmt.Printf("* 2 Public: \t[%s]\n", pub2)
		fmt.Printf("* 1 Chain: \t[%s]\n", encoded1)
		fmt.Printf("* 2 Chain: \t[%s]\n", encoded2)
	},
}
