package genie

import (
	"fmt"
	"io"
	"os"
	"strings"
)

type Error string

func (e Error) Error() string { return string(e) }

var (
	ErrNoOp                = Error("noop")
	ErrNoArgs              = Error("arguments not provided")
	ErrCommandDepthInvalid = Error("invalid command depth")
	ErrCommandNotFound     = Error("command not found")
	ErrCommandNotRunnable  = Error("command not runnable")
)

//Lamp is a very simple representation of a Command Line Interface or any interface with commands.
type Lamp struct {
	Name            string
	RootCommand     *Command
	Out             io.Writer
	Err             io.Writer
	Version         string
	SilenceFlags    bool
	MaxCommandDepth int
}

//NewLamp returns a Lamp with sensible defaults.
func NewLamp(name, version string, silenceFlags bool) *Lamp {
	return &Lamp{
		Name:            name,
		RootCommand:     NewCommand(name, silenceFlags),
		Out:             os.Stdout,
		Err:             os.Stderr,
		Version:         version,
		SilenceFlags:    silenceFlags,
		MaxCommandDepth: 3,
	}
}

//SetWriters will set the out and err writers on the interface and all commands.
//If a writer is nil it will not change current writer.
func (l *Lamp) SetWriters(o, e io.Writer) {
	if o != nil {
		l.Out = o
		l.RootCommand.SetOut(o)
	}

	if e != nil {
		l.Err = e
		l.RootCommand.SetErr(e)
	}
}

//Execute will execute the Lamp with os.Args as the provided arguments, returns the command executed if found.
func (l *Lamp) Execute() (*Command, error) {
	return l.ExecuteWith(os.Args)
}

//ExecuteWith executes the Lamp with the provided arguments, returns the command executed if found.
func (l *Lamp) ExecuteWith(args []string) (*Command, error) { //all: os.Args() = lamp command command -flag value -flag2 value2 arg1 arg2
	//if we have no root command there is nothing we can do
	if l.RootCommand == nil {
		return nil, ErrNoOp
	}

	//set root to true since we know for sure this is the root command, and anchor the paths from root
	l.RootCommand.root = true
	l.RootCommand.AnchorPaths()

	//no args will return an error, but some folks may not care
	if len(args) <= 0 {
		return nil, ErrNoArgs
	}

	//this could happen if the full path to binary used, or command acts as entrypoint to multiple lamps
	//in either case we want the first arg to reflect the lamp being used
	if args[0] != l.Name || strings.HasPrefix(args[0], "-") { //HRM: why did I allow for "-" as prefix?
		args[0] = l.Name
	}

	//just the lamp was provided so run root command
	if len(args) == 1 {
		return l.RootCommand, l.RootCommand.run([]string{})
	}

	//------------------------------------------------------------------------
	//at this point there are more than enough args to check for flags etc....
	//------------------------------------------------------------------------

	flagStart, hasFlags := l.hasFlags(args)

	//it's ok to not have flags, but if we do we assume a little structure
	if hasFlags {
		if flagStart > l.MaxCommandDepth {
			return nil, ErrCommandDepthInvalid
		}

		switch flagStart {
		case 1:
			if askedForVersion(args) {
				if l.Out != nil {
					_, _ = fmt.Fprintln(l.Out, l.Version)
				}
				return l.RootCommand, nil
			}
			//calling interface
			return l.RootCommand, l.RootCommand.run(args[flagStart:])
		default:
			command, found, _ := l.searchPathForCommand(args[1:flagStart], false)
			if !found {
				return nil, ErrCommandNotFound
			}

			command.root = false
			return command, command.run(args[flagStart:])
		}
	}

	//----------------------------------------------------------------------------------
	//no flags were provided, so we run things a little differently, with less structure
	//----------------------------------------------------------------------------------

	//we may have a command provided
	command, found, position := l.searchPathForCommand(args[1:], true)
	if found {
		if position+2 > l.MaxCommandDepth {
			return nil, ErrCommandDepthInvalid
		}
		command.root = false
		return command, command.run(args[position+2:])
	}

	//no command or subcommand found so run root
	return l.RootCommand, l.RootCommand.run(args[1:])
}

//TraverseCommands visits each command and its subcommands, and calls do with each command.
func (l *Lamp) TraverseCommands(do func(command *Command)) {
	l.RootCommand.AnchorPaths()
	do(l.RootCommand)
	for _, sc := range l.RootCommand.SubCommands {
		traverse(sc, do)
	}
}

func traverse(c *Command, do func(command *Command)) {
	do(c)
	for _, sc := range c.SubCommands {
		traverse(sc, do)
	}
}

//returns true if flags were found, and the start position of the first flag
func (l *Lamp) hasFlags(args []string) (int, bool) {
	for i, flag := range args {
		if strings.HasPrefix(flag, "-") {
			return i, true
		}
	}

	return -1, false
}

//path should not contain the interface name
//if partialAllowed == false, path should only contain commands, no flags/args (e.g. command subcommand)
//if partialAllowed == true, path can contain trailing flags/args as it will return last command found if any (e.g. command subcommand -flag value -flag2 value2 arg1 arg2)
func (l *Lamp) searchPathForCommand(path []string, partialAllowed bool) (*Command, bool, int) {
	var lastFoundCommand *Command
	lastFoundResult := false
	lastFoundAt := -1
	for i, pathPart := range path {
		if i == 0 {
			temp, found := l.RootCommand.findSubCommand(pathPart)
			//if very first element results in not found it makes no sense to continue at all
			if !found {
				return nil, false, -1
			}
			lastFoundCommand = temp
			lastFoundResult = found
			lastFoundAt = i
		} else {
			temp, found := lastFoundCommand.findSubCommand(pathPart)
			if !found {
				//now that we've at least found something, we check if a partial find was requested, if so return last found results
				if partialAllowed {
					break
				}
				return nil, false, -1
			}
			lastFoundCommand = temp
			lastFoundResult = found
			lastFoundAt = i
		}
	}

	return lastFoundCommand, lastFoundResult, lastFoundAt
}

func askedForVersion(args []string) bool {
	for _, flag := range args {
		if flag == "-version" || flag == "--version" {
			return true
		}
	}

	return false
}
