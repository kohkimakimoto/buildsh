# Buildsh

Buildsh is docker powered shell that make it easy to run isolated environment for building, testing and deploying softwares.

In implementation, buildsh is a wrapper of `docker run` command.

## Requirements

* Docker
* Bash

## Installation

Clone git repository and set a $PATH.

```
$ git clone https://github.com/kohkimakimoto/buildsh.git ~/.buildsh
$ echo 'export PATH="$HOME/.buildsh/bin:$PATH"' >> ~/.bash_profile
```

Run `buildsh -h` to check working.

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
Your current working direcotory is automatically mounted to `/build` directory in the container, and several language runtime (Go, Ruby, PHP, etc...) already be installed. 
So you can run your project's tests by the following commands.

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

If you use buildsh in non interactive mode, you use it with arguments that are the commands executed in the container.

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

Buildsh uses a docker image that is customized for some conventions.
See [build-images/standard/Dockerfile](build-images/standard/Dockerfile) and [build-images/standard/entrypoint.sh](build-images/standard/entrypoint.sh)

If you use your custom docker image, you should make a image like the [kohkimakimoto/buildsh:latest](https://hub.docker.com/r/kohkimakimoto/buildsh/).

WIP...

## Author

Kohki Makimoto <kohki.makimoto@gmail.com>

## License

The MIT License (MIT)
