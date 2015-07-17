# Docker Inject
Copy files/directories from hosts to running Docker containers.

## Description

Surely you sometimes desire to copy files
from your host machine to running Docker containers.
Indeed `docker cp` provide the way to copy files from containers to hosts,
but we have no command for the retrorse copying.

One of the ways to transport objects is to mount a directory
by the `-v` option of `docker run` and use the directory as a portal.
Meanwhile this method is, as you might notice,
completely helpless for running containers.

Docker Inject provide the way to copy objects from hosts to running containers.
With the command, you are able to inject not only a single file
but also directories recursively.

## Requirement

* [Docker 1.3+](https://www.docker.com/)
* [boot2docker](http://boot2docker.io/) (for Mac OS X and Windows)

## Installation

Download the appropriate zip package for your system.

* [Linux (64bit)](https://github.com/y-taka-23/docker-inject/releases/download/v0.1.0/docker-inject_0.1.0_linux_amd64.zip)
* [Mac OS X (64bit)](https://github.com/y-taka-23/docker-inject/releases/download/v0.1.0/docker-inject_0.1.0_darwin_amd64.zip)
* [Windows (64bit)](https://github.com/y-taka-23/docker-inject/releases/download/v0.1.0/docker-inject_0.1.0_windows_amd64.zip)

After downloading, unzip the package and copy the `docker-inject` binary
to somewhere on the `PATH` so that it can be executed.

Incidentally, `docker-inject` is written in Go,
thus you can build the latest version from the source by:

```
$ go get https://github.com/y-taka-23/docker-inject
```

## Usage

## License

[MIT](https://github.com/y-taka-23/docker-inject/blob/master/LICENSE)

## Author

[y-taka-23](https://github.com/y-taka-23)
