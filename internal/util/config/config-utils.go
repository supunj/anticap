// Copyright 2018 Supun Jayathilake(supunj@gmail.com). All rights reserved.

package config

import (
	"fmt"

	"github.com/spf13/viper"
)

// LoadConfig loads the configuration
func LoadConfig(configFile string) error {
	viper.SetConfigType("json")
	if configFile != "" {
		viper.SetConfigFile(configFile)
	} else {
		viper.SetConfigFile("/mnt/hdd1/gitlab/anticap/configs/config.dev.json")
	}
	// Find and read the config file
	if err := viper.ReadInConfig(); err != nil {
		fmt.Printf("Error reading config file, %s", err)
		return err
	}
	// Confirm which config file is used
	fmt.Printf("Using config: %s\n", viper.ConfigFileUsed())

	return nil
}
