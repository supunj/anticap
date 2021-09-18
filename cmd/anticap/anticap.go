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

	_ "github.com/supunj/anticap/api"
)

// @title ANTICAP API
// @version 1.0
// @description This is a sample serice for managing orders
// @termsOfService http://swagger.io/terms/
// @contact.name API Support
// @contact.email soberkoder@swagger.io
// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html
// @host localhost:8080
// @BasePath /
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

	if err != nil {
		panic(err)
	}

	//utils.AppLogger.Info("Log started....")
	config_util.Log.Println("Log started....")

	// Initialize the DB
	err = db_util.InitNodeDB()

	if err != nil {
		panic(err)
	}

	// Bind to a port and pass our router in
	config_util.Log.Println(http.ListenAndServe(fmt.Sprintf("%s%s", ":", viper.GetString("host.port")), api_util.RegisterRoutes()))
	//config_util.Log.Println(http.ListenAndServeTLS(fmt.Sprintf("%s%s", ":", viper.GetString("ssl.port")), fmt.Sprintf("%s%s", ":", viper.GetString("ssl.cert")), fmt.Sprintf("%s%s", ":", viper.GetString("ssl.key")), api_handler.RegisterRoutes()))

	// Close the db connection
	defer db_util.CloseDB()
}
