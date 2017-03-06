# Buildsh

Buildsh is docker powered shell that make it easy to run isolated environment for building, testing and deploying softwares.

In implementation, Buildsh is a wrapper of `docker run` command.

## Requirements

* Docker
* Bash

## Installation

Clone git repository and set a $PATH.

```
$ git clone https://github.com/kohkimakimoto/buildsh.git ~/.buildsh
$ echo 'export PATH="$HOME/.buildsh/bin:$PATH"' >> ~/.bash_profile
```

After get a `buildsh`, run `buildsh -h` to check working.

```
$ buildsh -h
Usage: buildsh [<options...>] -- [<commands...>]
...
```

## Usage

Try to run `buildsh` without any options.

```
$ buildsh
```

Buildsh boots a docker container using the default image [kohkimakimoto/buildsh:latest](https://hub.docker.com/r/kohkimakimoto/buildsh/), and starts bash with interactive mode.
This container is automatically mounted current working direcotory to `/build` directory,
And several language runtime (Go, Ruby, PHP, etc...) already be installed. 
So you can run your project's tests like the following.

```
# php
$ php phpunit

# go
$ go test ./...
```

As you exit the shell, the container will be removed automatically.

```
$ exit
```

Buildsh can be used with arguments.

```
$ buildsh php phpunit
```

## Configuration

Buildsh loads configuration from `.buildshrc` in your current working directory. 
Here is an example:

```sh
use_cache
docker_image      "kohkimakimoto/buildsh:latest"
docker_option     "--net=host"
docker_option     "-v=/var/run/docker.sock:/var/run/docker.sock"
envvar            "FOO=bar"
home_in_container "/build/src/github.com/kohkimakimoto/buildsh"
```

WIP...

## Supported Docker Images

* [kohkimakimoto/buildsh:latest](https://hub.docker.com/r/kohkimakimoto/buildsh/) *default

## Use Custom Image

WIP...

## Author

Kohki Makimoto <kohki.makimoto@gmail.com>

## License

The MIT License (MIT)
