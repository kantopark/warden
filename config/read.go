package config

import (
	"fmt"
	"log"

	"github.com/spf13/viper"
)

// Reads in the configuration settings. By default, the configuration is set with the
// config.yaml. This sets up the basic configurations. After which, another configuration
// file for secrets "config-secret" is read. This will override the basic configuration
// file.
//
// Configurations are searched for in their unix/Windows specific folders first before the
// local directory (/etc/kantopark/warden or C:\\kantopark\\warden first). Configuration
// files must be of the following format: json|yaml|toml.
func ReadConfig() {
	viper.SetConfigName("config")
	viper.AddConfigPath("/etc/kantopark/warden")           // unix
	viper.AddConfigPath("C:\\kantopark\\warden")           // windows
	viper.AddConfigPath("C:\\projects\\kantopark\\warden") // windows
	viper.AddConfigPath(".")

	if err := viper.ReadInConfig(); err != nil {
		log.Fatalln(fmt.Errorf("Error reading in config: %s\n", err))
	}

	viper.SetConfigName("config-secret")
	if err := viper.MergeInConfig(); err != nil {
		log.Fatalln(fmt.Errorf("Error merging in config: %s\n", err))
	}
}
