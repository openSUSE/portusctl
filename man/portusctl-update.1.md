PORTUSCTL 1 "portusctl User manuals" "SUSE LLC." "NOVEMBER 2017"
================================================================

# NAME
portusctl update \- Updates the given resource

# SYNOPSIS

**portusctl update** \<resource\> [arguments...]

# DESCRIPTION
The **update** command issues a PUT request to the RESP API for the given
resource. This way you can update a resource by simply passing its mandatory
fields as arguments plus the ID of the resource to be updated.

# EXAMPLES

  $ portusctl update user 3 display_name=User
  ...(output)...
  ID    Username    Email               Admin    NamespaceID    DisplayName
  3     user        user@portus.test    false    5              User

# BUGS

Please report bugs if you find them at https://github.com/openSUSE/portusctl/issues.

# AUTHOR

Miquel Sabaté Solà \<msabate@suse.com\>

# SEE ALSO

portusctl(1)

# HISTORY
November 2017, imported from SUSE/Portus by Miquel Sabaté Solà \<msabate@suse.com\>

August 2016, created by Miquel Sabaté Solà \<msabate@suse.com\>
