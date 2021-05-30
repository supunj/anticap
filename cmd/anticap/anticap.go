// Copyright 2018 Supun Jayathilake(supunj@gmail.com). All rights reserved.

package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/spf13/viper"
	api_util "github.com/supunj/anticap/internal/api"
	config_util "github.com/supunj/anticap/internal/util/config"
	db_util "github.com/supunj/anticap/internal/util/db"
)

func main() {
	args := os.Args[1:]

	fmt.Println(args)

	var configPath, logPath string = "", ""

	switch argCount := len(args); argCount {
	case 1:
		configPath = args[0]
	case 2:
		configPath = args[0]
		logPath = args[1]
	}

	// Initialize config
	err := config_util.LoadConfig(configPath)

	if err != nil {
		panic(err)
	}

	// Initialize log helper
	err = config_util.InitLog(logPath)
	err = config_util.InitLog(logPath)

	if err != nil {
		panic(err)
	}

	//utils.AppLogger.Info("Log started....")
	config_util.Log.Println("Log started....")

	// Initialize the DB
	err = db_util.InitDB()

	if err != nil {
		panic(err)
	}

	// Bind to a port and pass our router in
	config_util.Log.Println(http.ListenAndServe(fmt.Sprintf("%s%s", ":", viper.GetString("host.port")), api_util.RegisterRoutes()))
	//config_util.Log.Println(http.ListenAndServeTLS(fmt.Sprintf("%s%s", ":", viper.GetString("ssl.port")), fmt.Sprintf("%s%s", ":", viper.GetString("ssl.cert")), fmt.Sprintf("%s%s", ":", viper.GetString("ssl.key")), api_handler.RegisterRoutes()))
}
