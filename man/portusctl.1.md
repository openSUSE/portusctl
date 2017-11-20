PORTUSCTL 1 "portusctl User manuals" "SUSE LLC." "NOVEMBER 2017"
================================================================

# NAME
portusctl \- A client for your Portus instance

# SYNOPSIS

**portusctl** command [command options] [arguments...]

# DESCRIPTION
**portusctl** is a tool that helps administrators perform operation on the
installed Portus instance. For this, administrators will be able to execute
various commands, **configure** the Portus instance, fetch the **logs** produced
by Portus, etc.

# OPTIONS
**portusctl** has some global options that can be set or toggled.

**--server**, **-s**
  The location where the Portus instance serves requests.

**--token**, **-t**
  The authentication token of the user for the Portus REST API.

**--user**, **-u**
  The user of the Portus REST API.

**--quiet**, **-q**
  Prevent portusctl from outputting general information

# COMMANDS
There are two types of commands: those that act upon a resource, and those
others that don't.

The commands that act upon a resource are as follows: **create**, **delete**,
**get**, **update** and **validate**. These commands translate a syntax like
"\<command\> \<resource\> [args...]" into REST API calls, and then present the
given results back to the user. Therefore, you can think of these commands as a
more useful way to target the REST API.

The **exec** command is the only command that doesn't take a resource as its
first argument, and it transforms its arguments into a command that will be run
on the Portus' context. As you can guess, this command assumes that you have
Portus running locally.

# EXAMPLES
Refer to the man page of each command, but here it is a quick example on how to
use **portusctl**:

  $ portusctl get users
  ...(output)...
  ID    Username    Email                Admin    NamespaceID    DisplayName
  1     portus      portus@portus.com    true     2
  2     mssola      admin@portus.test    true     3

  $ portusctl create user username=user email="user@portus.test" password=12341234
  ...(output)...
  ID    Username    Email               Admin    NamespaceID    DisplayName
  3     user        user@portus.test    false    5

  $ portusctl update user 3 display_name=User
  ...(output)...
  ID    Username    Email               Admin    NamespaceID    DisplayName
  3     user        user@portus.test    false    5              User

  $ portusctl delete user 3
  ...(output)...
  Deleted 'user' successfully!

Besides commands dealing with the API itself, you can run arbitrary commands on
the same context as a local Portus instance:

  $ portusctl exec cat .ruby-version
  2.4.2

# BUGS

Please report bugs if you find them at https://github.com/openSUSE/portusctl/issues.

# AUTHOR

Miquel Sabaté Solà \<msabate@suse.com\>

# SEE ALSO

portusctl-exec(1), portusctl-create(1), portusctl-delete(1), portusctl-get(1),
portusctl-update(1) and portusctl-validate(1)

# HISTORY
November 2017, imported from SUSE/Portus by Miquel Sabaté Solà \<msabate@suse.com\>

August 2016, created by Miquel Sabaté Solà \<msabate@suse.com\>
