package main

import (
	"github.com/XUEGAONET/ifman/ctl/subcmd"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "ifman",
	Short: "Ifman is a Linux network resource manager",
	Long: `You can manage all the type of interface in Linux with ifman, to build 
your advanced network. 
For more documentation is available at https://blog.xuegaogg.com and https://github.com/XUEGAONET/ifman.`,
	Version: version,
}

func init() {
	rootCmd.AddCommand(subcmd.Reload)
	rootCmd.AddCommand(subcmd.Recheck)
	rootCmd.AddCommand(subcmd.Test)
	rootCmd.AddCommand(subcmd.Key)
}
