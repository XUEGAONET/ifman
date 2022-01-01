package subcmd

import (
	"fmt"
	"github.com/XUEGAONET/ifman/common"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
	"io/ioutil"
	"log"
	"os"
)

var Test = &cobra.Command{
	Use:     "test [path]",
	Short:   "Test the config file",
	Long:    "",
	Args:    cobra.MinimumNArgs(1),
	Example: "./ifman-ctl test ./config.yaml",
	Run: func(cmd *cobra.Command, args []string) {
		b, err := readFile(args[0])
		if err != nil {
			log.Fatalf("read file failed: %v", err)
		}

		c := common.Config{}
		err = yaml.Unmarshal(b, &c)
		if err != nil {
			log.Fatalf("parse yaml failed: %v", err)
		}

		fmt.Printf("OK\n")
	},
}

func readFile(path string) ([]byte, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	defer f.Close()

	return ioutil.ReadAll(f)
}
