#!/usr/bin/env bats -t
# Copyright (C) 2017-2018 Miquel Sabaté Solà <msabate@suse.com>
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

@test "create two users, fetch them, then delete one and fetch again" {
    portusctl create user username=msabate email=example@test.lan password=12341234
    [ $status -eq 0 ]
    portusctl create user username=another email=another@test.lan password=12341234 display_name=lala
    [ $status -eq 0 ]

    portusctl get users
    [ $status -eq 0 ]
    [[ "${lines[0]}" =~ "ID    Username    Email                  Admin    NamespaceID    DisplayName" ]]
    [[ "${lines[1]}" =~ "1     admin       admin@example.local    true     2" ]]
    [[ "${lines[2]}" =~ "2     msabate     example@test.lan       false    3" ]]
    [[ "${lines[3]}" =~ "3     another     another@test.lan       false    4              lala" ]]

    portusctl delete user 2
    [ $status -eq 0 ]

    portusctl get u
    [ $status -eq 0 ]
    [[ "${lines[0]}" =~ "ID    Username    Email                  Admin    NamespaceID    DisplayName" ]]
    [[ "${lines[1]}" =~ "1     admin       admin@example.local    true     2" ]]
    [[ "${lines[2]}" =~ "3     another     another@test.lan       false    4              lala" ]]
}

@test "create an user which has already been taken" {
    portusctl create users username=msabate email=example@test.lan password=12341234
    [ $status -eq 0 ]
    portusctl create users username=msabate email=example@test.lan password=12341234
    [ $status -eq 1 ]

    # NOTE: I'm not checking the "email" and the "username" lines themselves
    # because they might come unsorted, so it's not really safe.
    [[ "${lines[1]}" =~ "- Has already been taken" ]]
    [[ "${lines[3]}" =~ "- Has already been taken" ]]
}

@test "fetch a single user" {
    portusctl get user 1
    [ $status -eq 0 ]
    [[ "${lines[1]}" =~ "1     admin       admin@example.local    true     2" ]]
}

@test "fetch users in multiple formats" {
    portusctl get users
    [ $status -eq 0 ]
    [[ "${lines[1]}" =~ "1     admin       admin@example.local    true     2" ]]

    portusctl get users -f json
    [ $status -eq 0 ]
    [[ "${lines[0]}" =~ '[{"id":1,"username":"admin","email":"admin@example.local"' ]]
}

@test "updating a user" {
    portusctl create user username=msabate email=example@test.lan password=12341234
    [ $status -eq 0 ]

    portusctl update user 2 display_name=Miquel email=something@test.lan
    [ $status -eq 0 ]
    [[ "${lines[1]}" =~ "ID    Username    Email                 Admin    NamespaceID    DisplayName" ]]
    [[ "${lines[2]}" =~ "2     msabate     something@test.lan    false    3              Miquel" ]]

    # Now quiet
    portusctl -q update user 2 email=example@test.lan
    [ $status -eq 0 ]
    [[ "${lines[0]}" =~ "ID    Username    Email               Admin    NamespaceID    DisplayName" ]]
    [[ "${lines[1]}" =~ "2     msabate     example@test.lan    false    3              Miquel" ]]
}
