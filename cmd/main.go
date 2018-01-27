//
// Copyright (c) 2017
// Mainflux
// Cavium
//
// SPDX-License-Identifier: Apache-2.0
//

package main

import (
	"fmt"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/drasko/edgex-auth/auth"
	"github.com/drasko/edgex-auth/mongo"

	"go.uber.org/zap"
	"gopkg.in/mgo.v2"
)

const (
	authPort               int    = 8089
	authHost               string = "0.0.0.0"
	defMongoURL            string = "0.0.0.0"
	defMongoUsername       string = ""
	defMongoPassword       string = ""
	defMongoDatabase       string = "coredata"
	defMongoPort           int    = 27017
	defMongoConnectTimeout int    = 5000
	defMongoSocketTimeout  int    = 5000
	envMongoURL            string = "AUTH_MONGO_URL"
	envAuthHost          	 string = "AUTH_HOST_URL"
	envAuthPort          	 string = "AUTH_PORT"
)

type config struct {
	AuthPort            int
	AuthHost            string
	MongoURL            string
	MongoUser           string
	MongoPass           string
	MongoDatabase       string
	MongoPort           int
	MongoConnectTimeout int
	MongoSocketTimeout  int
}

func main() {
	logger, _ := zap.NewProduction()
	defer logger.Sync()

	auth.InitLogger(logger)

	cfg := loadConfig(logger)

	ms, err := connectToMongo(cfg)
	if err != nil {
		logger.Error("Failed to connect to Mongo.", zap.Error(err))
		return
	}
	defer ms.Close()

	repo := mongo.NewRepository(ms)
	auth.InitMongoRepository(repo)

	errs := make(chan error, 2)

	auth.StartHTTPServer(cfg.AuthHost, cfg.AuthPort, errs)

	go func() {
		c := make(chan os.Signal)
		signal.Notify(c, syscall.SIGINT)
		errs <- fmt.Errorf("%s", <-c)
	}()

	c := <-errs
	logger.Info("terminated", zap.String("error", c.Error()))
}

func loadConfig(logger *zap.Logger) (*config) {

	p, err := strconv.Atoi(env(envAuthPort, strconv.Itoa(authPort)))
	if err != nil {
		logger.Info("Failed to detect port number. Taking default.", zap.String("error", err.Error()))
		p = authPort
	}

	cfg := config{
		AuthPort: 					 p,
		AuthHost: 					 env(envAuthHost, authHost),
		MongoURL:            env(envMongoURL, defMongoURL),
		MongoUser:           defMongoUsername,
		MongoPass:           defMongoPassword,
		MongoDatabase:       defMongoDatabase,
		MongoPort:           defMongoPort,
		MongoConnectTimeout: defMongoConnectTimeout,
		MongoSocketTimeout:  defMongoSocketTimeout,
	}

	return &cfg
}

func env(key, fallback string) string {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}

	return value
}

func connectToMongo(cfg *config) (*mgo.Session, error) {
	mongoDBDialInfo := &mgo.DialInfo{
		Addrs:    []string{cfg.MongoURL + ":" + strconv.Itoa(cfg.MongoPort)},
		Timeout:  time.Duration(cfg.MongoConnectTimeout) * time.Millisecond,
		Database: cfg.MongoDatabase,
		Username: cfg.MongoUser,
		Password: cfg.MongoPass,
	}

	ms, err := mgo.DialWithInfo(mongoDBDialInfo)
	if err != nil {
		return nil, err
	}

	ms.SetSocketTimeout(time.Duration(cfg.MongoSocketTimeout) * time.Millisecond)
	ms.SetMode(mgo.Monotonic, true)

	return ms, nil
}
