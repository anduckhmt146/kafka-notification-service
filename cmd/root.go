package cmd

import (
	"expvar"
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cfgFile string

var rootCmd = &cobra.Command{
	Use:     viper.GetString("setting.service_name"),
	Short:   viper.GetString("setting.service_name") + " start",
	Long:    viper.GetString("setting.service_name") + " start",
	Version: viper.GetString("setting.version"),
}

func SetVersion(r string) {
	if len(r) > 0 {
		rootCmd.Version = r
	}
	expvar.NewString("service_version").Set(rootCmd.Version)
}

func GetVersion() string {
	return rootCmd.Version
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		panic(err)
	}
}

func initConfig() {
	if cfgFile == "" {
		cfgFile = os.Getenv("CONFIG_PATH")
	}
	fmt.Printf("CONFIG_PATH: %s", cfgFile)
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Use default config file
		viper.SetConfigFile("/config/local.yaml")
	}
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_", "-", "_"))
	viper.SetConfigType("yaml")
	viper.AutomaticEnv()

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Printf("Using config file: %s", viper.ConfigFileUsed())
	} else {
		fmt.Printf("failed to read config file. %+v", err)
	}
	service_name := viper.GetString("setting.service_name")
	fmt.Printf("Service name: %s", service_name)
}

func init() {
	os.Setenv("TZ", "Asia/Ho_Chi_Minh")
	cobra.OnInitialize(initConfig)
	rootCmd.AddCommand(startCmd)
	rootCmd.PersistentFlags().StringVar(&cfgFile, "CONFIG_PATH", "", "config file (default is $PWD/config/local.yaml)")
}
