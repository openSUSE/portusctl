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
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"text/tabwriter"

	"gopkg.in/urfave/cli.v1"
)

var gitCommit, version string

func versionString() string {
	str := version

	if gitCommit != "" {
		str += fmt.Sprintf(" with commit '%v'", gitCommit)
	}
	return fmt.Sprintf(`%v.
Copyright (C) 2017 Miquel Sabaté Solà <msabate@suse.com>
License GPLv3+: GNU GPL version 3 or later "<http://gnu.org/licenses/gpl.html>.
This is free software: you are free to change and redistribute it.
There is NO WARRANTY, to the extent permitted by law.`, str)
}

func printVersion(version Version) error {
	w := tabwriter.NewWriter(os.Stdout, 0, 4, 4, ' ', tabwriter.TabIndent)
	top := "SERVER VERSION\t"
	bottom := version.Version + "\t"

	git := version.Git
	if git.Tag != "" {
		top += "SERVER TAG\t"
		bottom += git.Tag + "\t"
	} else if git.Branch != "" && git.Commit != "" {
		top += "SERVER BRANCH\t"
		bottom += fmt.Sprintf("%v@%v\t", git.Branch, git.Commit)
	}

	top += "CLIENT VERSION\t"
	bottom += version.PortusctlVersion + "\t"

	apis := strings.Join(version.APIVersions, ", ")
	top += "SERVER API\tCLIENT API"
	bottom += fmt.Sprintf("%v\tv1", apis)

	fmt.Fprintln(w, top)
	fmt.Fprintln(w, bottom)
	return w.Flush()
}

func combineOutput(res *http.Response) error {
	data := Version{}

	dec := json.NewDecoder(res.Body)
	err := dec.Decode(&data)
	if err != nil {
		return err
	}

	data.PortusctlVersion = version
	if globalConfig.format == jsonFmt {
		enc := json.NewEncoder(os.Stdout)
		return enc.Encode(&data)
	}

	return printVersion(data)
}

func versionAction(ctx *cli.Context) error {
	if len(ctx.Args()) != 0 {
		return errors.New("you don't have to provide arguments for this command")
	}
	if err := setFlags(ctx); err != nil {
		return err
	}

	path := filepath.Join(apiPrefix, "version")
	res, err := request("GET", path, "", nil)
	if err != nil {
		return err
	}

	return combineOutput(res)
}
