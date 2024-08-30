package cmd

import (
	"fmt"
	"github.com/schmiddim/kibana-alert-exporter/helper"
	"runtime"

	"github.com/spf13/cobra"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "version info",
	Long:  `version info`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Code Version:", helper.GitCommit)
		fmt.Println("Go Version:", runtime.Version())
		fmt.Println("GOOS:", runtime.GOOS)
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
