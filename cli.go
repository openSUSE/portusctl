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
	"errors"
	"fmt"
	"path/filepath"
	"strings"

	"gopkg.in/urfave/cli.v1"
)

const (
	defaultFmt = iota
	jsonFmt
)

const (
	getAction = iota
	postAction
	putAction
	deleteAction
	validateAction
)

// Returns the string value of the given action identifier.
func action(act int) string {
	switch act {
	case getAction:
		return "get"
	case postAction:
		return "create"
	case putAction:
		return "update"
	case deleteAction:
		return "delete"
	case validateAction:
		return "validate"
	default:
		return "<unknown>"
	}
}

type config struct {
	user   string
	token  string
	server string
	format int
	quiet  bool
}

// globalConfig contains the configuration needed to perform requests to the
// Portus server.
var globalConfig config

// setFlags sets the global configuration with the values provided by the
// flags. An error will be returned when one of the mandatory flags is not
// provided.
func setFlags(ctx *cli.Context) error {
	if globalConfig.user = ctx.GlobalString("user"); globalConfig.user == "" {
		cli.ShowAppHelp(ctx)
		fmt.Println("")
		return errors.New("You have to set the user of the API")
	}
	if globalConfig.token = ctx.GlobalString("token"); globalConfig.token == "" {
		cli.ShowAppHelp(ctx)
		fmt.Println("")
		return errors.New("You have to set the token of your user")
	}
	if globalConfig.server = ctx.GlobalString("server"); globalConfig.server == "" {
		cli.ShowAppHelp(ctx)
		fmt.Println("")
		return errors.New("You have the deliver the URL of the Portus server")
	}

	globalConfig.quiet = ctx.GlobalBool("quiet")

	switch ctx.String("format") {
	case "json":
		globalConfig.format = jsonFmt
	default:
		globalConfig.format = defaultFmt
	}
	return nil
}

// parseArguments is only valid for the `get` command and parses the arguments
// for the given initial resource. It returns the real resource being targeted
// (e.g. it can be a subresource), whether we are fetching one or multiple
// elements (`true` for a single element, false otherwise) and an error if
// possible.
func parseArguments(resource *Resource, args []string) (*Resource, bool, error) {
	if len(args) > 1 {
		if len(resource.subresources) > 0 {
			resource = findResource(args[1])
			if resource == nil {
				return nil, false, fmt.Errorf("unknown subresource '%v'", args[1])
			}
			args[1] = resource.FullName()
		} else {
			return nil, false, errors.New("too many arguments")
		}
	}

	return resource, len(args)&1 == 1, nil
}

func extractArguments(resource *Resource, args []string, validate bool) (map[string]string, error) {
	id := ""
	values := make(map[string]string)
	unknown := make(map[string]string)

	if validate {
		resource.required = resource.validate
		resource.optional = []string{}
	}

	for _, a := range args {
		keyValue := strings.Split(a, "=")
		i := indexInSlice(resource.required, keyValue[0])
		if i >= 0 {
			if keyValue[0] == "id" {
				id = keyValue[1]
			} else {
				values[keyValue[0]] = keyValue[1]
			}
			resource.required = append(resource.required[:i], resource.required[i+1:]...)
		} else if len(keyValue) == 2 {
			unknown[keyValue[0]] = keyValue[1]
		}
	}

	finalUnknown := []string{}
	for k, v := range unknown {
		i := indexInSlice(resource.optional, k)
		if i >= 0 {
			values[k] = v
		} else {
			finalUnknown = append(finalUnknown, k)
		}
	}

	if len(finalUnknown) != 0 && !globalConfig.quiet {
		fmt.Printf("Ignoring the following keys: %v\n\n", strings.Join(finalUnknown, ", "))
	}
	if len(resource.required) != 0 {
		if resource.action == postAction {
			return nil, fmt.Errorf("The following mandatory fields are missing: %v",
				strings.Join(resource.required, ", "))
		}
		if len(values) == 0 {
			arguments := append(resource.required, resource.optional...)
			return nil, fmt.Errorf(
				"You have to provide at least one of the following arguments: %v",
				strings.Join(arguments, ", "))
		}
	}
	super := findResourceByID(resource.superresource)
	if super != nil && id != "" {
		resource.prefix = filepath.Join(super.FullName(), id)
	}
	return values, nil
}

// Returns a string like this: "(<prefix> '<element1>' or '<element2>' or ...)";
// where <element1> and <element2> are elements from the given synonims.
func listSynonims(prefix string, synonims []string) string {
	out := []string{}

	for i := 0; i < len(synonims); i++ {
		out = append(out, "'"+synonims[i]+"'")
	}

	return "(" + prefix + " " + strings.Join(out[:], " or ") + ")"
}

// resourceHelper returns an error containing the proper message when a resource
// was not given.
func resourceHelper(given string) error {
	if given != "" {
		given = fmt.Sprintf("Unknown resource '%v'. ", given)
	}
	given += "You must specify one of the following types of resource:\n\n"

	for _, resource := range availableResources {
		syn := resource.synonims[0 : len(resource.synonims)-1]
		given += "    * " + resource.FullName() + " " + listSynonims("aka", syn) + "\n"
	}
	given += "\nSee the man pages for help and examples."
	return errors.New(given)
}

// checkResource checks that the given resource is a valid one, and returns the
// identifier of the resource. If no resource was given, then the help command
// is executed.
func checkResource(resource string, ctx *cli.Context, action int) (*Resource, error) {
	res := findResource(resource)
	if res == nil {
		return nil, resourceHelper(resource)
	}

	// Check that the given action can be perform on the resource.
	err := res.SetAction(action)
	if err != nil {
		return nil, err
	}
	return res, nil
}

// resourceDecorator decorates the given function with some checks on the
// arguments and flags.
func resourceDecorator(f func(*Resource, []string) error, action int) func(*cli.Context) error {
	return func(ctx *cli.Context) error {
		if err := setFlags(ctx); err != nil {
			return err
		}

		resource, err := checkResource(ctx.Args().Get(0), ctx, action)
		if err != nil {
			return err
		}

		args := ctx.Args()
		return f(resource, args[1:])
	}
}
