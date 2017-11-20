PORTUSCTL 1 "portusctl User manuals" "SUSE LLC." "NOVEMBER 2017"
================================================================

# NAME
portusctl get \- Fetches info for the given resource

# SYNOPSIS

**portusctl get** \<resource\> [arguments...]

# DESCRIPTION
The **get** command issues a GET request to the RESP API for the given
resource. You can fetch multiple rows or just a single one by specifying the ID
of it. Moreover, some resources have sub-resources (e.g. users). In these cases,
in order to fetch sub-resources, you have to first specify the ID of the
super-resource object.

# OPTIONS
**--format**, **-f**
  The format in which the output given by the server should be presented. The
  default is a tabular-based output. You can specify **json**, which will print
  the raw JSON response.

# EXAMPLES

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

  $ ./portusctl get users 2 at
  ...(output)...
  ID    Application
  1     asdas
# BUGS

Please report bugs if you find them at https://github.com/openSUSE/portusctl/issues.

# AUTHOR

Miquel Sabaté Solà \<msabate@suse.com\>

# SEE ALSO

portusctl(1)

# HISTORY
November 2017, imported from SUSE/Portus by Miquel Sabaté Solà \<msabate@suse.com\>

August 2016, created by Miquel Sabaté Solà \<msabate@suse.com\>
