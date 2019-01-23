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

@test "explain show the possibilities when empty" {
    portusctl explain
    [ $status -eq 1 ]
    [[ "${lines[0]}" =~ "You must specify one of the following types of resource:" ]]
    [[ "${lines[1]}" =~ "* application_tokens (aka 'at' or 'application_token')" ]]
}

@test "explain actually explains commands" {
    portusctl explain u
    [ $status -eq 0 ]

    [[ "${lines[0]}" =~ "Resource: Users (reference it with: 'u' or 'user' or 'users')" ]]
    [[ "${lines[1]}" =~ "Supported commands: create, delete, get and update" ]]
    [[ "${lines[2]}" =~ "Required parameters:" ]]
    [[ "${lines[3]}" =~ "* In the 'create' and the 'update' commands: 'email', 'password' and 'username' (optional: 'display_name')" ]]
    [[ "${lines[4]}" =~ "Some commands will also accept the following 'subresources': 'application_tokens'" ]]
    [[ "${lines[5]}" =~ "For example, one might perform:" ]]
    [[ "${lines[6]}" =~ "$ portusctl get users <id> application_tokens" ]]
    [[ "${lines[7]}" =~ "Refer to the man pages for usage examples" ]]
}

@test "explain also explains validate parameters" {
    portusctl explain namespace
    [ $status -eq 0 ]

    [[ "${lines[3]}" =~ "* In the 'create' and the 'update' commands: 'name' and 'team' (optional: 'description')" ]]
    [[ "${lines[4]}" =~ "* In the 'validate' command: 'name'" ]]
}
