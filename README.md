# Buildsh

Buildsh is docker powered shell that makes it easy to run isolated environment for building, testing and deploying softwares.
Internally, buildsh is a wrapper of `docker run` command that is implemented in GO.

## Requirements

* Docker

## Installation

[Download latest version](https://github.com/kohkimakimoto/buildsh/releases/latest)

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

Buildsh loads configuration from `.buildsh.yml` in your current working directory. 

Example:

```yaml
use_cache: true
docker_image: kohkimakimoto/buildsh:latest
docker_options: --net=host -v=/var/run/docker.sock:/var/run/docker.sock
enviroment:
  FOO: bar
  FOO2: bar2
home_in_container: /build/src/github.com/kohkimakimoto/buildsh
```

Description:

* `use_cache`: If you set it, buildsh creates `.buildsh/cache` directory in a current directory. It also set the environment variable `BUILDSH_USE_CACHE=1` and `BUILDSH_CACHEDIR` which stores the path to the cache directory.

* `docekr_image`: Specifies a docker image to run. Default `kohkimakimoto/buildsh:latest`.

* `docker_options`: Options that are appended to the `docker run` that is executed by bashsh internally.

* `environment`: Specifies a environment variable in a container. 

* `home_in_container`: Changes mount point and current working directory in a container. Default `/build`.

## Using With Shebang

If you want to create a script file that are executed by buildsh, you can use a trick to interpret shebang with buildsh. See the following example code.

```sh
#!/bin/sh
[ -z "$BUILDSH" ] && exec buildsh "$0" "$@"

# your code is after here...
echo "I'm in a container!"
```

You can run it directly after adding an execution permission.

```sh
$ chmod 755 your_script.sh
$ ./your_script.sh
I'm in a container!
```




## Author

Kohki Makimoto <kohki.makimoto@gmail.com>

## License

The MIT License (MIT)
