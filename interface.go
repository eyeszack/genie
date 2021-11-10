package geenee

import (
	"io"
	"os"
	"strings"
)

type GeeneeError string

func (c GeeneeError) Error() string { return string(c) }

var (
	ErrNoOp                 = GeeneeError("noop")
	ErrInvalidInterfaceName = GeeneeError("invalid interface name provided")
	ErrNoArgs               = GeeneeError("arguments not provided")
	ErrInvalidSyntax        = GeeneeError("interface/command syntax was invalid")
	ErrCommandDepthInvalid  = GeeneeError("invalid command depth")
	ErrCommandNotFound      = GeeneeError("command not found")
	ErrCommandNotRunnable   = GeeneeError("command not runnable")
)

//Interface is a very simple representation of a Command Line Interface or any Interface with commands.
type Interface struct {
	Name        string
	RootCommand *Command
	Out         io.Writer
	Err         io.Writer
	Version     string
}

//NewInterface returns an Interface with sensible defaults.
func NewInterface(name, version string, withRoot bool) *Interface {
	i := &Interface{
		Name:    name,
		Out:     os.Stdout,
		Err:     os.Stderr,
		Version: version,
	}

	if withRoot {
		i.RootCommand = NewCommand(name)
	}

	return i
}

//Exec executes the Interface with the provided arguments.
func (c *Interface) Exec(args []string) error {
	//if we have no root command there is nothing we can do
	if c.RootCommand == nil {
		return ErrNoOp
	}

	if len(args) <= 0 {
		return ErrNoArgs
	}

	if args[0] != c.Name || strings.HasPrefix(args[0], "-") {
		return ErrInvalidInterfaceName
	}

	//just the interface was provided so run root command if available
	if len(args) == 1 {
		return c.RootCommand.run([]string{})
	}

	//there are more than enough args to check for flags etc....

	flagStart, hasFlags := c.hasFlags(args)

	//it's ok to not have flags, but if we do we assume a little structure
	if hasFlags {
		//flags if provided must immediately follow interface/command/subcommand
		if flagStart > 3 {
			return ErrCommandDepthInvalid
		}

		switch flagStart {
		case 1:
			//calling interface
			return c.RootCommand.run(args[flagStart:])
		case 2:
			//calling a command
			command, found := c.RootCommand.findSubCommand(args[1])
			if !found {
				return ErrCommandNotFound
			}

			return command.run(args[flagStart:])
		case 3:
			//calling a subcommand
			command, found := c.RootCommand.findSubCommand(args[1])
			if !found {
				return ErrCommandNotFound
			}

			subcommand, found := command.findSubCommand(args[2])
			if !found {
				return ErrCommandNotFound
			}

			return subcommand.run(args[flagStart:])
		}
	}

	//------
	//no flags were provided so we run things a little differently, less structure
	//-----

	//we may have a command provided
	command, found := c.RootCommand.findSubCommand(args[1])
	if found {
		//let's check for a subcommand
		if len(args) > 2 {
			subcommand, found := command.findSubCommand(args[2])
			if found {
				return subcommand.run(args[3:])
			}
		}

		//no subcommand found or requested
		return command.run(args[2:])
	}

	return c.RootCommand.run(args[1:])
}

func (c *Interface) hasFlags(args []string) (int, bool) {
	for i, flag := range args {
		if strings.HasPrefix(flag, "-") {
			return i, true
		}
	}

	return -1, false
}
