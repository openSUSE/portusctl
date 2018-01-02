#!/usr/bin/awk -f
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

# Make sure that the coverage is above a certain percentage and exit properly.

BEGIN {
    status = 0;
    coverage = (length(coverage) == 0) ? 100.0 : strtonum(coverage);
}

{
    # Ignore records which sum up everything.
    if ($1 ~ /^total/) {
        next;
    }

    # Strip away the '%' of the third column and convert it into a proper
    # numeric type.
    num = $3;
    gsub("\\%", "", num);
    num = strtonum(num);

    # Print the original record if the current row is below the desired coverage.
    if (num < coverage) {
        print;
        status = 1;
    }
}

END {
    # TODO: this will be removed in the future.
    if (length(allow_failure) == 0) {
        exit status;
    }
    exit 0;
}
