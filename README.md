# Buildsh

Buildsh is docker powered shell that makes it easy to run isolated environment for building, testing and deploying softwares.
Internally, buildsh is a wrapper of `docker run` command that is implemented in bash script.

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

Buildsh boots a docker container using the default image `kohkimakimoto/buildsh:latest`, and starts bash process with interactive mode.
Your current working direcotory is automatically mounted to `/build` directory in the container, and several programming language runtimes (Go, Ruby, PHP, etc...) already be installed in the container. 
So you can run your project's tests by the following commands.

```
# php
$ php phpunit

# go
$ go test ./...
```

When you exit the shell, the container will be removed automatically.

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
envvar            "FOO2=bar2"
home_in_container "/build/src/github.com/kohkimakimoto/buildsh"
```

WIP...

## Supported Docker Images

You can use the following docker images with buildsh.

* `kohkimakimoto/buildsh:latest`: CentOS7 with some runtimes (*default) ([Dockerfile](build-images/standard/Dockerfile))
* `kohkimakimoto/buildsh:centos7-minimal`: CentOS7 minimal ([Dockerfile](build-images/centos7-minimal/Dockerfile))

Buildsh uses a docker image that is customized for some rules.

* Have `/build` directory.
* Default working directory is `/build`.
* Run a process by the user that is specified by `BUILDSH_USER` environment variables.
* If you doesn't specify the command. A container starts interactive shell.

For more detail, see [build-images/centos7-minimal/Dockerfile](build-images/centos7-minimal/Dockerfile) and [build-images/centos7-minimal/entrypoint.sh](build-images/centos7-minimal/entrypoint.sh)

If you want to use your custom docker image with buildsh, you should make a image like the [kohkimakimoto/buildsh](https://hub.docker.com/r/kohkimakimoto/buildsh/).

## Author

Kohki Makimoto <kohki.makimoto@gmail.com>

## License

The MIT License (MIT)
