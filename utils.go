// Copyright (C) 2017-2018 Miquel Sabaté Solà <msabate@suse.com>
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

// indexInSlice returns the index of the given needle inside of the given
// slice. If this is not possible, then it returns -1.
func indexInSlice(ary []string, needle string) int {
	for i := 0; i < len(ary); i++ {
		if ary[i] == needle {
			return i
		}
	}
	return -1
}
