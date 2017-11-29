#!/usr/bin/env bats -t
# Copyright (C) 2017 Miquel Sabaté Solà <msabate@suse.com>
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

@test "health works" {
    portusctl health
    [ $status -eq 0 ]
    [[ "${lines[1]}" =~ 'Database is up-to-date' ]]
    [[ "${lines[1]}" =~ 'is reachable' ]]
}

@test "health does not accept further arguments" {
    portusctl health another
    [ $status -eq 1 ]
    [[ "${lines[0]}" =~ "you don't have to provide arguments for this command" ]]
}

@test "API user, server, token have to be provided" {
    unset PORTUSCTL_API_USER
    portusctl health
    [ $status -eq 1 ]
    [[ "${lines[-1]}" =~ "You have to set the user of the API" ]]

    export PORTUSCTL_API_USER="something"
    unset PORTUSCTL_API_TOKEN
    portusctl health
    [ $status -eq 1 ]
    [[ "${lines[-1]}" =~ "You have to set the token of your user" ]]

    export PORTUSCTL_API_TOKEN="something"
    unset PORTUSCTL_API_SERVER
    portusctl health
    [ $status -eq 1 ]
    [[ "${lines[-1]}" =~ "You have the deliver the URL of the Portus server" ]]
}
