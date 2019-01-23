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

func TestIndexInSlice(t *testing.T) {
	if indexInSlice([]string{"a", "b"}, "a") != 0 {
		t.Fatalf("Wrong index")
	}
	if indexInSlice([]string{"a", "b"}, "b") != 1 {
		t.Fatalf("Wrong index")
	}
	if indexInSlice([]string{"a", "b"}, "c") != -1 {
		t.Fatalf("Wrong index")
	}
}
