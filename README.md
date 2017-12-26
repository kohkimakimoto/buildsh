# Buildsh

Several script files make it easy to run a script in an isolated docker container for building, testing and deploying softwares. 

## Requirements

* Docker

## Usage

Try to run `buildsh-default` without any options.

```
$ buildsh-default
```

The `buildsh-default` boots a docker container using the docker image `kohkimakimoto/buildsh:latest`, and starts a bash process with interactive mode. Your current working direcotory is automatically mounted to `/build` directory in the container, and several programming language runtimes (Go, Ruby, PHP, etc...) already be installed in the container. 
So you can run your project's tests by the following commands.

```
# php
$ php phpunit

# go
$ go test ./...

# shell script
$ bash ./tests.sh
```

When you exit the shell, the container will be removed automatically.

```
$ exit
```

If you use `buildsh-default` in non interactive mode, you use it with '-c' option with commands executed in the container.

```
$ buildsh-default -c 'php phpunit'
```

## Author

Kohki Makimoto <kohki.makimoto@gmail.com>

## License

The MIT License (MIT)
