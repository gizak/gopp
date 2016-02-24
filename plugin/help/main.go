package goplugin

import (
	"fmt"
	"os/exec"
)

type help struct {
	text string
}

var PLUGIN = help{`Go is a tool for managing Go source code.

This is the gopp version of go command with some additional features.
Vanilla go commands are expecting to work as usual.

Usage:

	gopp command [arguments]

The commands are:

	build       compile packages and dependencies
	clean       remove object files
	doc         show documentation for package or symbol
	env         print Go environment information
	fix         run go tool fix on packages
	fmt         run gofmt on package sources
	generate    generate Go files by processing source
	get         download and install packages and dependencies
	install     compile and install packages and dependencies
	list        list packages
	plugin      manage plugins
	rm          remove a installed pacakge
	run         compile and run Go program
	test        test packages
	tool        run specified go tool
	version     print Go version
	vet         run go tool vet on packages

Use "go help [command]" for more information about a command.

Additional help topics:

	c           calling between Go and C
	buildmode   description of build modes
	filetype    file types
	gopath      GOPATH environment variable
	environment environment variables
	importpath  import path syntax
	packages    description of package lists
	testflag    description of testing flags
	testfunc    description of testing functions

Use "gopp help [topic]" for more information about that topic.
`}

func (p help) Subcmd() string {
	return "help"
}

func (p help) Usage() string {
	return p.text
}

var cmds = []string{
	"plugin",
	"rm",
}

func (p help) RunSubcmd(args []string) error {
	if len(args) == 1 {
		fmt.Print(p.text)
		return nil
	}

	cmd := args[1]
	has := false

	for _, c := range cmds {
		if c == cmd {
			has = true
			break
		}
	}

	if has {
		b, err := exec.Command("gopp", args[1], "--help").CombinedOutput()
		fmt.Print(string(b))
		return err
	}

	// original go cmds will fall through
	b, err := exec.Command("go", args...).CombinedOutput()
	fmt.Print(string(b))
	return err
}
