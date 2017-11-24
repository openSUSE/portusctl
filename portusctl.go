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
	"fmt"

	"gopkg.in/urfave/cli.v1"
)

func main() {
	app := cli.NewApp()
	app.Name = "portusctl"
	app.Usage = "Client for your Portus instance"
	app.UsageText = "portusctl <command> [arguments...]"
	app.HideHelp = true
	app.Version = versionString()

	app.CommandNotFound = func(context *cli.Context, cmd string) {
		fmt.Printf("Incorrect usage: command '%v' does not exist.\n\n", cmd)
		cli.ShowAppHelp(context)
	}

	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:   "server, s",
			Usage:  "The location where the Portus instance serves requests",
			EnvVar: "PORTUSCTL_API_SERVER",
		},
		cli.StringFlag{
			Name:   "token, t",
			Usage:  "The authentication token of the user for the Portus REST API",
			EnvVar: "PORTUSCTL_API_TOKEN",
		},
		cli.StringFlag{
			Name:   "user, u",
			Usage:  "The user of the Portus REST API",
			EnvVar: "PORTUSCTL_API_USER",
		},
		cli.BoolFlag{
			Name:   "quiet, q",
			Usage:  "Prevent portusctl from outputting general information",
			EnvVar: "PORTUSCTL_QUIET",
		},
	}

	app.Commands = []cli.Command{
		{
			Name:   "create",
			Usage:  "Create the given resource",
			Action: resourceDecorator(create, postAction),
			ArgsUsage: `<resource> [arguments...]

Where <resource> is the resource that you want to create.`,
		},
		{
			Name:   "delete",
			Usage:  "Delete the given resource",
			Action: resourceDecorator(delete, deleteAction),
			ArgsUsage: `<resource> [arguments...]

Where <resource> is the resource that you want to delete.`,
		},
		{
			Name:   "exec",
			Usage:  "Execute an arbitrary command on the environment of your Portus instance",
			Action: execCmd,
			ArgsUsage: `<command> [arguments...]

Where <command> is the command that you want to run on the environment of your
Portus instance. The successive arguments will be passed also to this command.`,
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:   "local, l",
					Value:  defaultPath,
					Usage:  "The location on the current host of your Portus instance",
					EnvVar: "PORTUSCTL_EXEC_LOCATION",
				},
				cli.BoolFlag{
					Name:   "vendor, v",
					Usage:  "Use the local 'vendor' directory as the gem environment",
					EnvVar: "PORTUSCTL_EXEC_VENDOR",
				},
			},
		},
		{
			Name:   "get",
			Usage:  "Fetches info for the given resource",
			Action: resourceDecorator(get, getAction),
			ArgsUsage: `<resource> [arguments...]

Where <resource> is the resource that you want to fetch.`,
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:   "format, f",
					Value:  "default",
					Usage:  "The output format. Available options: default and json",
					EnvVar: "PORTUSCTL_FORMAT",
				},
			},
		},
		{
			Name:      "health",
			Usage:     "Get health info from Portus",
			Action:    healthAction,
			ArgsUsage: " ",
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:   "format, f",
					Value:  "default",
					Usage:  "The output format. Available options: default and json",
					EnvVar: "PORTUSCTL_FORMAT",
				},
			},
		},
		{
			Name:   "update",
			Usage:  "Update the given resource",
			Action: resourceDecorator(update, putAction),
			ArgsUsage: `<resource> [arguments...]

Where <resource> is the resource that you want to update.`,
		},
		{
			Name:   "validate",
			Usage:  "Validate the given resource",
			Action: resourceDecorator(validate, validateAction),
			ArgsUsage: `<resource> [arguments...]

Where <resource> is the resource that you want to validate.`,
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:   "format, f",
					Value:  "default",
					Usage:  "The output format. Available options: default and json",
					EnvVar: "PORTUSCTL_FORMAT",
				},
			},
		},
		{
			Name:      "version",
			Usage:     "Print the client and server version information",
			Action:    versionAction,
			ArgsUsage: " ",
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:   "format, f",
					Value:  "default",
					Usage:  "The output format. Available options: default and json",
					EnvVar: "PORTUSCTL_FORMAT",
				},
			},
		},
	}

	app.RunAndExitOnError()
}
