// Copyright 2018 Supun Jayathilake(supunj@gmail.com). All rights reserved.

package db

import (
	"fmt"

	"github.com/go-redis/redis"
	"github.com/spf13/viper"
	config_util "github.com/supunj/anticap/internal/util/config"
)

// NodesDB - Stores node information
var NodesDB redis.Client

// ReqResDB - Stores requests and responses
//var ReqResDB redis.Client

// InitDB - Initialize databases
func InitDB() error {
	var err error
	NodesDB, err = initNodeDB()
	return err
}

func initNodeDB() (redis.Client, error) {

	config_util.Log.Println(fmt.Sprintf("%s%s%s", viper.GetString("datastore.redis.host"), ":", viper.GetString("datastore.redis.port")))

	ndb := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s%s%s", viper.GetString("datastore.redis.host"), ":", viper.GetString("datastore.redis.port")),
		Password: viper.GetString("datastore.redis.password"), // no password set
		DB:       viper.GetInt("datastore.redis.nodes-db"),    // use default DB
	})

	/*
		ReqResDB := redis.NewClient(&redis.Options{
			Addr:     fmt.Sprintf("%s%s%s", viper.GetString("datastore.redis.host"), ":", viper.GetString("datastore.redis.port")),
			Password: viper.GetString("datastore.redis.password"), // no password set
			DB:       viper.GetInt("datastore.redis.req-res-db"),  // use default DB
		}) */

	//pong1, err := ndb.Ping().Result()
	//pong2, err := ReqResDB.Ping().Result()

	//config_util.Log.Println(pong1, err)

	return *ndb, nil
}
