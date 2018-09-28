PORTUSCTL 1 "portusctl User manuals" "SUSE LLC." "SEPTEMBER 2018"
=================================================================

# NAME
portusctl health \- Get health info from Portus

# SYNOPSIS

**portusctl bootstrap** [arguments...]

# DESCRIPTION
The **bootstrap** command allows you to create the first administrator of your
Portus instance. This will only work if there are no users in the system and if
the *first_user_admin* option from Portus is enabled. The arguments to be passed
are the same as for creating a new user: username, email and password are
mandatory; display_name is optional.

# EXAMPLES
This command is pretty straight-forward:

    $ portusctl bootstrap username=admin email=admin@portus.test password=12341234

On success, you will get the token to be used for the newly created user.

# BUGS

Please report bugs if you find them at https://github.com/openSUSE/portusctl/issues.

# AUTHOR

Miquel Sabaté Solà \<msabate@suse.com\>

# SEE ALSO

portusctl(1)

# HISTORY
September 2018, created by Miquel Sabaté Solà \<msabate@suse.com\>
