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

@test "create application token" {
    portusctl create application_token id=1 application=something
    [ $status -eq 0 ]
    [[ "${lines[1]}" =~ "PlainToken" ]]

    portusctl get users 1 application_tokens
    [ $status -eq 0 ]
    [[ "${lines[0]}" =~ "ID    Application" ]]
    [[ "${lines[1]}" =~ "1     app" ]]
    [[ "${lines[2]}" =~ "2     something" ]]
}

@test "creating too many application tokens is not allowed" {
    portusctl create application_token id=1 application=something
    [ $status -eq 0 ]
    portusctl create at id=1 application=something2
    [ $status -eq 0 ]
    portusctl create application_tokens id=1 application=something3
    [ $status -eq 0 ]
    portusctl create application_tokens id=1 application=something4
    [ $status -eq 0 ]

    portusctl create application_tokens id=1 application=something5
    [ $status -eq 1 ]
    [[ "${lines[0]}" =~ "- Base: Users cannot have more than 5 application tokens" ]]
}

@test "an application token cannot be updated" {
    portusctl create application_token id=1 application=something
    [ $status -eq 0 ]

    portusctl update application_tokens 2 application=another
    [ $status -eq 1 ]
    [[ "${lines[0]}" =~ "Action not supported for resource 'application token'" ]]
}
