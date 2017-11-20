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

# This Makefile has taken lots of ideas and code from openSUSE/umoci by Aleksa Sarai.

GO ?= go
GO_MD2MAN ?= go-md2man
CMD ?= portusctl
GO_SRC = $(shell find . -name \*.go)

# Version information.
VERSION := $(shell cat ./VERSION)
COMMIT_NO := $(shell git rev-parse HEAD 2> /dev/null || true)
COMMIT := $(if $(shell git status --porcelain --untracked-files=no),"${COMMIT_NO}-dirty","${COMMIT_NO}")

# Test integration
SKIP_ENV_TESTS ?=
TEARDOWN_TESTS ?= 1

# Build flags and settings.
BUILD_FLAGS ?=
DYN_BUILD_FLAGS := $(BUILD_FLAGS) -buildmode=pie -ldflags "-s -w -X main.gitCommit=${COMMIT} -X main.version=${VERSION}" -tags "$(BUILDTAGS)"

.DEFAULT: portusctl
portusctl: $(GO_SRC)
	@$(GO) build ${DYN_BUILD_FLAGS} -o $(CMD)

.PHONY: clean
clean:
	@rm -rf $(CMD)
	@rm -f ./man/*.1
	@rm -f ./man/*.out

#
# Unit & integration tests.
#

.PHONY: test-unit
test-unit:
	@go test -v ./...

.PHONY: test-integration
test-integration: portusctl
	@chmod +x ./test/bin/test-integration.sh
	SKIP_ENV_TESTS="${SKIP_ENV_TESTS}" TEARDOWN_TESTS="${TEARDOWN_TESTS}" ./test/bin/test-integration.sh

.PHONY: test
test: test-unit test-integration

#
# Validation tools.
#

EPOCH_COMMIT ?= b6061031a600
.PHONY: validate-git
validate-git:
	@which git-validation > /dev/null 2>/dev/null || (echo "ERROR: git-validation not found." && false)
ifdef TRAVIS_COMMIT_RANGE
	git-validation -q
else
	git-validation -q -range $(EPOCH_COMMIT)..HEAD
endif

.PHONY: validate-go
validate-go:
	@which gofmt >/dev/null 2>/dev/null || (echo "ERROR: gofmt not found." && false)
	test -z "$$(gofmt -s -l . | grep -vE '^vendor/' | tee /dev/stderr)"
	@which golint >/dev/null 2>/dev/null || (echo "ERROR: golint not found." && false)
	test -z "$$(golint . | grep -v vendor | tee /dev/stderr)"
	@go doc cmd/vet >/dev/null 2>/dev/null || (echo "ERROR: go vet not found." && false)
	test -z "$$(go vet . | grep -v vendor | tee /dev/stderr)"

.PHONY: validate
validate: validate-git validate-go

#
# Man pages
#

MANPAGES_MD := $(wildcard man/*.md)
MANPAGES    := $(MANPAGES_MD:%.md=%)

man/%.1: man/%.1.md
	@$(GO_MD2MAN) -in $< -out $@.out
	@go run man/sanitize.go $@.out &> $@
	@rm $@.out

.PHONY: doc
doc: $(MANPAGES)

#
# Travis-CI
#

.PHONY: ci
ci: portusctl doc validate test
