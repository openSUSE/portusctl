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

# TODO
#@test "unsupported action" {
    # For now registries cannot be removed.
    #portusctl delete registries 1
    #[ $status -eq 1 ]
    #[[ "${lines[0]}" =~ "Resource 'registries' cannot be removed" ]]
#}

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
