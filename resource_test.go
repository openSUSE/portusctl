// Copyright (C) 2017-2019 Miquel Sabaté Solà <msabate@suse.com>
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

import "testing"

func TestResourceString(t *testing.T) {
	expected := [][]string{
		{"application token", "application_tokens"},
		{"namespace", "namespaces"},
		{"repository", "repositories"},
		{"registry", "registries"},
		{"tag", "tags"},
		{"team", "teams"},
		{"user", "users"},
	}

	if len(availableResources) != len(expected) {
		t.Errorf("You have to update me :)")
	}

	for k, v := range availableResources {
		if v.String() != expected[k][0] {
			t.Errorf("Expected: '%v' - Given: '%v'", expected[k][0], v.String())
		}
		if v.FullName() != expected[k][1] {
			t.Errorf("Expected: '%v' - Given: '%v'", expected[k][1], v.FullName())
		}
	}
}

func TestResourcePath(t *testing.T) {
	res := availableResources[0]
	path := res.Path([]string{"path", "1"})
	if path != "/api/v1/application_tokens/path/1" {
		t.Errorf("Bad path for '%v'", path)
	}

	res.prefix = "prefix/1"
	path = res.Path([]string{"path", "1"})
	if path != "/api/v1/prefix/1/application_tokens/path/1" {
		t.Errorf("Bad path for '%v'", path)
	}

	res.prefix = ""
	path = res.Path(nil)
	if path != "/api/v1/application_tokens" {
		t.Errorf("Bad path for '%v'", path)
	}
}

func TestFindResource(t *testing.T) {
	cases := []struct {
		given        string
		expectedKind int
	}{
		// Application tokens
		{"at", kindApplicationToken},
		{"application_token", kindApplicationToken},
		{"application_tokens", kindApplicationToken},

		// Namespaces
		{"n", kindNamespace},
		{"namespace", kindNamespace},
		{"namespaces", kindNamespace},

		// Repository
		{"r", kindRepository},
		{"repository", kindRepository},
		{"repositories", kindRepository},

		// Registry
		{"re", kindRegistry},
		{"registry", kindRegistry},
		{"registries", kindRegistry},

		// Tag
		{"tag", kindTag},
		{"tags", kindTag},

		// Team
		{"t", kindTeam},
		{"team", kindTeam},
		{"teams", kindTeam},

		// User
		{"u", kindUser},
		{"user", kindUser},
		{"users", kindUser},

		// Invalid
		{"", -1},
		{"k", -1},
		{"mssola", -1},
		{"1", -1},
	}

	for _, v := range cases {
		res := findResource(v.given)
		if res == nil {
			if v.expectedKind != -1 {
				t.Errorf("CASE: '%v' returned nil when it shouldn't", v.given)
			}
		} else if res.kind != v.expectedKind {
			t.Errorf("CASE '%v' returned '%v' but '%v' was expected", v.given, res.kind, v.expectedKind)
		}
	}
}
