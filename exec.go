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
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"gopkg.in/urfave/cli.v1"
)

const (
	// defaultPath is the path taken by default as the server location. This is
	// the path where Portus will be located in the Docker image, so it's
	// probably not a good idea to move away from this default value.
	defaultPath = "/srv/Portus"
)

// rubyVersion returns the current ruby version.
func rubyVersion(local string) (string, error) {
	path := filepath.Join(local, ".ruby-version")
	contents, err := ioutil.ReadFile(path)
	if err != nil {
		return "", err
	}

	return strings.TrimSpace(string(contents)), nil
}

// Returns both the ruby version and the base path for bundled ruby gems.
func pathBase(dir string) (string, string, error) {
	ruby, err := rubyVersion(dir)
	if err != nil {
		return "", "", err
	}

	base := filepath.Join(dir, "vendor/bundle/ruby", ruby)
	return ruby, base, nil
}

// vendorRubyPath returns some paths in the form of environment variables
// targeting a local `vendor` directory of the Portus instance.
func vendorRubyPath(dir string) ([]string, error) {
	_, base, err := pathBase(dir)
	if err != nil {
		return nil, err
	}

	lib, err := filepath.Glob(filepath.Join(base, "gems", "bundler*"))
	if err != nil {
		return nil, err
	}

	return []string{
		"RAILS_ENV=production",
		"GEM_PATH=" + base,
		"RUBYLIB=" + filepath.Join(lib[0], "lib"),
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

// Returns the `bundle` command to be used. The command being returned depends
// on whether the vendor flag was set to true (default) or not. If the vendor
// flag is set but no local bundler binary could be found, then the global
// version is returned.
func bundleCommand(ctx *cli.Context) string {
	if !ctx.Bool("vendor") {
		return "bundle"
	}

	ruby, base, err := pathBase(ctx.String("local"))
	if err != nil {
		log.Printf("Could not find local bundler: %v", err.Error())
		return "bundle"
	}
	return filepath.Join(base, "bin/bundler.ruby"+ruby[:3])
}

func execCmd(c *cli.Context) error {
	if len(c.Args()) == 0 {
		return errors.New("you have to provide the command to be executed")
	}

	cmd := exec.Command(bundleCommand(c), append([]string{"exec"}, c.Args()...)...)
	cmd.Dir = c.String("local")
	cmd.Env = environment(c)
	cmd.Stdout = os.Stdout
	cmd.Stdin = os.Stdin
	cmd.Stderr = os.Stderr

	err := cmd.Run()
	return err
}
