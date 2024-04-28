package cmd

import (
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// startCmd represents the start command
var startCmd = &cobra.Command{
	Use:   "start",
	Short: "Start service",
	Long:  "Start " + viper.GetString("setting.service_name") + " service",
	Run: func(cmd *cobra.Command, args []string) {
		server := NewServer()
		server.Start()
	},
}
