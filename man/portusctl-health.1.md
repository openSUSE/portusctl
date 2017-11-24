PORTUSCTL 1 "portusctl User manuals" "SUSE LLC." "NOVEMBER 2017"
================================================================

# NAME
portusctl health \- Get health info from Portus

# SYNOPSIS

**portusctl health**

# DESCRIPTION
Get health information from Portus. This information can vary depending on
whether some features are enabled or not. The minimal information being shown is
the status of the registry and the database. Other than that, for example, you
will get the status of security scanners if you have this feature enabled.

# EXAMPLES
This command is pretty straight-forward:

```
$ portusctl health
...(output)...
DATABASE                  REGISTRY
Database is up-to-date    Registry is reachable
```

# BUGS

Please report bugs if you find them at https://github.com/openSUSE/portusctl/issues.

# AUTHOR

Miquel Sabaté Solà \<msabate@suse.com\>

# SEE ALSO

portusctl(1)

# HISTORY
November 2017, imported from SUSE/Portus by Miquel Sabaté Solà \<msabate@suse.com\>

August 2016, created by Miquel Sabaté Solà \<msabate@suse.com\>
