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

@test "create a namespace and fetch it" {
    portusctl create team name=newteam
    portusctl create namespace name=newnamespace team=newteam
    [ $status -eq 0 ]
    [[ "${lines[1]}" =~ "ID    Name            Description    TeamID    Visibility    Global" ]]
    [[ "${lines[2]}" =~ "3     newnamespace                   2         private       false" ]]

    portusctl get namespaces
    [[ "${lines[1]}" =~ "1     registry:5000    The global namespace for the registry Registry.    0         public        true" ]]
    [[ "${lines[2]}" =~ "2     admin            This personal namespace belongs to admin.          1         private       false" ]]
    [[ "${lines[3]}" =~ "3     newnamespace                                                        2         private       false" ]]
}

@test "lists repositories in namespace" {
    portusctl get n 2 repositories
    [ $status -eq 0 ]
    [[ "${lines[0]}" =~ "ID    Name    FullName      Namespace        Stars    TagsCount" ]]
    [[ "${lines[1]}" =~ "1     fake    admin/fake    admin (ID: 2)    0        1" ]]
}

@test "validate namespace" {
    portusctl validate namespace name=something
    [ $status -eq 0 ]
    [[ "${lines[0]}" =~ "Valid" ]]

    portusctl validate namespace name=admin
    [ $status -eq 1 ]
    [[ "${lines[0]}" =~ "name:" ]]
    [[ "${lines[1]}" =~ "- Has already been taken" ]]
}
