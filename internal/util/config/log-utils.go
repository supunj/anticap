// Copyright 2018 Supun Jayathilake(supunj@gmail.com). All rights reserved.

package config

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/spf13/viper"
)

// AppLogger - My logger object
type AppLogger struct {
}

var (
	// Log - Logger instance
	Log *log.Logger
)

// InitLog - Create the log file
func InitLog(logpath string) error {
	if viper.GetString("log.outout") == "console" {
		Log = log.New(os.Stdout, "", log.LstdFlags|log.Lshortfile)
		return nil
	}

	fmt.Println("LogFile: " + logpath)

	// Get the default log file path from the config if not passed
	if logpath == "" {
		logpath = viper.GetString("log.default-log-file-path")
	}

	lp := flag.String("logpath", logpath, "Log Path")
	flag.Parse()

	logfile, err := os.Create(*lp)
	if err != nil {
		return err
	}

	Log = log.New(logfile, "", log.LstdFlags|log.Lshortfile)

	return nil
}
