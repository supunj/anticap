// Copyright 2018 Supun Jayathilake(supunj@gmail.com). All rights reserved.

package db

import (
	"context"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/spf13/viper"
	config_util "github.com/supunj/anticap/internal/util/config"
)

// NodesDB - Stores node information
var NodesDB *redis.Client

// ReqResDB - Stores requests and responses
//var ReqResDB redis.Client

// InitDB - Initialize databases
/*func InitDB() error {
	var err error
	NodesDB, err = initNodeDB()
	return err
}*/

// InitNodeDB - Initialise redis
func InitNodeDB() error {
	config_util.Log.Println(fmt.Sprintf("%s%s%s%s", "Redis connection opened -> ", viper.GetString("datastore.redis.host"), ":", viper.GetString("datastore.redis.port")))

	/*NodesDB = redis.NewClient(&redis.Options{
		Network: "",
		Addr:    fmt.Sprintf("%s%s%s", viper.GetString("datastore.redis.host"), ":", viper.GetString("datastore.redis.port")),
		Dialer: func() (ctx, net.Conn, error) {
			return nil, nil
		},
		OnConnect: func(ctx, *redis.Conn) error {
			return nil
		},
		Password:           viper.GetString("datastore.redis.password"),
		DB:                 viper.GetInt("datastore.redis.nodes-db"),
		MaxRetries:         0,
		MinRetryBackoff:    0,
		MaxRetryBackoff:    0,
		DialTimeout:        0,
		ReadTimeout:        0,
		WriteTimeout:       0,
		PoolSize:           0,
		MinIdleConns:       0,
		MaxConnAge:         0,
		PoolTimeout:        0,
		IdleTimeout:        0,
		IdleCheckFrequency: 0,
		TLSConfig:          &tls.Config{},
	})*/

	NodesDB = redis.NewClient(&redis.Options{
		Addr:         fmt.Sprintf("%s%s%s", viper.GetString("datastore.redis.host"), ":", viper.GetString("datastore.redis.port")),
		DB:           viper.GetInt("datastore.redis.nodes-db"),
		DialTimeout:  10 * time.Second,
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 30 * time.Second,
		PoolSize:     10,
		PoolTimeout:  30 * time.Second,
	})

	_, err := NodesDB.Ping(context.Background()).Result()

	return err
}

// CloseDB - Close the databases connection
func CloseDB() error {
	return closeNodeDB()
}

func closeNodeDB() error {
	NodesDB.Close()
	config_util.Log.Println(fmt.Sprintf("%s%s%s%s", "Redis connection closed -> ", viper.GetString("datastore.redis.host"), ":", viper.GetString("datastore.redis.port")))
	return nil
}
