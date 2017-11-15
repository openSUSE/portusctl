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
	"testing"

	"github.com/mssola/capture"
)

func expectParseArguments(t *testing.T, resource *Resource, args []string, kind int, single bool, err string) {
	r, s, e := parseArguments(resource, args)

	if e == nil && err != "" {
		t.Fatalf("Expected an error but none was given")
	} else if e != nil {
		if e.Error() != err {
			t.Fatalf("'%v' != '%v'", e.Error(), err)
		}
		return
	}

	if r.kind != kind {
		t.Fatalf("Expected to be %v, but %v was given", kind, r.kind)
	}
	if s != single {
		t.Fatalf("Expected to be %v", single)
	}
}

// The extractArguments function destroys some values from the original array,
// so instead of using the global variable directly, we'll use hard copies of it.
func newResource(kind int) Resource {
	resource := findResourceByID(kind)
	r := Resource{
		bundler:       resource.bundler,
		kind:          resource.kind,
		prefix:        resource.prefix,
		returnedKind:  resource.returnedKind,
		superresource: resource.superresource,
	}
	r.required = append([]string{}, resource.required...)
	r.synonims = append([]string{}, resource.synonims...)
	r.optional = append([]string{}, resource.optional...)
	r.subresources = append([]string{}, resource.subresources...)
	return r
}

func testMap(t *testing.T, args map[string]string, values [][]string) {
	if len(values) != len(args) {
		t.Fatalf("Wrong length: given values '%v', expected '%v'", len(values), len(args))
	}

	for _, v := range values {
		given, ok := args[v[0]]
		if !ok {
			t.Fatalf("Key '%v' does not exist!", v[0])
		}

		if given != v[1] {
			t.Fatalf("Got '%v', expected '%v'; for key '%v'", given, v[1], v[0])
		}
	}
}

func TestParseArguments(t *testing.T) {
	user := findResourceByID(kindUser)

	expectParseArguments(t, user, nil, kindUser, false, "")
	expectParseArguments(t, user, []string{"1"}, kindUser, true, "")
	expectParseArguments(t, user, []string{"1", "at"}, kindApplicationToken, false, "")
	expectParseArguments(t, user, []string{"1", "at", "1"}, kindApplicationToken, true, "")
	expectParseArguments(t, user, []string{"1", "1"}, kindApplicationToken, true, "unknown subresource '1'")

	at := findResourceByID(kindApplicationToken)
	expectParseArguments(t, at, []string{"1", "at"}, kindApplicationToken, true, "too many arguments")
}

func TestExtractArgumentsUsers(t *testing.T) {
	var err error

	user := newResource(kindUser)
	user.action = postAction

	captured := capture.All(func() {
		_, err = extractArguments(&user, []string{"name=msabate", "email=lala@example.org"})
	})
	msg := "The following mandatory fields are missing: username, password"
	if err == nil {
		t.Fatalf("Expected '%v', got nil", msg)
	}
	if err.Error() != msg {
		t.Fatalf("Expected '%v', got '%v'", msg, err.Error())
	}
	msg = "Ignoring the following keys: name\n\n"
	if string(captured.Stdout) != msg {
		t.Fatalf("Expected '%v', got '%v'", msg, string(captured.Stdout))
	}
}

func TestExtractArgumentsUsersOK(t *testing.T) {
	var args map[string]string
	var err error

	user := newResource(kindUser)
	user.action = postAction

	captured := capture.All(func() {
		args, err = extractArguments(&user,
			[]string{"username=msabate", "email=lala@example.org", "password=12341234"})
	})
	if string(captured.Stdout) != "" {
		t.Fatalf("Expected no message printed, got '%v'", string(captured.Stdout))
	}
	if err != nil {
		t.Fatalf("Expected no error, got '%v'", err.Error())
	}

	testMap(t, args, [][]string{
		{"username", "msabate"},
		{"email", "lala@example.org"},
		{"password", "12341234"},
	})
}

func TestExtractArgumentsUsersOptional(t *testing.T) {
	var args map[string]string
	var err error

	user := newResource(kindUser)
	user.action = postAction

	captured := capture.All(func() {
		args, err = extractArguments(&user,
			[]string{"username=msabate", "email=lala@example.org",
				"password=12341234", "display_name=miquel"})
	})
	if string(captured.Stdout) != "" {
		t.Fatalf("Expected no message printed, got '%v'", string(captured.Stdout))
	}
	if err != nil {
		t.Fatalf("Expected no error, got '%v'", err.Error())
	}

	testMap(t, args, [][]string{
		{"username", "msabate"},
		{"email", "lala@example.org"},
		{"password", "12341234"},
		{"display_name", "miquel"},
	})
}

func TestExtractArgumentsWithID(t *testing.T) {
	var args map[string]string
	var err error

	at := newResource(kindApplicationToken)
	at.action = postAction

	captured := capture.All(func() {
		args, err = extractArguments(&at, []string{"application=lala", "id=2"})
	})
	if string(captured.Stdout) != "" {
		t.Fatalf("Expected no message printed, got '%v'", string(captured.Stdout))
	}
	if err != nil {
		t.Fatalf("Expected no error, got '%v'", err.Error())
	}

	testMap(t, args, [][]string{{"application", "lala"}})

	if at.prefix != "users/2" {
		t.Fatalf("Wrong prefix: expected 'users/2'; got '%v'", at.prefix)
	}
}

func TestExtractArgumentsNoArgumentsToPut(t *testing.T) {
	var err error

	user := newResource(kindUser)
	user.action = putAction

	captured := capture.All(func() {
		_, err = extractArguments(&user, []string{"name=msabate"})
	})
	msg := "Ignoring the following keys: name\n\n"
	if string(captured.Stdout) != msg {
		t.Fatalf("Expected '%v', got '%v'", msg, string(captured.Stdout))
	}

	msg = "You have to provide at least one of the following arguments: username, email, password, display_name"
	if err.Error() != msg {
		t.Fatalf("Expected '%v', got '%v'", msg, err.Error())
	}
}

func TestExtractArgumentsPut(t *testing.T) {
	var args map[string]string
	var err error

	user := newResource(kindUser)
	user.action = putAction

	captured := capture.All(func() {
		args, err = extractArguments(&user, []string{"username=msabate"})
	})
	if string(captured.Stdout) != "" {
		t.Fatalf("Expected no message printed, got '%v'", string(captured.Stdout))
	}
	if err != nil {
		t.Fatalf("Expected no error, got '%v'", err.Error())
	}

	testMap(t, args, [][]string{{"username", "msabate"}})
}
