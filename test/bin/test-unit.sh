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

# This file is largely based on previous work by Aleksa Sarai <asarai@suse.de>

set -e

GO="${GO:-go}"

# Set up the root and coverage directories.
export ROOT="$(readlink -f "$(dirname "$(readlink -f "$BASH_SOURCE")")/../..")"
export COVERAGE_DIR=$(mktemp --tmpdir -d portusctl-coverage.XXXXXX)

# Run the tests and collate the results.
$GO test -v -cover -covermode=count -coverprofile="$(mktemp --tmpdir=$COVERAGE_DIR cov.XXXXX)" -coverpkg=. . 2>/dev/null
chmod +x $ROOT/test/bin/collate.awk
$ROOT/test/bin/collate.awk $COVERAGE_DIR/* $COVERAGE | sponge $COVERAGE

# Clean up the coverage directory.
rm -rf "$COVERAGE_DIR"
