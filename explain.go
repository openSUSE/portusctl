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
	"fmt"
	"sort"
	"strings"

	"gopkg.in/urfave/cli.v1"
)

func enumerate(items []string) string {
	length := len(items)

	switch length {
	case 0:
		return ""
	case 1:
		return items[0]
	default:
		sort.Strings(items)
		subset := strings.Join(items[:length-1], ", ")
		return subset + " and " + items[length-1]
	}
}

func quote(items []string) []string {
	res := []string{}

	for _, v := range items {
		res = append(res, "'"+v+"'")
	}
	return res
}

func printSupportedCommands(res *Resource) {
	cmds := []string{}
	for _, v := range res.supportedActions {
		cmds = append(cmds, action(v))
	}

	fmt.Printf("Supported commands: %v\n", enumerate(cmds))
}

func printRequiredArguments(res *Resource) {
	var createRequirements, optional, validateRequirements string
	shouldPrint := false

	for _, v := range res.supportedActions {
		switch v {
		case postAction, putAction:
			shouldPrint = true
			createRequirements = enumerate(quote(res.required))
			optional = enumerate(quote(res.optional))
		case validateAction:
			shouldPrint = true
			validateRequirements = enumerate(quote(res.validate))
		}
	}

	if shouldPrint {
		fmt.Printf("Required parameters:\n")

		if createRequirements != "" {
			if optional == "" {
				fmt.Printf("  * In the 'create' and the 'update' commands: %v\n", createRequirements)
			} else {
				msg := "  * In the 'create' and the 'update' commands: %v (optional: %v)\n"
				fmt.Printf(msg, createRequirements, optional)
			}
		}
		if validateRequirements != "" {
			fmt.Printf("  * In the 'validate' command: %v\n", validateRequirements)
		}
	}
}

func printSubresources(res *Resource) {
	if len(res.subresources) == 0 {
		return
	}

	subs := enumerate(quote(res.subresources))
	fmt.Printf("\nSome commands will also accept the following 'subresources': %v\n", subs)

	sr := findResource(res.subresources[0])
	fmt.Printf("For example, one might perform:\n")
	fmt.Printf("\t$ portusctl get %v <id> %v\n", res.FullName(), sr.FullName())
}

func explain(res *Resource) error {
	fmt.Printf("Resource: %v %v\n", res.PluralName(), listSynonims("reference it with:", res.synonims))
	printSupportedCommands(res)
	printRequiredArguments(res)
	printSubresources(res)

	fmt.Printf("\nRefer to the man pages for usage examples\n")

	return nil
}

func explainAction(ctx *cli.Context) error {
	resource := ctx.Args().Get(0)

	res := findResource(resource)
	if res == nil {
		return resourceHelper(resource)
	}

	return explain(res)
}
