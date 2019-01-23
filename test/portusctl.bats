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

@test "server could not be found" {
    portusctl -s "http://localhost:1234" get users
    [ $status -eq 1 ]
    [[ "${lines[0]}" =~ "connection refused" ]]
}

@test "wrong credentials" {
    portusctl -t "asdasdasd" get users
    [ $status -eq 1 ]
    [[ "${lines[0]}" =~ "Authentication fails" ]]
}

@test "unknown command" {
    # It's deleted, not destroyed
    portusctl destroy registries 1
    [ $status -eq 0 ]
    [[ "${lines[0]}" =~ "Incorrect usage: command 'destroy' does not exist." ]]
}

@test "unsupported action" {
    # For now registries cannot be removed.
    portusctl delete registries 1
    [ $status -eq 1 ]
    [[ "${lines[0]}" =~ "Action not supported for resource 'registry'" ]]
}

@test "user does not have permission to perform some operation" {
    # Let's create a non-admin and try to perform an admin-only action with the
    # newly created credentials.
    portusctl create users username=msabate email=example@test.lan password=12341234
    portusctl create application_tokens id=2 application=something
    token="$(echo -e "${lines[2]}" | tr -d '[:space:]')"

    portusctl -u msabate -t "$token" get users
    [ $status -eq 1 ]
    [[ "${lines[0]}" =~ "Authorization fails" ]]
}

@test "quiet flag" {
    portusctl -q create user username=msabate lala=this email=example@test.lan password=12341234
    [ $status -eq 0 ]

    # Note that we no longer have a first line being "Created ... successfully"
    # Moreover, it doesn't complain about the useless `lala` argument.
    [[ "${lines[0]}" =~ "ID    Username    Email               Admin    NamespaceID    DisplayName" ]]
    [[ "${lines[1]}" =~ "2     msabate     example@test.lan    false    3" ]]

    # No lines from deleting a resource (well, two because of the two cover sentences...)
    portusctl -q delete user 2
    [ $status -eq 0 ]
    [[ "${#lines[@]}" -eq 2 ]]
}

@test "API user, server, token have to be provided" {
    unset PORTUSCTL_API_USER
    portusctl get users
    [ $status -eq 1 ]
    [[ "${lines[-1]}" =~ "You have to set the user of the API" ]]

    export PORTUSCTL_API_USER="something"
    unset PORTUSCTL_API_TOKEN
    portusctl get users
    [ $status -eq 1 ]
    [[ "${lines[-1]}" =~ "You have to set the token of your user" ]]

    export PORTUSCTL_API_TOKEN="something"
    unset PORTUSCTL_API_SERVER
    portusctl get users
    [ $status -eq 1 ]
    [[ "${lines[-1]}" =~ "You have the deliver the URL of the Portus server" ]]
}

@test "Error on unknown resource" {
    portusctl get whatever
    [ $status -eq 1 ]
    [[ "${lines[0]}" =~ "Unknown resource 'whatever'" ]]
}

@test "Error on unsupported action" {
    portusctl get at
    [ $status -eq 1 ]
    [[ "${lines[0]}" =~ "Action not supported for resource 'application token'" ]]

}
