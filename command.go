package geenee

import (
	"flag"
	"io"
)

//CheckFunc is used to check stuff before run is called, if error is returned run will not be called.
type CheckFunc func(command *Command) error

//RunFunc is the function that will be called when the command is run.
type RunFunc func(command *Command) error

//Command represents a command or subcommand of the interface.
type Command struct {
	Name        string
	Flags       *flag.FlagSet
	SubCommands []*Command
	Out         io.Writer
	Err         io.Writer
	Check       CheckFunc
	Run         RunFunc
}

//FlagWasProvided returns true if the flag was actually provided at execution time.
func (c *Command) FlagWasProvided(name string) bool {
	if c.Flags == nil {
		return false
	}

	set := false
	c.Flags.Visit(func(f *flag.Flag) {
		if name == f.Name {
			set = true
		}
	})

	return set
}

func (c *Command) findSubCommand(name string) (*Command, bool) {
	for _, command := range c.SubCommands {
		if name == command.Name {
			return command, true
		}
	}

	return nil, false
}

func (c *Command) run(args []string) error {
	if c.Run == nil {
		return ErrCommandNotRunnable
	}
	if c.Flags != nil {
		c.Flags.Parse(args)
	}

	if c.Check != nil {
		err := c.Check(c)
		if err != nil {
			return err
		}
	}

	return c.Run(c)
}
