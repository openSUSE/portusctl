#!/usr/bin/env bats -t
# Copyright (C) 2017-2019 Miquel Sabaté Solà <msabate@suse.com>
#
# This program is free software: you can redistribute it and/or modify
# it under the terms of the GNU General Public License as published by
# the Free Software Foundation, either version 3 of the License, or
# (at your option) any later version.
#
# This program is distributed in the hope that it will be useful,
# but WITHOUT ANY WARRANTY; without even the implied warranty of
# MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
# GNU General Public License for more details.
#
# You should have received a copy of the GNU General Public License
# along with this program.  If not, see <http://www.gnu.org/licenses/>.

load helpers

function setup() {
    __setup_db
    __source_environment
}

@test "exec works" {
    docker_run portusctl exec rake portus:info
    [[ "${lines[2]}" =~ "Evaluated configuration:" ]]
}

@test "changing the local flag fails because bundle is not installed there" {
    docker_run portusctl exec -l /srv rake portus:info
    [ $status -eq 1 ]
    # Since it's not in the local directory, it returns the global.
    [[ "${lines[1]}" =~ "\"bundle\": executable file not found" ]]
}

@test "changing the vendor flag results in error since bundler is not installed globally" {
    docker_run portusctl exec --vendor=false rake portus:info
    [ $status -eq 1 ]
    [[ "${lines[0]}" =~ "\"bundle\": executable file not found" ]]
}
