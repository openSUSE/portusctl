#!/usr/bin/env bash
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

set -e

start=`date +%s`

ROOT_DIR="$( cd "$( dirname "$0" )/../.." && pwd )"
CNAME="portus_portus_1"

export COVERAGE_DIR=$(mktemp --tmpdir -d portusctl-coverage.XXXXXX)
cp "$ROOT_DIR/portusctl" "$ROOT_DIR/test/portus/"

export DOCKER_COVERAGE_DIR=/srv/Portus/tmp/coverage
LOCAL_COVERAGE_DIR="$ROOT_DIR/test/portus/tmp/coverage"
rm -rf "$LOCAL_COVERAGE_DIR/*"

# Setup the environment
if [[ ! "$SKIP_ENV_TESTS" ]]; then
    pushd "$ROOT_DIR/test/portus"
    docker-compose kill
    docker-compose rm -f
    docker-compose up -d
    popd

    RETRY=1
    while [ $RETRY -ne 0 ]; do
        msg=$(SKIP_MIGRATION=1 docker exec $CNAME portusctl exec rails r /srv/Portus/bin/check_db.rb)
        case $(echo "$msg" | grep DB) in
            "DB_READY")
                echo "Database ready"
                break
                ;;
            *)
                echo "Database is not ready yet:"
                echo $msg
                ;;
        esac

        sleep 5
    done

    echo "You may want to set the 'SKIP_ENV_TESTS' env. variable for successive runs..."

    # Travis oddities...
    if [ ! -z "$CI" ]; then
        sleep 10
    fi
fi


# Run tests.
tests=()
if [[ -z "$TESTS" ]]; then
	tests=($ROOT_DIR/test/*.bats)
else
	for f in $TESTS; do
		tests+=("$ROOT_DIR/test/$f.bats")
	done
fi
bats -t ${tests[*]}
status=$?

# Tear down
if [[ "$TEARDOWN_TESTS" ]]; then
    pushd "$ROOT_DIR/test/portus"
    docker-compose kill
    docker-compose rm -f
    popd
fi

# Output coverage and clean.
$ROOT_DIR/test/bin/collate.awk $COVERAGE_DIR/* $LOCAL_COVERAGE_DIR/* $COVERAGE | sponge $COVERAGE
#rm -rf "$COVERAGE_DIR"

end=`date +%s`
dt=$(echo "$end - $start" | bc)
dm=$(echo "$dt/60" | bc)
ds=$(echo "$dt-60*$dm" | bc)
echo "Running integration tests took ${dm}m${ds}s"

exit $status
