package main

import (
	"bytes"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"text/template"

	"github.com/Songmu/wrapcommander"
	"github.com/pkg/errors"
	"gopkg.in/yaml.v2"
)

var (
	Version    = "0.4.0"
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

Examples:
    buildsh
    buildsh -- ls -la

Configuration:
    Buildsh loads .buildsh.yml file if it is existed in your current directory.

Description:
    Run an arbitrary command in the isolated container.
    If you run a command without any options,
    buildsh boots the container with interactive shell (default bash).

See also:
    https://github.com/kohkimakimoto/buildsh
`)
		os.Exit(0)
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

	// override config by the config file.
	if optConfig != "" {
		p, err := filepath.Abs(optConfig)
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

	if config.HomeInContainer != "" {
		config.ContainerHome = config.HomeInContainer
		config.ContainerWorkdir = config.HomeInContainer
	}

	// appned env
	var envOptions = ""
	if len(config.Environment) > 0 {
		for k, v := range config.Environment {
			envOptions = envOptions + " -e " + k + "=" + shellEscape(v)
		}
	}

	if len(optEnv) > 0 {
		for _, v := range optEnv {
			envOptions = envOptions + " -e " + shellEscape(v)
		}
	}

	// cache config
	var envForCache = ""
	if config.UseCache {
		var dir string
		if !filepath.IsAbs(config.Cachedir) {
			dir = filepath.Join(config.Home, config.Cachedir)
		} else {
			dir = config.Cachedir
		}

		if err := os.MkdirAll(dir, 0777); err != nil {
			panic(err)
		}

		envForCache = "-e BUILDSH_USE_CACHE=1 -e BUILDSH_CACHEDIR=" + config.ContainerHome + "/.buildsh/cache"
	}

	entrypoint, cmd, err := makeEntryPointAndCmd(flag.Args(), config)
	if err != nil {
		panic(err)
	}

	// construct docker run command.
	cmdline := "docker run -w " + config.ContainerWorkdir +
		" -v " + config.Home + ":" + config.ContainerHome +
		" " + config.DockerOptions +
		" " + config.AdditionalDockerOptions +
		" " + envForCache +
		" " + envOptions +
		` -e BUILDSH=1` +
		` -e BUILDSH_USER=$(id -u):$(id -g)` +
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

func makeEntryPointAndCmd(args []string, c *Config) (string, string, error) {
	var entrypoint = "/bin/bash"
	var cmd string

	funcMap := template.FuncMap{
		"ShellEscape": shellEscape,
	}

	tmpl, err := template.New("T").Funcs(funcMap).Parse(realEntrypointTemplate)
	if err != nil {
		return "", "", err
	}

	var mainCommand string
	if len(args) > 0 {
		mainCommand = strings.Join(args, " ")
	} else {
		mainCommand = ""
	}

	dict := map[string]interface{}{
		"Cmd": mainCommand,
	}

	var b bytes.Buffer
	err = tmpl.Execute(&b, dict)
	if err != nil {
		return "", "", err
	}

	realEntrypoint := shellEscape(b.String())
	cmd = " -c " + realEntrypoint

	return entrypoint, cmd, nil
}

type Config struct {
	DockerImage             string            `yaml:"docker_image"`
	DockerOptions           string            `yaml:"docker_options"`
	AdditionalDockerOptions string            `yaml:"additional_docker_options"`
	Home                    string            `yaml:"home"`
	Environment             map[string]string `yaml:"environment"`
	ContainerHome           string            `yaml:"container_home"`
	ContainerWorkdir        string            `yaml:"container_workdir"`
	HomeInContainer         string            `yaml:"home_in_container"`
	UseCache                bool              `yaml:"use_cache"`
	Cachedir                string            `yaml:"cahcedir"`
}

func NewConfig() (*Config, error) {
	wd, err := os.Getwd()
	if err != nil {
		return nil, errors.Wrap(err, "failed to get working directory.")
	}

	c := &Config{
		DockerImage:             "kohkimakimoto/buildsh:latest",
		DockerOptions:           "-i -t --rm -e TZ=Asia/Tokyo",
		AdditionalDockerOptions: "",
		Home:             wd,
		Environment:      map[string]string{},
		ContainerHome:    "/build",
		ContainerWorkdir: "/build",
		HomeInContainer:  "",
		UseCache:         false,
		Cachedir:         ".buildsh/cache",
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
	if d := os.Getenv("BUILDSH_HOME"); d != "" {
		c.Home = d
	}
	if d := os.Getenv("BUILDSH_CONTAINER_HOME"); d != "" {
		c.ContainerHome = d
	}
	if d := os.Getenv("BUILDSH_CONTAINER_WORKDIR"); d != "" {
		c.ContainerWorkdir = d
	}
	if d := os.Getenv("BUILDSH_HOME_IN_CONTAINER"); d != "" {
		c.HomeInContainer = d
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

func shellEscape(s string) string {
	return "'" + strings.Replace(s, "'", "'\"'\"'", -1) + "'"
}

var realEntrypointTemplate = `
set -e

# Workaround to use 'sudo' with arbitrary user id that is specified host machine.
# You should set '-e' docker run option like the following:
#   -e BUILDSH_USER="<user_id>:<group_id>"
if [ -n "$BUILDSH_USER" ]; then
    # split user_id and group_id
    OLD_IFS="$IFS"
    IFS=:
    arr=($BUILDSH_USER)
    IFS="$OLD_IFS"

    if [ ${#arr[@]} -ne 2 ]; then
        echo "'BUILDSH_USER' must be formatted '<user_id>:<group_id>', but $BUILDSH_USER" 1>&2
        exit 1
    fi

    # Create buildbot user
    groupadd --non-unique --gid ${arr[1]} buildbot
    useradd --non-unique --uid ${arr[0]} --gid ${arr[1]} buildbot

    if type sudo >/dev/null 2>&1; then
        echo 'buildbot	ALL=(ALL) NOPASSWD:ALL' >> /etc/sudoers
    fi

    {{ if .Cmd }}
    exec su buildbot -c {{ .Cmd | ShellEscape}}
    {{ else }}
    exec su buildbot
    {{ end }}
else
    {{ if .Cmd }}
    exec {{ .Cmd }}
    {{ else }}
    exec /bin/bash
    {{ end }}
fi
`
