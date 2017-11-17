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
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"

	"gopkg.in/urfave/cli.v1"
)

// TODO: mssola: the whole `vendor` feature needs proper integration
// testing. None has been provided because the development docker image being
// used is different on this regard. I'll provide integration tests once we fix
// the production docker image.

const (
	defaultPath = "/srv/Portus"
)

var rubyRegexp = regexp.MustCompile(`^ruby\s+(\d+)\.(\d+)\.(\d+).+`)

// getOutputFromCommand runs a command under the `bundle exec` environment, and
// at the directory specified by the dir parameter. Returns the combined output
// if successful, otherwise it returns an error.
func getOutputFromCommand(args []string, dir string) (string, error) {
	cmd := exec.Command("bundle", append([]string{"exec"}, args...)...)
	cmd.Dir = dir
	b, err := cmd.CombinedOutput()
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(b)), nil
}

// rubyVersion returns the current ruby version.
func rubyVersion(local string) (string, error) {
	output, err := getOutputFromCommand([]string{"ruby", "-v"}, local)
	if err != nil {
		return "", err
	}

	match := rubyRegexp.FindStringSubmatch(output)
	if len(match) != 4 {
		return "", errors.New("something went wrong")
	}
	return strings.Join(match[1:], "."), nil
}

// vendorRubyPath returns some paths in the form of environment variables
// targeting a local `vendor` directory of the Portus instance.
func vendorRubyPath(dir string) ([]string, error) {
	ruby, err := rubyVersion(dir)
	if err != nil {
		return nil, err
	}

	base := filepath.Join(dir, "vendor/bundle/ruby", ruby)

	cmd := []string{"ls", "-d", filepath.Join(base, "bundler*")}
	lib, err := getOutputFromCommand(cmd, dir)
	if err != nil {
		return nil, err
	}

	return []string{
		"GEM_PATH=" + base,
		"BUNDLER_BIN=" + filepath.Join(base, "bin/bundler.ruby"+ruby[:3]),
		"RUBYLIB=" + filepath.Join(lib, "lib"),
	}, nil
}

// environment returns the environment variables to be set for the main `exec`
// command.
func environment(ctx *cli.Context) []string {
	environment := os.Environ()

	if ctx.Bool("vendor") {
		dir := ctx.String("local")
		if vendor, err := vendorRubyPath(dir); err == nil {
			environment = append(environment, vendor...)
		}
	}
	return environment
}

func execCmd(c *cli.Context) error {
	if len(c.Args()) == 0 {
		return errors.New("you have to provide the command to be executed")
	}

	cmd := exec.Command("bundle", append([]string{"exec"}, c.Args()...)...)
	cmd.Dir = c.String("local")
	cmd.Env = environment(c)
	cmd.Stdout = os.Stdout
	cmd.Stdin = os.Stdin
	cmd.Stderr = os.Stderr

	err := cmd.Run()
	return err
}
