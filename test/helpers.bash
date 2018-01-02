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

# Function taken from openSUSE/umoci. See:
# https://github.com/openSUSE/umoci/blob/57c73c27fe3c13d80e1fb7f82c9a046a2bc2b6f1/test/helpers.bash#L116-L125
function sane_run() {
	local cmd="$1"
	shift

	run "$cmd" "$@"

	# Some debug information to make life easier.
	echo "$(basename "$cmd") $@ (status=$status)" >&2
	echo "$output" >&2
}

# Source the config file so we have available some relevant environment
# variables.
function __source_environment() {
    . "$BATS_TEST_DIRNAME/portus/tmp/config.sh"
}

# Setup the database for each test case.
function __setup_db() {
    ROOT_DIR="$( cd "$( dirname "$BATS_TEST_DIRNAME" )" && pwd )"
    PORTUSCTL="$ROOT_DIR/portusctl"
    docker exec \
           -e PORTUSCTL=$PORTUSCTL \
           portus_portus_1 \
           bundle exec rails r /srv/Portus/bin/runner.rb
}

# Wrapper for the main command. Use this instead of running the binary by
# hand. Coverage part taken from openSUSE/umoci.
function portusctl() {
	local cover=()
    cover+=("-test.coverprofile=$(mktemp -p "$COVERAGE_DIR" portusctl.cov.XXXXXX)")
    cover+=("__DEVEL--cover-tests")

    sane_run $PORTUSCTL "${cover[@]}" $@
}

# Wrapper that allows us to run commands inside of the Portus container.
function docker_run() {
    # Doing the mapping between our local `tmp` file, and the one from the
    # docker container.
    ROOT_DIR="$( cd "$( dirname "$BATS_TEST_DIRNAME" )" && pwd )"
    tmp=$(mktemp -p "$ROOT_DIR/test/portus/tmp/coverage" portusctl.cov.XXXXXX)
    test=$(basename $tmp)

	local cover=()
    cover+=("-test.coverprofile=$DOCKER_COVERAGE_DIR/$test")
    cover+=("__DEVEL--cover-tests")

    # Discard the initial `portusctl` argument, so it's easier to call `sane_run`.
    shift

    sane_run docker exec -e PORTUSCTL=$PORTUSCTL portus_portus_1 portusctl "${cover[@]}" $@
}
