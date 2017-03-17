package main

import (
	"flag"
	"fmt"
	"os"
	"os/exec"
	"runtime"
)

var (
	Version    = "0.1.0"
	CommitHash = "unknown"
)

func main() {
	os.Exit(realMain())
}

func realMain() (status int) {
	defer func() {
		if err := recover(); err != nil {
			printError(err)
			status = 1
		}
	}()

	// parse flags...
	var optVersion, optDebug, optClean bool
	var optConfig string
	var optEnv stringSlice

	flag.BoolVar(&optVersion, "v", false, "")
	flag.BoolVar(&optVersion, "version", false, "")
	flag.BoolVar(&optDebug, "d", false, "")
	flag.BoolVar(&optDebug, "debug", false, "")
	flag.BoolVar(&optClean, "clean", false, "")

	flag.StringVar(&optConfig, "config", ".buildsh.yml", "")
	flag.StringVar(&optConfig, "c", ".buildsh.yml", "")

	flag.Var(&optEnv, "e", "")
	flag.Var(&optEnv, "env", "")

	flag.Usage = func() {
		fmt.Println(`Usage: buildsh [<options...>] -- [<commands...>]

Buildsh is docker powered shell that make it easy to run isolated
environment for building, testing and deploying softwares.

The MIT License (MIT)
Kohki Makimoto <kohki.makimoto@gmail.com>

version ` + Version + ` (` + CommitHash + `)

Options:
    -e, --env <KEY=VALUE>      Set custom environment variables.
    -d, --debug                Use debug mode.
    -c, --config <FILE>        Load configuration from the FILE instead of .buildsh.yml
    --clean                    Remove cache.
    -h, --help                 Show help.
    -v, --version              Print the version
`)
	}
	flag.Parse()

	if optVersion {
		fmt.Println("buildsh version " + Version + " (" + CommitHash + ")")
		return 0
	}

	return status
}

type stringSlice []string

// Now, for our new type, implement the two methods of
// the flag.Value interface...
// The first method is String() string
func (s *stringSlice) String() string {
	return fmt.Sprintf("%s", *s)
}

// The second method is Set(value string) error
func (s *stringSlice) Set(value string) error {
	*s = append(*s, value)
	return nil
}

func printError(err interface{}) {
	fmt.Fprintf(os.Stderr, "buildsh error: %v\n", err)
}

func spawn(command string) error {
	var shell, flag string
	if runtime.GOOS == "windows" {
		shell = "cmd"
		flag = "/C"
	} else {
		shell = "bash"
		flag = "-c"
	}
	cmd := exec.Command(shell, flag, command)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Start()
	if err != nil {
		return err
	}
	return cmd.Wait()
}
