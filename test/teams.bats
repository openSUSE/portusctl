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

@test "create a team and fetch it" {
    portusctl create team name=newteam
    [ $status -eq 0 ]
    [[ "${lines[1]}" =~ "ID    Name       Hidden" ]]
    [[ "${lines[2]}" =~ "2     newteam    false" ]]

    portusctl get teams
    [ $status -eq 0 ]
    [[ "${lines[0]}" =~ "ID    Name       Hidden" ]]
    [[ "${lines[1]}" =~ "2     newteam    false" ]]

    portusctl get team 2
    [ $status -eq 0 ]
    [[ "${lines[0]}" =~ "ID    Name       Hidden" ]]
    [[ "${lines[1]}" =~ "2     newteam    false" ]]
}

@test "create a namespace inside of a new team and list it" {
    portusctl create team name=newteam
    [ $status -eq 0 ]

    portusctl create namespace name=newnamespace team=newteam
    [ $status -eq 0 ]
    [[ "${lines[1]}" =~ "ID    Name            Description    TeamID    Visibility    Global" ]]
    [[ "${lines[2]}" =~ "3     newnamespace                   2         private       false" ]]

    portusctl get teams 2 namespaces
    [ $status -eq 0 ]
    [[ "${lines[0]}" =~ "ID    Name            Description    TeamID    Visibility    Global" ]]
    [[ "${lines[1]}" =~ "3     newnamespace                   2         private       false" ]]
}

# TODO: shortcuts
