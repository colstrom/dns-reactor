# dns-reactor

Monitors one or more hostnames for IP changes, and (optionally) execute a command.

Usage
-----

```dns-reactor -interval 30 -execute "service haproxy reload" github.com```

This would check the IP of github.com every 30 seconds, and execute ```service haproxy reload``` if it has changed since last check.

License
-------
[MIT](https://tldrlegal.com/license/mit-license)

Contributors
------------
* [Chris Olstrom](https://colstrom.github.io/) | [e-mail](mailto:chris@olstrom.com) | [Twitter](https://twitter.com/ChrisOlstrom)
