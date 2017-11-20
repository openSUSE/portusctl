PORTUSCTL 1 "portusctl User manuals" "SUSE LLC." "NOVEMBER 2017"
================================================================

# NAME
portusctl delete \- Deletes the given resource

# SYNOPSIS

**portusctl delete** \<resource\> [arguments...]

# DESCRIPTION
The **delete** command issues a DELETE request to the RESP API for the given
resource. This way you can delete a resource by simply passing the ID.

# EXAMPLES

  $ portusctl delete user 3

# BUGS

Please report bugs if you find them at https://github.com/openSUSE/portusctl/issues.

# AUTHOR

Miquel Sabaté Solà \<msabate@suse.com\>

# SEE ALSO

portusctl(1)

# HISTORY
November 2017, imported from SUSE/Portus by Miquel Sabaté Solà \<msabate@suse.com\>

August 2016, created by Miquel Sabaté Solà \<msabate@suse.com\>
