PORTUSCTL 1 "portusctl User manuals" "SUSE LLC." "AUGUST 2018"
================================================================

# NAME
portusctl explain \- Fetch the documentation of the available resources

# SYNOPSIS

**portusctl explain** \<resource\>

# DESCRIPTION
The **explain** command fetches the documentation of the available resources and
prints a helpful output. If no resource is given, then it simply lists the
available resources. Otherwise, if a resource was given, then it lists all the
available options around said resource.

# EXAMPLES

  $ portusctl explain
  ...(output)...
  You must specify one of the following types of resource:

    * application_tokens (aka 'at' or 'application_token')
    * namespaces (aka 'n' or 'namespace')
    * repositories (aka 'r' or 'repository')
    * registries (aka 're' or 'registry')
    * tags (aka 'tag')
    * teams (aka 't' or 'team')
    * users (aka 'u' or 'user')

  See the man pages for help and examples.

  $ portusctl explain namespace
  ...(output)...
  Resource: Namespaces (reference it with: 'n' or 'namespace' or 'namespaces')
  Supported commands: create, get and validate
  Required parameters:
    * In the 'create' and the 'update' commands: 'name' and 'team' (optional: 'description')
    * In the 'validate' command: 'name'

  Some commands will also accept the following 'subresources': 'repositories'
  For example, one might perform:
        $ portusctl get namespaces <id> repositories

  Refer to the man pages for usage examples

# BUGS

Please report bugs if you find them at https://github.com/openSUSE/portusctl/issues.

# AUTHOR

Miquel Sabaté Solà \<msabate@suse.com\>

# SEE ALSO

portusctl(1)

# HISTORY
August 2018, created by Miquel Sabaté Solà \<msabate@suse.com\>
