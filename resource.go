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
	"io/ioutil"
	"net/http"
	"net/url"
	"path/filepath"
	"strings"
	"time"
	"unicode"
	"unicode/utf8"
)

const (
	pathPrefix = "/api/v1"
)

var requestTimeout = 15 * time.Second

var subresources = map[string][]string{
	"namespaces":   []string{"repositories"},
	"repositories": []string{"tags"},
	"teams":        []string{"namespaces", "members"},
	"users":        []string{"application_tokens"},
}

type argumentRequirement struct {
	required []string
	optional []string
	prefix   string
	resource string
	name     string
	returned string
}

// TODO: more OO when refactoring
// TODO: restrict the methods that can be used (e.g. a team cannot be destroyed
// as of now)
var bodyArguments = map[string]argumentRequirement{
	"users": argumentRequirement{
		required: []string{"username", "email", "password"},
		optional: []string{"display_name"},
		prefix:   "user",
		name:     "user",
	},
	"application_tokens": argumentRequirement{
		required: []string{"id", "application"},
		resource: "users",
		name:     "application token",
		returned: "plain_tokens",
	},
	"teams": argumentRequirement{
		required: []string{"name"},
		optional: []string{"description"},
		name:     "team",
		returned: "teams",
	},
	"namespaces": argumentRequirement{
		required: []string{"name", "team"},
		optional: []string{"description"},
		name:     "namespace",
		returned: "namespaces",
	},
}

func buildRequest(method, prefix, resource string, args []string, body []byte) (*http.Request, error) {
	u, _ := url.Parse(globalConfig.server)
	u.Path = filepath.Join(pathPrefix, prefix, resource)
	for _, v := range args {
		u.Path = filepath.Join(u.Path, v)
	}

	reader := bytes.NewBuffer(body)
	req, err := http.NewRequest(method, u.String(), reader)
	if err != nil {
		return nil, fmt.Errorf("Could not build request: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("PORTUS-AUTH", fmt.Sprintf("%s:%s", globalConfig.user, globalConfig.token))
	return req, nil
}

func request(method, prefix, resource string, args []string, body []byte) (*http.Response, error) {
	req, err := buildRequest(method, prefix, resource, args, body)
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
		str := ""
		for k, list := range target["errors"] {
			str += k + ":\n"
			for _, e := range list {
				r, n := utf8.DecodeRuneInString(e)
				str += "  - " + string(unicode.ToUpper(r)) + e[n:] + "\n"
			}
		}
		return errors.New(strings.TrimRight(str, "\n"))
	} else if code == http.StatusUnauthorized || code == http.StatusForbidden {
		target := make(map[string]string)
		defer response.Body.Close()
		json.NewDecoder(response.Body).Decode(&target)
		return errors.New(target["error"])
	} else if code == http.StatusNotFound {
		return errors.New("Resource not found")
	} else if code == http.StatusMethodNotAllowed {
		return errors.New("Action not allowed on given resource")
	}
	return nil
}

func indexInSlice(ary []string, needle string) int {
	for i := 0; i < len(ary); i++ {
		if ary[i] == needle {
			return i
		}
	}
	return -1
}

func extractArguments(require argumentRequirement, args []string) (map[string]string, string, error) {
	id := ""
	values := make(map[string]string)
	unknown := make(map[string]string)
	for _, a := range args {
		keyValue := strings.Split(a, "=")
		i := indexInSlice(require.required, keyValue[0])
		if i >= 0 {
			if keyValue[0] == "id" {
				id = keyValue[1]
			} else {
				values[keyValue[0]] = keyValue[1]
			}
			require.required = append(require.required[:i], require.required[i+1:]...)
		} else {
			unknown[keyValue[0]] = keyValue[1]
		}
	}

	finalUnknown := []string{}
	for k, v := range unknown {
		i := indexInSlice(require.optional, k)
		if i >= 0 {
			values[k] = v
		} else {
			finalUnknown = append(finalUnknown, k)
		}
	}

	unkn := strings.Join(finalUnknown, ",")
	if unkn != "" {
		fmt.Printf("Ignoring the following keys: %v\n\n", unkn)
	}
	mandatory := strings.Join(require.required, ",")
	if mandatory == "" {
		path := ""
		if require.resource != "" && id != "" {
			path = filepath.Join(require.resource, id)
		}
		return values, path, nil
	}
	return nil, "", fmt.Errorf("The following mandatory fields are missing: %v", mandatory)
}

func generateBody(resource string, args []string) ([]byte, string, error) {
	require, ok := bodyArguments[resource]
	if !ok {
		return nil, "", fmt.Errorf("Resource '%v' does not support this action", resource)
	}

	values, path, err := extractArguments(require, args)
	if err != nil {
		return nil, path, err
	}

	var b []byte
	var e error

	if require.prefix == "" {
		b, e = json.Marshal(values)
	} else {
		b, e = json.Marshal(map[string]map[string]string{require.prefix: values})
	}
	if e != nil {
		return nil, path, fmt.Errorf("Internal error: ", e)
	}
	return b, path, nil
}

func create(resource string, args []string) error {
	b, path, err := generateBody(resource, args)
	if err != nil {
		return err
	}

	res, err := request("POST", path, resource, nil, b)
	if err != nil {
		return err
	}

	// TODO: quiet flag ?
	fmt.Printf("Created '%v' successfully!\n\n", bodyArguments[resource].name)

	body, _ := ioutil.ReadAll(res.Body)
	switch globalConfig.format {
	case jsonFmt:
		fmt.Print(string(body))
	default:
		return prettyPrint(bodyArguments[resource].returned, body, true)
	}
	return nil
}

func delete(resource string, args []string) error {
	prefix := bodyArguments[resource].resource
	if len(args) != 1 {
		return fmt.Errorf("Expecting exactly 1 argument, %v given", len(args))
	}

	res, err := request("DELETE", prefix, resource, args, nil)
	if err != nil {
		return err
	}
	body, _ := ioutil.ReadAll(res.Body)
	fmt.Printf("%v\n", string(body))

	// TODO: quiet flag ?
	fmt.Printf("Deleted '%v' successfully!\n", bodyArguments[resource].name)
	return nil
}

func parseArguments(resource string, args []string) (string, bool, error) {
	if len(args) > 1 {
		if val, ok := subresources[resource]; ok {
			resource = args[1]
			found := false
			for _, v := range val {
				if v == resource {
					found = true
					break
				}
			}
			if !found {
				return resource, false, fmt.Errorf("unknown subresource '%v'", resource)
			}
		} else {
			return resource, false, errors.New("too many given arguments")
		}
	}

	return resource, len(args)&1 == 1, nil
}

func get(resource string, args []string) error {
	rsrc, single, err := parseArguments(resource, args)
	if err != nil {
		return err
	}

	res, err := request("GET", "", resource, args, nil)
	if err != nil {
		return err
	}

	body, _ := ioutil.ReadAll(res.Body)
	switch globalConfig.format {
	case jsonFmt:
		fmt.Print(string(body))
	default:
		return prettyPrint(rsrc, body, single)
	}

	return nil
}

func update(resource string, args []string) error {
	return nil
}
