### Description

A general description of the problem.

### Steps to reproduce

1. First I did this...
2. Then that...
3. And this happened!

- **Expected behavior**: I expected this to happen!
- **Actual behavior**: But this happened...

### Deployment information

If relevant, provide the following information from the Portus deployment:

**Deployment method**: how have you deployed Portus? Are you using one of the
[examples](https://github.com/SUSE/Portus/tree/master/examples) as a base? If
possible, could you paste your configuration? (don't forget to strip passwords
or other sensitive data!)

**Configuration of Portus**:

You can get this information like this:

- In bare metal execute: `bundle exec rake portus:info`.
- In a container:
  - Using the development `docker-compose.yml` file: `docker exec -it <container-id> bundle exec rake portus:info`.
  - Using the [production image](https://hub.docker.com/r/opensuse/portus/): `docker exec -it <container-id> portusctl exec rake portus:info`.

```yml
CONFIG HERE
```

**Portus version**: with commit if possible. You can get this info too with the
above commands.
**portusctl version**: with commit if possible. You can get this info too with the
above commands.
