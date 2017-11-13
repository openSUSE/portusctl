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

	"gopkg.in/urfave/cli.v1"
)

// TODO:
//
// export RAILS_ENV=production
// export GEM_PATH=/srv/Portus/vendor/bundle/ruby/2.1.0/
// export BUNDLER_BIN=/srv/Portus/vendor/bundle/ruby/2.1.0/bin/bundler.ruby2.1
// export RUBYLIB=$(ls -d /srv/Portus/vendor/bundle/ruby/2.1.0/gems/bundler*)/lib/

func vendorBasePath() string {
	// TODO: vendor flag
	return ""
}

func execCmd(c *cli.Context) error {
	if len(c.Args()) == 0 {
		return errors.New("you have to provide the command to be executed")
	}

	cmd := exec.Command("bundle", append([]string{"exec"}, c.Args()...)...)
	cmd.Dir = c.String("local")
	cmd.Env = append(os.Environ(),
		"FOO=duplicate_value", // TODO: bundle environment
	)
	cmd.Stdout = os.Stdout
	cmd.Stdin = os.Stdin
	cmd.Stderr = os.Stderr

	err := cmd.Run()
	return err
}
