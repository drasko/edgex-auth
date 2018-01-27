//
// Copyright (c) 2017
// Cavium
//
// SPDX-License-Identifier: Apache-2.0
//

package auth

type Config struct {
	Port int
	Host string
}

var AuthCfg Config

func InitConfig(host string, port int) Config {
	return Config{
		Host: host,
		Port: port,
	}
}
