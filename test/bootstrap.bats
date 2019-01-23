#!/usr/bin/env bats -t
# Copyright (C) 2018-2019 Miquel Sabaté Solà <msabate@suse.com>
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
    __setup_db nopopulate
    __source_environment
}

@test "bootstrap works" {
    portusctl bootstrap username=admin password=12341234 email=admin@example.com
    [ $status -eq 0 ]
    [[ "${lines[-3]}" =~ "You can now use the user 'admin' with the following token" ]]
}

@test "refuse to bootstrap when a user has already been created" {
    portusctl bootstrap username=admin password=12341234 email=admin@example.com
    [ $status -eq 0 ]

    portusctl bootstrap username=another password=12341234 email=another@example.com
    [ $status -eq 1 ]
    [[ "${lines[-1]}" =~ "You can only use this when there are no users on the system" ]]
}
