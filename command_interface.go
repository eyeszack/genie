package geenee

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

//CommandInterface is a very simple representation of a Command Line Interface or any interface with commands.
type CommandInterface struct {
	Name            string
	RootCommand     *Command
	Out             io.Writer
	Err             io.Writer
	Version         string
	SilenceFlags    bool
	MaxCommandDepth int
}

//NewCommandInterface returns a CommandInterface with sensible defaults.
func NewCommandInterface(name, version string, silenceFlags bool) *CommandInterface {
	return &CommandInterface{
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
func (ci *CommandInterface) SetWriters(o, e io.Writer) {
	if o != nil {
		ci.Out = o
		ci.RootCommand.SetOut(o)
	}

	if e != nil {
		ci.Err = e
		ci.RootCommand.SetErr(e)
	}
}

//Execute executes the CommandInterface with the provided arguments, returns the command executed if found.
func (ci *CommandInterface) Execute(args []string) (*Command, error) { //all: os.Args() = interface command command -flag value -flag2 value2 arg1 arg2
	//if we have no root command there is nothing we can do
	if ci.RootCommand == nil {
		return nil, ErrNoOp
	}

	//set root to true since we know for sure this is the root command, and anchor the paths from root
	ci.RootCommand.root = true
	ci.RootCommand.AnchorPaths()

	//no args will return an error, but some folks may not care
	if len(args) <= 0 {
		return nil, ErrNoArgs
	}

	//this could happen if the full path to binary used, or command acts as entrypoint to multiple command interfaces
	//in either case we want the first arg to reflect the command interface being used
	if args[0] != ci.Name || strings.HasPrefix(args[0], "-") { //HRM: why did I allow for "-" as prefix?
		args[0] = ci.Name
	}

	//just the interface was provided so run root command
	if len(args) == 1 {
		return ci.RootCommand, ci.RootCommand.run([]string{})
	}

	//------------------------------------------------------------------------
	//at this point there are more than enough args to check for flags etc....
	//------------------------------------------------------------------------

	flagStart, hasFlags := ci.hasFlags(args)

	//it's ok to not have flags, but if we do we assume a little structure
	if hasFlags {
		if flagStart > ci.MaxCommandDepth {
			return nil, ErrCommandDepthInvalid
		}

		switch flagStart {
		case 1:
			if askedForVersion(args) {
				if ci.Out != nil {
					fmt.Fprintln(ci.Out, ci.Version)
				}
				return ci.RootCommand, nil
			}
			//calling interface
			return ci.RootCommand, ci.RootCommand.run(args[flagStart:])
		default:
			command, found, _ := ci.searchPathForCommand(args[1:flagStart], false)
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
	command, found, position := ci.searchPathForCommand(args[1:], true)
	if found {
		if position+2 > ci.MaxCommandDepth {
			return nil, ErrCommandDepthInvalid
		}
		command.root = false
		return command, command.run(args[position+2:])
	}

	//no command or subcommand found so run root
	return ci.RootCommand, ci.RootCommand.run(args[1:])
}

//returns true if flags were found, and the start position of the first flag
func (ci *CommandInterface) hasFlags(args []string) (int, bool) {
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
func (ci *CommandInterface) searchPathForCommand(path []string, partialAllowed bool) (*Command, bool, int) {
	var lastFoundCommand *Command
	lastFoundResult := false
	lastFoundAt := -1
	for i, pathPart := range path {
		if i == 0 {
			temp, found := ci.RootCommand.findSubCommand(pathPart)
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
