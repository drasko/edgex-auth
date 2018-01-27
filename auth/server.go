//
// Copyright (c) 2017 Mainflux
//
// SPDX-License-Identifier: Apache-2.0
//

package auth

import (
	"fmt"
	"net/http"

	"github.com/go-zoo/bone"
	"go.uber.org/zap"
)

// HTTPServer function
func httpServer() http.Handler {
	mux := bone.New()

	mux.Get("/status", http.HandlerFunc(getStatus))
	mux.Get("/users", http.HandlerFunc(getAllUsers))
	mux.Get("/users/:id", http.HandlerFunc(getUserByID))
	mux.Post("/users", http.HandlerFunc(createUser))
	mux.Delete("/users/:id", http.HandlerFunc(deleteUser))
	mux.Post("/login", http.HandlerFunc(login))
	mux.Post("/auth", http.HandlerFunc(authorize))
	return mux
}

func StartHTTPServer(host string, port int, errChan chan error) {
	go func() {
		p := fmt.Sprintf("%s:%d", host, port)
		logger.Info("Starting EdgeX Auth", zap.String("url", p))
		errChan <- http.ListenAndServe(p, httpServer())
	}()
}
