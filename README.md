# Buildsh

Buildsh is docker powered shell that make it easy to run isolated environment for building, testing and deploying softwares.

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

WIP...

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
