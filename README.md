# buildsh

A docker image and shell script for building testing and deploying softwares.

WIP...

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

Run an arbitrary command in the isolated build container.
If you run a command without any options,
buildsh boots the container with interactive shell (default bash).

Options:
    -e, --env <KEY=VALUE>   Set custom environment variables
    -h, --help              Show help

Examples:
    buildsh
    buildsh -- ls -la

Configuration:
    buildsh loads .buildshrc file if it is existed in your current directory.
    You can set custom environment to change buildsh behavior.
```

## Configuration

WIP...

## Create Custom Image

WIP...

## Author

Kohki Makimoto <kohki.makimoto@gmail.com>

## License

The MIT License (MIT)
