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

// File taken from openSUSE/umoci.

package main

import (
	"os"
	"strings"
	"testing"
)

// Test is a hack that allows us to figure out what the coverage is during
// integration tests. I would not recommend that you use a binary built using
// this hack outside of a test suite. Read the following link for more info:
//   https://www.cyphar.com/blog/post/20170412-golang-integration-coverage
func TestMain(t *testing.T) {
	var args []string
	run := os.Getenv("DOCKER_DEVEL_COVER_TESTS") != ""

	for _, arg := range os.Args {
		switch {
		case arg == "__DEVEL--cover-tests":
			run = true
		case strings.HasPrefix(arg, "-test"):
		case strings.HasPrefix(arg, "__DEVEL"):
		default:
			args = append(args, arg)
		}
	}
	os.Args = args

	if run {
		main()
	}
}
