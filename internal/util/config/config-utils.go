// Copyright 2018 Supun Jayathilake(supunj@gmail.com). All rights reserved.

package config

import (
	"fmt"
	"os"

	"github.com/spf13/viper"
)

// LoadConfig loads the configuration
func LoadConfig(configFile string) error {
	viper.SetConfigType("json")
	if configFile != "" {
		viper.SetConfigFile(configFile)
	} else {
		currentfolderPath, _ := os.Getwd()
		viper.SetConfigFile(currentfolderPath + "/configs/config.dev.json")
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
