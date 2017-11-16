# portusctl [![Build Status](https://travis-ci.org/openSUSE/portusctl.svg?branch=master)](https://travis-ci.org/openSUSE/portusctl)

Under construction. To do:

- Polish exec command & integration tests for it
- Refactoring
- Coverage
- Man pages
- RPM spec

Missing resources on both sides:

- Webhooks
- Search and explore
- Ping and health
- Activities

Final considerations:

- What about typeahead kind of queries ? Should they be put in some special
  explore command ?

## Development

You could build this project as any other Go binary with `go build`, but this is
not recommended. Instead, use the default make target:

```
$ make
# or the equivalent `make portusctl`
```

With this command, `portusctl` will be built with the desired build flags and
setting the proper version for it. Note that the build can be further customized
with the `BUILD_FLAGS` variable. So, you could pass extra arguments like so:

```
$ make BUILD_FLAGS="-v"
```

When doing this you should be careful to not conflict with a default build flag.

### Unit testing

Unit testing is performed through the `test-unit` make target:

```
$ make test-unit
```

### Integration testing

We use [bats](https://github.com/sstephenson/bats.git) for integration testing
and `docker-compose`. There is a specific target on the `Makefile` called
`test-integration`, which will run the test integration suite:

```
$ make test-integration
```

This will setup a Portus instance running in the background as a Docker
container (through `docker-compose`), and then it will run the tests targeting
this Portus instance.

There are some nice flags when running integration tests:

- `SKIP_ENV_TESTS`: when set to true, it will skip the process of bringing your
  Portus instance up. You should use this flag when you want to re-use a
  previous environment for your tests (bear in mind that between each test case
  the database will be cleaned up). It is set to false by default.
- `TEARDOWN_TESTS`: when set to true, it will destroy the Portus instance at the
  end. It is set to true by default.
- `TESTS`: you can define a list of tests to be run. This way you can specify
  that you want to run only a specific subset of tests instead of the whole suite.

Taken that into account, you should perform the following when running tests for
the first time:

```
$ make test-integration TEARDOWN_TESTS=
```

This way you will have the Portus instance available and it won't be destroyed
at the end. Then, for successive runs you can perform:

```
$ make test-integration SKIP_ENV_TESTS=1 TEARDOWN_TESTS=
```

If only you care about a specific test (e.g. the `test/users.bats` file):

```
$ make test-integration TESTS=users SKIP_ENV_TESTS=1 TEARDOWN_TESTS=
```

### Validation

Besides running tests, we also perform some validation tests on the code. These
tests can be run like this:

```
$ make validate
```

### Other make targets

If you want to perform both unit and integration testing, then you can simply
call the `test` target like this:

```
$ make test
```

If you want to run all tests (including validations), you can run the same
target the CI environment will use:

```
$ make ci
```

## License

This project is based on work we did for the
[Portus](https://github.com/SUSE/Portus) project. However, all the code has been
re-written from scratch, so the entire project is subject to the GPLv3 license:

```
Copyright (C) 2017 Miquel Sabaté Solà <msabate@suse.com>

This program is free software: you can redistribute it and/or modify
it under the terms of the GNU General Public License as published by
the Free Software Foundation, either version 3 of the License, or
(at your option) any later version.

This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU General Public License for more details.

You should have received a copy of the GNU General Public License
along with this program.  If not, see <http://www.gnu.org/licenses/>.
```
