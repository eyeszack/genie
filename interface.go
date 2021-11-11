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
	Name         string
	RootCommand  *Command
	Out          io.Writer
	Err          io.Writer
	Version      string
	SilenceFlags bool
}

//NewInterface returns an Interface with sensible defaults.
func NewInterface(name, version string, silenceFlags bool) *Interface {
	return &Interface{
		Name:         name,
		RootCommand:  NewCommand(name, silenceFlags),
		Out:          os.Stdout,
		Err:          os.Stderr,
		Version:      version,
		SilenceFlags: silenceFlags,
	}
}

//SetWriters will set the out and err writers on the interface and all commands.
//If either writer is nil this will not change current writer.
func (c *Interface) SetWriters(o, e io.Writer) {
	if o != nil {
		c.Out = o
		c.RootCommand.SetOut(o)
	}

	if e != nil {
		c.Err = e
		c.RootCommand.SetErr(e)
	}
}

//Exec executes the Interface with the provided arguments.
func (c *Interface) Exec(args []string) error {
	//if we have no root command there is nothing we can do
	if c.RootCommand == nil {
		return ErrNoOp
	}

	//no args will return an error, but some folks may not care
	if len(args) <= 0 {
		return ErrNoArgs
	}

	//this could happen if the root command is misnamed, maybe in the future we can ignore
	if args[0] != c.Name || strings.HasPrefix(args[0], "-") {
		return ErrInvalidInterfaceName
	}

	//just the interface was provided so run root command if available
	if len(args) == 1 {
		return c.RootCommand.run([]string{})
	}

	//------------------------------------------------------------------------
	//at this point there are more than enough args to check for flags etc....
	//------------------------------------------------------------------------

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
		default:
			command, found := c.searchPathForCommand(args[1:flagStart])
			if !found {
				return ErrCommandNotFound
			}

			return command.run(args[flagStart:])
			/*case 2:
				//calling a command (technically a subcommand on root)
				command, found := c.RootCommand.findSubCommand(args[1])
				if !found {
					return ErrCommandNotFound
				}

				return command.run(args[flagStart:])
			case 3:
				//calling a subcommand so first find command
				command, found := c.RootCommand.findSubCommand(args[1])
				if !found {
					return ErrCommandNotFound
				}

				//then find the subcommand
				subcommand, found := command.findSubCommand(args[2])
				if !found {
					return ErrCommandNotFound
				}

				return subcommand.run(args[flagStart:])*/
		}
	}

	//----------------------------------------------------------------------------
	//no flags were provided so we run things a little differently, less structure
	//----------------------------------------------------------------------------

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

		//no subcommand found or requested, so run the command we found
		return command.run(args[2:])
	}

	//no command or subcommand found so run root
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

func (c *Interface) searchPathForCommand(path []string) (*Command, bool) { //this is path excluding the interface name
	var lastFound *Command
	found := false
	for i, pathPart := range path {
		if i == 0 {
			lastFound, found = c.RootCommand.findSubCommand(pathPart)
			if !found {
				return nil, false
			}
		} else {
			lastFound, found = lastFound.findSubCommand(pathPart)
			if !found {
				return nil, false
			}
		}
	}

	return lastFound, found
}
