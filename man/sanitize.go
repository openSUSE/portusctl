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

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
)

// Returns false if printing the current line is a stupid decision (a.k.a. why
// would go-md2man behave like this...)
func canPrint(previous, next string) bool {
	if previous == "" {
		for _, v := range []string{"SH", "PP", "fi"} {
			if strings.HasPrefix(next, "."+v) {
				return false
			}
		}
		return next != ""
	}
	return true
}

func main() {
	if len(os.Args) < 1 {
		log.Fatal("Give me exactly one argument")
	}

	arg := os.Args[len(os.Args)-1]
	file, err := os.Open(arg)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	var previous, next string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		previous = next
		next = scanner.Text()

		if canPrint(previous, next) {
			fmt.Println(previous)
		}
	}
	fmt.Println(next)

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
}
