package geenee

import (
	"flag"
	"fmt"
	"io"
	"os"
)

//CheckFunc is used to check stuff before run is called, if error is returned run will not be called.
type CheckFunc func(command *Command) error

//RunFunc is the function that will be called when the command is run.
type RunFunc func(command *Command) error

//UsageFunc is the function that will be called when the command is run.
type UsageFunc func(command *Command) string

//Command represents a command or subcommand of the interface.
type Command struct {
	Name           string
	Aliases        []string
	RunSyntax      string
	Description    string
	Path           string
	ExtraInfo      string
	Flags          *flag.FlagSet
	SubCommands    []*Command
	Out            io.Writer
	Err            io.Writer
	Check          CheckFunc
	Run            RunFunc
	Usage          UsageFunc
	MergeFlagUsage bool
	Secret         bool
}

//NewCommand returns a Command with sensible defaults.
func NewCommand(name string) *Command {
	c := &Command{
		Name:  name,
		Path:  name, //default path to name, but this will be altered by AddCommand method
		Flags: flag.NewFlagSet(name, flag.ExitOnError),
		Out:   os.Stdout,
		Err:   os.Stderr,
		Usage: DefaultCommandUsageFunc,
	}
	c.Flags.Usage = func() {
		fmt.Fprint(c.Err, DefaultCommandUsageFunc(c))
	}

	return c
}

//AddSubCommand adds the command to the available subcommands. This will also adjust the path to the added command.
func (c *Command) AddSubCommand(command *Command) {
	command.Path = fmt.Sprintf("%s %s", c.Path, command.Name)
	c.SubCommands = append(c.SubCommands, command)
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

		for _, alias := range command.Aliases {
			if name == alias {
				return command, true
			}
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
