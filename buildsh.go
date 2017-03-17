package main

import (
	"flag"
	"fmt"
	"github.com/Songmu/wrapcommander"
	"github.com/pkg/errors"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
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

	flag.StringVar(&optConfig, "config", "", "")
	flag.StringVar(&optConfig, "c", "", "")

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

	config, err := NewConfig()
	if err != nil {
		panic(err)
	}

	if optConfig != "" {
		config.ConfigFile = optConfig
	}

	// override config by the config file.
	if config.ConfigFile != "" {
		p, err := filepath.Abs(config.ConfigFile)
		if err != nil {
			panic(err)
		}

		b, err := ioutil.ReadFile(p)
		if err != nil {
			panic(err)
		}

		if err := yaml.Unmarshal(b, config); err != nil {
			panic(err)
		}
	} else {
		wd, err := os.Getwd()
		if err != nil {
			panic(errors.Wrap(err, "failed to get working directory."))
		}

		p := filepath.Join(wd, ".buildsh.yml")
		if _, err := os.Stat(p); err == nil {
			b, err := ioutil.ReadFile(p)
			if err != nil {
				panic(err)
			}

			if err := yaml.Unmarshal(b, config); err != nil {
				panic(err)
			}
		}
	}

	var dockerOptions = config.DockerOptions
	if flag.NArg() == 0 {
		dockerOptions = dockerOptions + " -i -t"
	}

	entrypoint, cmd := makeEntryPointAndCmd(flag.Args(), config)

	// construct docker run command.
	cmdline := "docker run -w " + config.ContainerWorkdir +
		" -v " + config.Home + ":" + config.ContainerHome +
		" " + dockerOptions +
		" " + config.AdditionalDockerOptions +
		` -e "BUILDSH=1"` +
		` -e "BUILDSH_USER=$(id -u):$(id -g)"` +
		` --entrypoint="` + entrypoint + `"` +
		" " + config.DockerImage +
		" " + cmd

	if optDebug {
		fmt.Println("[debug] " + cmdline)
	}

	if err := spawn(cmdline); err != nil {
		status = wrapcommander.ResolveExitCode(err)
	}

	return status
}

func makeEntryPointAndCmd(args []string, c *Config) (string, string) {
	var entrypoint = "/bin/bash"
	var cmd string

	if len(args) == 0 {
		cmd = ""
	} else {
		cmd = "-c '" + strings.Join(args, " ") + "'"
	}

	return entrypoint, cmd
}

type Config struct {
	DockerImage             string `yaml:"docker_image"`
	DockerOptions           string `yaml:"docker_options"`
	AdditionalDockerOptions string `yaml:"additional_docker_options"`
	ConfigFile              string `yaml:"-"`
	Home                    string `yaml:"home"`
	ContainerHome           string `yaml:"container_home"`
	ContainerWorkdir        string `yaml:"container_workdir"`
	Cmd                     string `yaml:"cmd"`
	UseCache                bool   `yaml:"use_cache"`
}

func NewConfig() (*Config, error) {
	wd, err := os.Getwd()
	if err != nil {
		return nil, errors.Wrap(err, "failed to get working directory.")
	}

	c := &Config{
		DockerImage:             "kohkimakimoto/buildsh:latest",
		DockerOptions:           "--rm -e TZ=Asia/Tokyo",
		AdditionalDockerOptions: "",
		ConfigFile:              "",
		Home:                    wd,
		ContainerHome:           "/build",
		ContainerWorkdir:        "/build",
		Cmd:                     "",
		UseCache:                false,
	}

	// Override default config by the environment variables.
	if d := os.Getenv("BUILDSH_DOCKER_IMAGE"); d != "" {
		c.DockerImage = d
	}
	if d := os.Getenv("BUILDSH_DOCKER_OPTIONS"); d != "" {
		c.DockerOptions = d
	}
	if d := os.Getenv("BUILDSH_ADDITIONAL_DOCKER_OPTIONS"); d != "" {
		c.AdditionalDockerOptions = d
	}
	if d := os.Getenv("BUILDSH_CONFIG"); d != "" {
		c.ConfigFile = d
	}
	if d := os.Getenv("BUILDSH_HOME"); d != "" {
		c.Home = d
	}
	if d := os.Getenv("BUILDSH_CONTAINER_HOME"); d != "" {
		c.ContainerHome = d
	}
	if d := os.Getenv("BUILDSH_CONTAINER_WORKDIR"); d != "" {
		c.ContainerWorkdir = d
	}
	if d := os.Getenv("BUILDSH_CMD"); d != "" {
		c.Cmd = d
	}
	if d := os.Getenv("BUILDSH_USE_CACHE"); d != "" {
		c.UseCache = true
	}

	return c, nil
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
