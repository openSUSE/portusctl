// Copyright (C) 2017 Miquel Sabaté Solà <msabate@suse.com>
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License
// along with this program.  If not, see <http://www.gnu.org/licenses/>.

package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"time"
)

const (
	apiPrefix = "/api"
	v1Prefix  = apiPrefix + "/v1"
)

var requestTimeout = 15 * time.Second

func generateBody(resource *Resource, args []string) ([]byte, error) {
	values, err := extractArguments(resource, args, false)
	if err != nil {
		return nil, err
	}

	var b []byte
	var e error

	if resource.bundler == "" {
		b, e = json.Marshal(values)
	} else {
		b, e = json.Marshal(map[string]map[string]string{resource.bundler: values})
	}
	if e != nil {
		return nil, fmt.Errorf("Internal error: %v", e)
	}
	return b, nil
}

func buildRequest(method, path, query string, body []byte) (*http.Request, error) {
	u, _ := url.Parse(globalConfig.server)
	u.Path = path
	u.RawQuery = query

	reader := bytes.NewBuffer(body)
	req, err := http.NewRequest(method, u.String(), reader)
	if err != nil {
		return nil, fmt.Errorf("Could not build request: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("PORTUS-AUTH", fmt.Sprintf("%s:%s", globalConfig.user, globalConfig.token))
	return req, nil
}

func request(method, path, query string, body []byte) (*http.Response, error) {
	req, err := buildRequest(method, path, query, body)
	if err != nil {
		return nil, err
	}

	client := http.Client{Timeout: requestTimeout}
	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	if err := handleCode(res); err != nil {
		return nil, err
	}
	return res, nil
}

func handleCode(response *http.Response) error {
	code := response.StatusCode

	if code == http.StatusBadRequest {
		target := make(map[string]map[string][]string)
		defer response.Body.Close()

		json.NewDecoder(response.Body).Decode(&target)
		return messagesError(target["errors"])
	} else if code == http.StatusUnauthorized || code == http.StatusForbidden {
		key := "errors"
		if code == http.StatusUnauthorized {
			key = "error"
		}
		target := make(map[string]string)
		defer response.Body.Close()
		json.NewDecoder(response.Body).Decode(&target)
		return errors.New(target[key])
	} else if code == http.StatusNotFound {
		return errors.New("Resource not found")
	} else if code == http.StatusMethodNotAllowed {
		return errors.New("Action not allowed on given resource")
	} else if code >= http.StatusInternalServerError && code <= http.StatusNetworkAuthenticationRequired {
		fmt.Printf("Portus faced an error:\n\n")
		target := make(map[string]string)
		defer response.Body.Close()
		json.NewDecoder(response.Body).Decode(&target)
		return errors.New(target["error"])
	}
	return nil
}
