PORTUSCTL 1 "portusctl User manuals" "SUSE LLC." "NOVEMBER 2017"
================================================================

# NAME
portusctl validate \- Validates the given arguments for a resource

# SYNOPSIS

**portusctl validate** \<resource\> [arguments...]

# DESCRIPTION
The **validate** command issues a GET request targeting the `/validate`
endpoints of the RESP API. This command is useful for validating certain data
from some resouces. Note that the exit status of this command will also tell you
whether the given validation passed or not.

# EXAMPLES

  $ portusctl validate namespace name=asd
  ...(output)...
  Valid

  $ portusctl validate namespace name=mssola
  ...(output)...
  name:
    - Has already been taken

  $ portusctl validate namespace name=/
  ...(output)...
  name:
    - Can only contain lower case alphanumeric characters, with optional underscores and dashes in the middle.

# BUGS

Please report bugs if you find them at https://github.com/openSUSE/portusctl/issues.

# AUTHOR

Miquel Sabaté Solà \<msabate@suse.com\>

# SEE ALSO

portusctl(1)

# HISTORY
November 2017, imported from SUSE/Portus by Miquel Sabaté Solà \<msabate@suse.com\>

August 2016, created by Miquel Sabaté Solà \<msabate@suse.com\>
