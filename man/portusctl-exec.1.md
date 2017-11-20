PORTUSCTL 1 "portusctl User manuals" "SUSE LLC." "NOVEMBER 2017"
================================================================

# NAME
portusctl exec \- Execute an arbitrary command on the environment of your Portus instance

# SYNOPSIS

**portusctl exec** \<command\> [arguments...]

# DESCRIPTION
The **exec** command allows users to execute an arbitrary command on the same
environment of your Portus instance.

# OPTIONS
**--local**, **-l**
  You can use this flag to set the location of the current host of your Portus
  instance. Only use this flag for non-containerized scenarios, since portusctl
  will default to the proper location in this case.

**--vendor**, **-v**
  With this flag you will instruct portusctl to use a local vendor directory as
  the gem environment.

# EXAMPLES
The **exec** command might come in handy when you want to execute some
administrative command in the same context as Portus. For example, if you want
to take a closer look at the current state of Portus, you can perform the
following:

```
$ portusctl exec rails c
```

This will put you inside of a Ruby on Rails console with Portus' code loaded in
it. This way, you will be able to perform deeper inspections like:

```
> puts Team.find_by(name: "myteam").namespaces
```

However, if you are not that experienced with Ruby on Rails and you want to
check the database directly as Portus sees it, you can perform:

```
$ portusctl exec rails db
```

With the command above, you will be able to access a MariaDB prompt that is
connected to the database that Portus is using.

# BUGS

Please report bugs if you find them at https://github.com/openSUSE/portusctl/issues.

# AUTHOR

Miquel Sabaté Solà \<msabate@suse.com\>

# SEE ALSO

portusctl(1)

# HISTORY
November 2017, imported from SUSE/Portus by Miquel Sabaté Solà \<msabate@suse.com\>

August 2016, created by Miquel Sabaté Solà \<msabate@suse.com\>
