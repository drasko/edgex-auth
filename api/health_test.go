/**
 * Copyright (c) Mainflux
 *
 * Mainflux server is licensed under an Apache license, version 2.0.
 * All rights not explicitly granted in the Apache license, version 2.0 are reserved.
 * See the included LICENSE file for more details.
 */

package api_test

import (
	"io/ioutil"
	"net/http"
	"testing"
)

func TestHealth(t *testing.T) {
	cases := []struct {
		body string
		code int
	}{
		{`{"running": true}`, 200},
	}

	url := ts.URL + "/health"

	for i, c := range cases {
		res, err := http.Get(url)
		if err != nil {
			t.Errorf("case %d: %s", i+1, err.Error())
		}

		if res.StatusCode != c.code {
			t.Errorf("case %d: expected status %d got %d", i+1, c.code, res.StatusCode)
		}

		body, err := ioutil.ReadAll(res.Body)
		res.Body.Close()
		if err != nil {
			t.Fatalf("case %d: %s", i+1, err.Error())
		}

		if c.body != string(body) {
			t.Errorf("case %d: expected response %s got %s", i+1, c.body, string(body))
		}
	}
}
