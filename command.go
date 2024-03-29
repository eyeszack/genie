package genie

import (
	"flag"
	"fmt"
	"io"
	"os"
	"strings"
)

// PipedInFunc will be called if stdin was piped in to this command, runs before check. If error is returned command
// execution stops here.
type PipedInFunc func(command *Command) (error, bool)

// CheckFunc is used to check stuff before run is called, if error is returned command execution stops here.
type CheckFunc func(command *Command) error

// RunFunc is the function that will be called when the command is run.
type RunFunc func(command *Command) error

// UsageFunc is the function that will be called when the command is run.
type UsageFunc func(command *Command) string

// NoopUsage if used for usage silence.
func NoopUsage() {}

// NoopWriter is used to silence output.
type NoopWriter struct{}

func (n *NoopWriter) Write([]byte) (int, error) {
	return 0, nil
}

// Command represents a command or subcommand of the interface.
type Command struct {
	Name           string
	Aliases        []string
	RunSyntax      string
	Description    string
	ArgInfo        string
	ExtraInfo      string
	Flags          *flag.FlagSet
	SubCommands    []*Command
	Out            io.Writer
	Err            io.Writer
	PipedIn        PipedInFunc
	Check          CheckFunc
	Run            RunFunc
	Usage          UsageFunc
	MergeFlagUsage bool
	SilenceFlags   bool
	Secret         bool
	root           bool   //this is set at execution time
	path           string //this is set at execution time
	depth          int    //this is set at execution time
	secretFlags    []string
}

// NewCommand returns a Command with sensible defaults.
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
			_, _ = fmt.Fprint(c.Err, DefaultCommandUsageFunc(c))
		}
	}

	return c
}

// SetOut sets the command's output writer, and all of it's subcommand's as well.
func (c *Command) SetOut(o io.Writer) {
	c.Out = o
	for _, sc := range c.SubCommands {
		sc.SetOut(o)
	}
}

// SetErr sets the command's error writer, flags output writer, and all of it's subcommand's as well.
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

// FlagsAndArgs returns a string representation of any flags and arguments provided to the command.
// If using the Lamp provided execution methods, this method will be ready to use in Check and Run.
func (c *Command) FlagsAndArgs() string {
	if c.Flags == nil {
		return ""
	}

	sb := strings.Builder{}
	c.Flags.Visit(func(f *flag.Flag) {
		sb.WriteString(fmt.Sprintf("%s %s ", f.Name, f.Value.String()))
	})
	if len(c.Flags.Args()) > 0 {
		for _, arg := range c.Flags.Args() {
			sb.WriteString(fmt.Sprintf("%s ", arg))
		}
	}
	s := strings.Trim(sb.String(), " ")

	return s
}

// FlagWasProvided returns true if the flag was actually provided at execution time.
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

// SecretFlag will hide the named flag from usage.
func (c *Command) SecretFlag(name string) {
	c.secretFlags = append(c.secretFlags, name)
}

// ShowUsage runs the provided usage function, or the default if none provided.
func (c *Command) ShowUsage() string {
	if c.Usage != nil {
		return c.Usage(c)
	}

	return DefaultCommandUsageFunc(c)
}

// Path returns the path to this command from the "anchor" command.
// The path will be blank until AnchorPaths is called here or on a parent command, or when command interface is executed.
func (c *Command) Path() string {
	return c.path
}

// AnchorPaths will set this command as the start of the command path.
func (c *Command) AnchorPaths() {
	c.path = c.Name
	for _, sc := range c.SubCommands {
		sc.adjustPath(c.path, c.depth+1)
	}
}

func (c *Command) adjustPath(path string, depth int) {
	c.path = fmt.Sprintf("%s %s", path, c.Name)
	c.depth = depth
	for _, sc := range c.SubCommands {
		sc.adjustPath(c.path, c.depth+1)
	}
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

func (c *Command) flagIsSecret(name string) bool {
	for _, f := range c.secretFlags {
		if f == name {
			return true
		}
	}

	return false
}

func (c *Command) run(args []string) error { //only flags/args: -flag value -flag2 value2 arg1 arg2
	return DefaultCommandRunner(c, args)
}

var DefaultCommandRunner = func(command *Command, args []string) error { //only flags/args: -flag value -flag2 value2 arg1 arg2
	if ContainsFlag("help", args) {
		if command.Out != nil {
			_, _ = fmt.Fprint(command.Out, command.ShowUsage())
			return nil
		}
		return flag.ErrHelp
	}

	if command.Run == nil {
		return ErrCommandNotRunnable
	}

	if command.Flags != nil {
		err := command.Flags.Parse(args)
		//Technically we'd not get here if flagset error handling is set to flag.ExitOnError, or flag.PanicOnError,
		//but for folks who use ContinueOnError we can return the error for custom handling if desired, so we pack it
		//in a geenee.Error for easier identification
		if err != nil {
			return Error(err.Error())
		}
	}

	fileInfo, err := os.Stdin.Stat()
	skipCheck := false
	if err == nil {
		if command.PipedIn != nil && ((fileInfo.Mode() & os.ModeCharDevice) == 0) {
			// we're receiving input (stdin) via pipe
			err, skipCheck = command.PipedIn(command)
			if err != nil {
				return err
			}
		}
	}

	if command.Check != nil && !skipCheck {
		err = command.Check(command)
		if err != nil {
			return err
		}
	}

	return command.Run(command)
}

func ContainsFlag(flag string, flags []string) bool {
	for _, f := range flags {
		if f == "-"+flag || f == "--"+flag {
			return true
		}
	}

	return false
}
