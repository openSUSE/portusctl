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

@test "fetch repositories" {
    portusctl get repositories
    [ $status -eq 0 ]
    [[ "${lines[0]}" =~ "ID    Name    FullName      Namespace        Stars    TagsCount" ]]
    [[ "${lines[1]}" =~ "1     fake    admin/fake    admin (ID: 2)    0        1" ]]

    portusctl get repository 1
    [ $status -eq 0 ]
    [[ "${lines[0]}" =~ "ID    Name    FullName      Namespace        Stars    TagsCount" ]]
    [[ "${lines[1]}" =~ "1     fake    admin/fake    admin (ID: 2)    0        1" ]]
}

@test "fetch tags from repository" {
    portusctl get r 1 tags
    [ $status -eq 0 ]
    [[ "${lines[0]}" =~ "ID    Name    Author           Digest    ImageID" ]]
    [[ "${lines[1]}" =~ "1     tag1    admin (ID: 1)    digest    imageid" ]]
}

@test "fetch tags from unknown repository" {
    portusctl get repositories 99 tags
    [ $status -eq 1 ]
    [[ "${lines[0]}" =~ "Resource not found" ]]
}
