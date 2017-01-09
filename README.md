# Mainflux Auth Service

[![License](https://img.shields.io/badge/license-Apache%20v2.0-blue.svg)](LICENSE)
[![Build Status](https://travis-ci.org/mainflux/mainflux-auth.svg?branch=master)](https://travis-ci.org/mainflux/mainflux-auth)
[![Go Report Card](https://goreportcard.com/badge/github.com/Mainflux/mainflux-auth)](https://goreportcard.com/report/github.com/Mainflux/mainflux-auth)

Authentication server in Go.

### Install
Go Auth Server uses [Redis](https://redis.io/), so insure that it is installed on your system.

Installing Go Auth Server is trivial [`go get`](https://golang.org/cmd/go/):
```bash
go get github.com/drasko/go-auth
$GOBIN/mainflux-auth
```

### Features
- JWT
- User account creation
- Redis persistance
- Login and logout
- Token revocation

### License
[Apache License, version 2.0](LICENSE)
