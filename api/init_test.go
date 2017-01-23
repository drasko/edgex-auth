/**
 * Copyright (c) Mainflux
 *
 * Mainflux server is licensed under an Apache license, version 2.0.
 * All rights not explicitly granted in the Apache license, version 2.0 are reserved.
 * See the included LICENSE file for more details.
 */

package api_test

import (
	"fmt"
	"log"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/mainflux/mainflux-core/api"
	mfdb "github.com/mainflux/mainflux-core/db"

	"gopkg.in/mgo.v2"
	"gopkg.in/ory-am/dockertest.v3"
)

var ts *httptest.Server

func TestMain(m *testing.M) {
	var (
		db  *mgo.Session
		err error
	)

	// uses a sensible default on windows (tcp/http) and linux/osx (socket)
	pool, err := dockertest.NewPool("")
	if err != nil {
		log.Fatalf("Could not connect to docker: %s", err)
	}

	// pulls an image, creates a container based on it and runs it
	resource, err := pool.Run("mongo", "3.4", nil)
	if err != nil {
		log.Fatalf("Could not start resource: %s", err)
	}

	// exponential backoff-retry, because the application in the container might not be ready to accept connections yet
	if err := pool.Retry(func() error {
		var err error
		db, err = mgo.Dial(fmt.Sprintf("localhost:%s", resource.GetPort("27017/tcp")))
		if err != nil {
			return err
		}

		mfdb.SetMainSession(db)
		mfdb.SetMainDb("mainflux_test")

		return db.Ping()
	}); err != nil {
		log.Fatalf("Could not connect to docker: %s", err)
	}

	// Start the HTTP server
	ts = httptest.NewServer(api.HTTPServer())
	defer ts.Close()

	code := m.Run()

	// Close database connection.
	db.Close()

	// You can't defer this because os.Exit doesn't care for defer
	if err := pool.Purge(resource); err != nil {
		log.Fatalf("Could not purge resource: %s", err)
	}

	// Exit tests
	os.Exit(code)
}
