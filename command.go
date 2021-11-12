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

//NoopUsage if used for usage silence.
func NoopUsage() {}

//NoopWriter is used to silence output.
type NoopWriter struct{}

func (n *NoopWriter) Write(b []byte) (int, error) {
	return 0, nil
}

//Command represents a command or subcommand of the interface.
type Command struct {
	Name           string
	Aliases        []string
	RunSyntax      string
	Description    string
	ExtraInfo      string
	Flags          *flag.FlagSet
	SubCommands    []*Command
	Out            io.Writer
	Err            io.Writer
	Check          CheckFunc
	Run            RunFunc
	Usage          UsageFunc
	MergeFlagUsage bool
	SilenceFlags   bool
	Secret         bool
}

//NewCommand returns a Command with sensible defaults.
func NewCommand(name string, silenceFlags bool) *Command {
	c := &Command{
		Name:         name,
		Out:          os.Stdout,
		Err:          os.Stderr,
		Usage:        DefaultCommandUsageFunc,
		SilenceFlags: silenceFlags,
	}

	if silenceFlags {
		//this allows for flag parsing errors to continue through to caller to handle what/if anything is output to user
		c.Flags = flag.NewFlagSet(name, flag.ContinueOnError)
		c.Flags.SetOutput(&NoopWriter{})
		c.Flags.Usage = NoopUsage
	} else {
		c.Flags = flag.NewFlagSet(name, flag.ExitOnError)
		c.Flags.Usage = func() {
			fmt.Fprint(c.Err, DefaultCommandUsageFunc(c))
		}
	}

	return c
}

//SetOut sets the command's output writer, and all of it's subcommand's as well.
func (c *Command) SetOut(o io.Writer) {
	c.Out = o
	for _, sc := range c.SubCommands {
		sc.SetOut(o)
	}
}

//SetOut sets the command's error writer, flags output writer, and all of it's subcommand's as well.
func (c *Command) SetErr(e io.Writer) {
	c.Err = e
	//if flags set and not noop via silenced, set output
	if c.Flags != nil && !c.SilenceFlags {
		c.Flags.SetOutput(e)
	}
	for _, sc := range c.SubCommands {
		sc.SetErr(e)
	}
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
		//just in case this is a call for help :)
		if askedForHelp(args) {
			return flag.ErrHelp
		}

		return ErrCommandNotRunnable
	}

	if c.Flags != nil {
		err := c.Flags.Parse(args)
		//technically we'd not get here if flagset error handling is set to flag.ExitOnError, or flag.PanicOnError,
		//but for folks who use ContinueOnError we can return the error for custom handling if desired
		if err != nil {
			if err == flag.ErrHelp && c.SilenceFlags {
				//this assumes that one of the New*(name, true) funcs was used so we can give folks a free halp command
				fmt.Fprint(c.Out, c.Usage(c))
				return nil
			}
			return err
		}
	}

	if c.Check != nil {
		err := c.Check(c)
		if err != nil {
			return err
		}
	}

	return c.Run(c)
}

func askedForHelp(args []string) bool {
	for _, flag := range args {
		if flag == "-help" || flag == "--help" || flag == "-h" || flag == "--h" {
			return true
		}
	}

	return false
}
