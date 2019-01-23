// Copyright (C) 2018-2019 Miquel Sabaté Solà <msabate@suse.com>
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
	"path/filepath"

	"gopkg.in/urfave/cli.v1"
)

var retainedValues map[string]string

func bootstrapCmd(c *cli.Context) error {
	if err := setFlags(c, true); err != nil {
		return err
	}

	retainedValues = make(map[string]string)

	resource := findResourceByID(kindUser)
	b, err := generateBody(resource, c.Args(), true)
	if err == nil {
		path := filepath.Join(resource.Path(nil), "bootstrap")
		res, err := request("POST", path, "", b)
		if err != nil {
			return err
		}
		return printAndQuit(res, kindBootstrap, true)
	}
	return err
}
