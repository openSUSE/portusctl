// Copyright (C) 2018 Miquel Sabaté Solà <msabate@suse.com>
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

func TestEnumerate(t *testing.T) {
	res := enumerate([]string{})
	if res != "" {
		t.Fatalf("Expecting to be empty")
	}

	res = enumerate([]string{"a"})
	if res != "a" {
		t.Fatalf("Expecting to be 'a', '%v' given", res)
	}

	res = enumerate([]string{"a", "b"})
	if res != "a and b" {
		t.Fatalf("Expecting to be 'a and b', '%v' given", res)
	}

	res = enumerate([]string{"a", "b", "c"})
	if res != "a, b and c" {
		t.Fatalf("Expecting to be 'a, b and c', '%v' given", res)
	}
}
