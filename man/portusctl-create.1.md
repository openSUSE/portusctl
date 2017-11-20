PORTUSCTL 1 "portusctl User manuals" "SUSE LLC." "NOVEMBER 2017"
================================================================

# NAME
portusctl create \- Creates the given resource

# SYNOPSIS

**portusctl create** \<resource\> [arguments...]

# DESCRIPTION
The **create** command issues a POST request to the RESP API for the given
resource. This way you can create a resource by simply passing its mandatory
fields as arguments. It will output the newly created resource on success.

# EXAMPLES

  $ portusctl create user username=user email="user@portus.test" password=12341234
  ...(output)...
  ID    Username    Email               Admin    NamespaceID    DisplayName
  3     user        user@portus.test    false    5

  $ portusctl create at id=2 application=another
  ...(output)...
  PlainToken
  SsC4PAx-b_9RsgbJ3si9

# BUGS

Please report bugs if you find them at https://github.com/openSUSE/portusctl/issues.

# AUTHOR

Miquel Sabaté Solà \<msabate@suse.com\>

# SEE ALSO

portusctl(1)

# HISTORY
November 2017, imported from SUSE/Portus by Miquel Sabaté Solà \<msabate@suse.com\>

August 2016, created by Miquel Sabaté Solà \<msabate@suse.com\>
