# portusctl [![Build Status](https://travis-ci.org/openSUSE/portusctl.svg?branch=master)](https://travis-ci.org/openSUSE/portusctl) [![Go Report Card](https://goreportcard.com/badge/github.com/openSUSE/portusctl)](https://goreportcard.com/report/github.com/openSUSE/portusctl)

**portusctl** is a client for your [Portus](https://github.com/SUSE/Portus)
instance. It allows you to access the REST API offered by Portus and present the
results in a friendly manner. For example:

```bash
$ portusctl get users
...(output)...
ID    Username    Email                Admin    NamespaceID    DisplayName
1     portus      portus@portus.com    true     2
2     mssola      admin@portus.test    true     3

$ portusctl get -f json users 2 | jq
...(output)...
{
    "id": 2,
    "username": "mssola",
    "email": "admin@portus.test",
    "current_sign_in_at": "2017-11-22T14:30:52.000Z",
    "last_sign_in_at": "2017-11-16T09:23:34.000Z",
    "created_at": "2017-11-14T14:22:01.000Z",
    "updated_at": "2017-11-22T14:30:52.000Z",
    "admin": true,
    "enabled": true,
    "locked_at": null,
    "namespace_id": 3,
    "display_name": ""
}

$ portusctl create user username=user email=user@portus.test password=12341234
...(output)...
ID    Username    Email               Admin    NamespaceID    DisplayName
3     user        user@portus.test    false    5

$ portusctl update user 3 display_name=User
...(output)...
ID    Username    Email               Admin    NamespaceID    DisplayName
3     user        user@portus.test    false    5              User
```

Moreover, it also allows you to execute arbitrary commands on the Portus'
context if your instance is running locally:

```bash
$ portusctl exec cat .ruby-version
2.4.2
```

**portusctl** does **not** implement all the API entrypoints available as of
Portus `v2.4`. You can follow the progress of this
[here](https://github.com/openSUSE/portusctl/issues/8). The **exec** command is
considered stable already, and that's why **portusctl** is included inside of
the Portus Docker image as of `v2.4`.

## Installation

You can install `portusctl` from source by cloning this repository and then
performing the following command:

```bash
$ make install
```

You can also install `portusctl` from
[obs://Virtualization:containers:Portus](https://build.opensuse.org/package/show/Virtualization:containers:Portus/portusctl). In
order to install it with zypper you need to perform the following commands:

```bash
% zypper ar -f https://download.opensuse.org/repositories/Virtualization:/containers:/Portus/openSUSE_Leap_42.3/ portus
% zypper install portusctl
```

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

### Man pages

Man pages have been written using Markdown, and they can be found in the `man`
directory. In order to generate man pages from these markdown files, you have to
run the following command:

```
$ make doc
```

### Code coverage

If you want to perform both unit and integration testing, then you can simply
call the `test` target like this:

```
$ make test
```

This target has one extra benefit: it will also check for code coverage. If code
coverage is below an expected threshold, then you will get a report about.

Note though that you need the `sponge` command installed on your system.

### What the CI will end up running

The CI will only run the `ci` target, which will in turn:

1. Run all validators.
2. Run unit & integration tests.
3. Perform checks on code coverage.

It is recommended that you perform `make ci` before submitting a pull request,
and check that it ran successfully.

## License

This project is based on work we did for the
[Portus](https://github.com/SUSE/Portus) project. However, all the code has been
re-written from scratch, so the entire project is subject to the GPLv3 license:

```
Copyright (C) 2017-2018 Miquel Sabaté Solà <msabate@suse.com>

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
