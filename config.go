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
	"errors"
	"fmt"

	"gopkg.in/urfave/cli.v1"
)

const (
	defaultFmt = iota
	jsonFmt
)

type config struct {
	user   string
	token  string
	server string
	format int64
}

// globalConfig contains the configuration needed to perform requests to the
// Portus server.
var globalConfig config

// availableResources contains the available resources with their
// shortcuts. This map can be accessed with the first character as the key.
var availableResources = map[string][]string{
	// TODO
	"app": []string{"application_tokens", "application_token", "at"}, // TODO: not on GET
	"nam": []string{"namespaces", "namespace", "n"},
	"rep": []string{"repositories", "repository", "r"},
	"reg": []string{"registries", "registry", "re"},
	"tag": []string{"tags", "tag", "ta"},
	"tea": []string{"teams", "team", "t"},
	"use": []string{"users", "user", "u"},
}

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
	switch ctx.String("format") {
	case "json":
		globalConfig.format = jsonFmt
	default:
		globalConfig.format = defaultFmt
	}
	return nil
}

// checkResource checks that the given resource is a valid one, and returns the
// identifier of the resource. If no resource was given, then the help command
// is executed.
func checkResource(resource string, ctx *cli.Context) (string, error) {
	if resource == "" {
		cli.ShowAppHelp(ctx)
		fmt.Println("")
		return "", errors.New("You have to provide at least one argument")
	}

	// TODO: broken one-letter shortcuts...
	// TODO: ./portusctl delete 3 <- CRASH
	bet := availableResources[resource[0:3]]
	for _, synonim := range bet {
		if synonim == resource {
			return bet[0], nil
		}
	}
	return "", fmt.Errorf("Unknown resource '%v'\n", resource)
}

// resourceDecorator decorates the given function with some checks on the
// arguments and flags.
func resourceDecorator(f func(string, []string) error) func(*cli.Context) error {
	return func(ctx *cli.Context) error {
		if err := setFlags(ctx); err != nil {
			return err
		}

		if resource, err := checkResource(ctx.Args().Get(0), ctx); err == nil {
			args := ctx.Args()
			return f(resource, args[1:len(args)])
		} else {
			return err
		}
	}
}
