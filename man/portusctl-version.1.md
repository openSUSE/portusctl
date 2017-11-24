PORTUSCTL 1 "portusctl User manuals" "SUSE LLC." "NOVEMBER 2017"
================================================================

# NAME
portusctl version \- Print the client and server version information

# SYNOPSIS

**portusctl version**

# DESCRIPTION
This command differs from the plain **-v, --version** global flag in the fact
that it outputs a full description of the versions being targeted. In
particular, it prints:

- The version of the **server**, and the APIs supported by it.
- Some **Git** information from the server if possible.
- The version of the **client** (portusctl), and the APIs supported by it.

# EXAMPLES
This command is pretty straight-forward:

```
$ portusctl version
...(output)...
SERVER VERSION    SERVER BRANCH          CLIENT VERSION    SERVER API    CLIENT API
2.3.0-dev         master@b0b55d57e39e    0.1.0             v1            v1
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
