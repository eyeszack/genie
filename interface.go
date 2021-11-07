package geenee

import (
	"io"
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
	Name             string
	RootCommand      *Command
	Commands         []*Command
	ShowUsageOnError bool
	Out              io.Writer
	Err              io.Writer
}

//Exec executes the Interface with the provided arguments.
func (c *Interface) Exec(args []string) error {
	if len(args) <= 0 {
		return ErrNoArgs
	}

	if args[0] != c.Name || strings.HasPrefix(args[0], "-") {
		return ErrInvalidInterfaceName
	}

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
			if c.RootCommand == nil {
				return ErrInvalidSyntax
			}

			return c.RootCommand.run(args[flagStart:])
		case 2:
			//calling a command
			command, found := c.findCommand(args[1])
			if !found {
				return ErrCommandNotFound
			}

			return command.run(args[flagStart:])
		case 3:
			//calling a subcommand
			command, found := c.findCommand(args[1])
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

	//no command was provided run the root command without args if provided
	if len(args) == 1 {
		if c.RootCommand != nil {
			return c.RootCommand.run([]string{})
		}

		return ErrNoOp
	}

	//we may have a command provided
	command, found := c.findCommand(args[1])
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

	//if a command was not found check to see if there is a root command to run and pass the args
	if c.RootCommand == nil {
		return ErrNoOp //or should this be command not found?
	}

	return c.RootCommand.run(args[1:])
}

func (c *Interface) findCommand(name string) (*Command, bool) {
	for _, command := range c.Commands {
		if name == command.Name {
			return command, true
		}
	}

	return nil, false
}

func (c *Interface) hasFlags(args []string) (int, bool) {
	for i, flag := range args {
		if strings.HasPrefix(flag, "-") {
			return i, true
		}
	}

	return -1, false
}
