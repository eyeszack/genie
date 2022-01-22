package geenee

import (
	"bytes"
	"flag"
	"fmt"
	"io/ioutil"
	"testing"
	"time"
)

func Test_DefaultMarkedUsage(t *testing.T) {
	t.Run("validate default marked usage", func(t *testing.T) {
		tf := testFlag{}
		want := `::DESCRIPTION::The test command is for testing.::DESCRIPTION-END::

::HEADER::USAGE:::HEADER-END::
command [-flags...] [args...]

::HEADER::HEADING:::HEADER-END::
This can be used to provide all kinds of extra usage info.

::HEADER::HEADING2:::HEADER-END::
Good stuff!

::HEADER::FLAGS:::HEADER-END::
::FLAG::--count::FLAG-END::       int         this is an int count (default 100)
::FLAG::--help::FLAG-END::                    display help for command
::FLAG::--price::FLAG-END::       float       this is a float flag (default 1.5)
::FLAG::--testing::FLAG-END::     string      this is a testing flag
::FLAG::--time::FLAG-END::        duration    this is a duration flag (default 1h0m0s)
::FLAG::-n::FLAG-END::            uint        this is a uint flag (default 0)
::FLAG::-v::FLAG-END::                        this is another testing flag (default false)
::FLAG::-z::FLAG-END::            testFlag    this is a custom flag

::HEADER::COMMANDS:::HEADER-END::
::SUBCMD::subcommand::SUBCMD-END::    The test command subcommand.

Use "command <command> --help" for more information.
`
		subject := &Command{
			Name:        "command",
			RunSyntax:   "[-flags...] [args...]",
			Description: "The test command is for testing.",
			ExtraInfo: `HEADING:
This can be used to provide all kinds of extra usage info.

HEADING2:
Good stuff!`,
			Flags: flag.NewFlagSet("command", flag.ExitOnError),
			SubCommands: []*Command{
				{
					Name:        "subcommand",
					RunSyntax:   "[-flags...] [args...]",
					Description: "The test command subcommand.",
				},
				{
					Name:        "shhh",
					RunSyntax:   "[-flags...] [args...]",
					Description: "shhh.",
					Secret:      true,
				},
			},
			Usage: DefaultCommandUsageMarkedFunc,
		}

		subject.Flags.String("testing", "", "this is a testing flag")
		subject.Flags.Bool("v", false, "this is another testing flag")
		subject.Flags.Int("count", 100, "this is an int count")
		subject.Flags.Float64("price", 1.5, "this is a float flag")
		subject.Flags.Duration("time", time.Hour, "this is a duration flag")
		subject.Flags.Uint("n", 0, "this is a uint flag")
		subject.Flags.Var(&tf, "z", "this is a custom flag")
		subject.AnchorPaths()

		got := subject.ShowUsage()
		if got != want {
			t.Errorf("want: %s got: %s", want, got)
		}
	})

	t.Run("validate default marked usage - merged", func(t *testing.T) {
		tf := testFlag{}
		want := `::DESCRIPTION::The test command is for testing.::DESCRIPTION-END::

::HEADER::USAGE:::HEADER-END::
command [-flags...] [args...]

::HEADER::HEADING:::HEADER-END::
This can be used to provide all kinds of extra usage info.

::HEADER::FLAGS:::HEADER-END::
::FLAG::--count::FLAG-END::       int         this is an int count (default 100)
::FLAG::--help::FLAG-END::                    display help for command
::FLAG::--price::FLAG-END::       float       this is a float flag (default 1.5)
::FLAG::--testing::FLAG-END::     string      this is a testing flag
::FLAG::--time::FLAG-END::        duration    this is a duration flag (default 1h0m0s)
::FLAG::-n::FLAG-END::            uint        this is a uint flag (default 0)
::FLAG::-v::FLAG-END::                        this is another testing flag (default false)
::FLAG::-z::FLAG-END::            testFlag    this is a custom flag

::HEADER::COMMANDS:::HEADER-END::
::SUBCMD::subcommand::SUBCMD-END::    The test command subcommand.

Use "command <command> --help" for more information.
`
		subject := &Command{
			Name:        "command",
			RunSyntax:   "[-flags...] [args...]",
			Description: "The test command is for testing.",
			ExtraInfo: `HEADING:
This can be used to provide all kinds of extra usage info.`,
			Flags: flag.NewFlagSet("command", flag.ExitOnError),
			SubCommands: []*Command{
				{
					Name:        "subcommand",
					RunSyntax:   "[-flags...] [args...]",
					Description: "The test command subcommand.",
				},
				{
					Name:        "shhh",
					RunSyntax:   "[-flags...] [args...]",
					Description: "shhh.",
					Secret:      true,
				},
			},
			Usage: DefaultCommandUsageMarkedFunc,
		}

		subject.MergeFlagUsage = true
		subject.Flags.String("testing", "", "this is a testing flag")
		subject.Flags.Bool("v", false, "this is another testing flag")
		subject.Flags.Int("count", 100, "this is an int count")
		subject.Flags.Float64("price", 1.5, "this is a float flag")
		subject.Flags.Duration("time", time.Hour, "this is a duration flag")
		subject.Flags.Uint("n", 0, "this is a uint flag")
		subject.Flags.Var(&tf, "z", "this is a custom flag")
		subject.AnchorPaths()

		got := subject.ShowUsage()
		if got != want {
			t.Errorf("want: %s got: %s", want, got)
		}
	})

	t.Run("validate default marked usage - merged flags and path set", func(t *testing.T) {
		want := `::DESCRIPTION::The test command is for testing.::DESCRIPTION-END::

::HEADER::USAGE:::HEADER-END::
command [-flags...] [args...]

::HEADER::HEADING:::HEADER-END::
This can be used to provide all kinds of extra usage info.

::HEADER::FLAGS:::HEADER-END::
::FLAG::--count -c::FLAG-END::       int         this is an int count (default 100)
::FLAG::--help::FLAG-END::                       display help for command
::FLAG::--number -n::FLAG-END::      uint        this is a uint flag (default 0)
::FLAG::--price -p::FLAG-END::       float       this is a float flag (default 1.5)
::FLAG::--testing -t::FLAG-END::     string      this is a testing flag
::FLAG::--time -d::FLAG-END::        duration    this is a duration flag (default 1h0m0s)
::FLAG::--verbose -v::FLAG-END::                 this is another testing flag (default false)
::FLAG::--version::FLAG-END::                    display version information

::HEADER::COMMANDS:::HEADER-END::
::SUBCMD::subcommand::SUBCMD-END::    The test command subcommand.

Use "command <command> --help" for more information.
`
		subject := &Command{
			Name:        "command",
			RunSyntax:   "[-flags...] [args...]",
			Description: "The test command is for testing.",
			ExtraInfo: `HEADING:
This can be used to provide all kinds of extra usage info.`,
			Flags: flag.NewFlagSet("command", flag.ExitOnError),
			SubCommands: []*Command{
				{
					Name:        "subcommand",
					RunSyntax:   "[-flags...] [args...]",
					Description: "The test command subcommand.",
					root:        false,
				},
			},
			Usage:          DefaultCommandUsageMarkedFunc,
			MergeFlagUsage: true,
			root:           true,
		}

		subject.Flags.String("testing", "", "this is a testing flag")
		subject.Flags.Bool("v", false, "this is another testing flag")
		subject.Flags.Int("count", 100, "this is an int count")
		subject.Flags.Float64("price", 1.5, "this is a float flag")
		subject.Flags.Duration("time", time.Hour, "this is a duration flag")
		subject.Flags.Uint("n", 0, "this is a uint flag")
		subject.Flags.String("t", "", "this is a testing flag")
		subject.Flags.Bool("verbose", false, "this is another testing flag")
		subject.Flags.Int("c", 100, "this is an int count")
		subject.Flags.Float64("p", 1.5, "this is a float flag")
		subject.Flags.Duration("d", time.Hour, "this is a duration flag")
		subject.Flags.Uint("number", 0, "this is a uint flag")
		subject.AnchorPaths()

		got := subject.ShowUsage()
		if got != want {
			t.Errorf("want: %s got: %s", want, got)
		}
	})

	t.Run("validate default marked usage minimal", func(t *testing.T) {
		want := `::DESCRIPTION::The test command is for testing.::DESCRIPTION-END::

::HEADER::USAGE:::HEADER-END::
command

::HEADER::FLAGS:::HEADER-END::
::FLAG::--help::FLAG-END::        display help for command
`
		subject := &Command{
			Name:        "command",
			Description: "The test command is for testing.",
			Usage:       DefaultCommandUsageMarkedFunc,
		}

		got := subject.ShowUsage()
		if got != want {
			t.Errorf("want: %s got: %s", want, got)
		}
	})

	t.Run("validate default usage minimal - nil Usage", func(t *testing.T) {
		want := `The test command is for testing.

USAGE:
command

FLAGS:
--help        display help for command
`
		subject := &Command{
			Name:        "command",
			Description: "The test command is for testing.",
		}

		got := subject.ShowUsage()
		if got != want {
			t.Errorf("want: %s got: %s", want, got)
		}
	})

	t.Run("validate default marked usage - flags error", func(t *testing.T) {
		b := bytes.NewBufferString("")
		want := `::DESCRIPTION::The test command is for testing.::DESCRIPTION-END::

::HEADER::USAGE:::HEADER-END::
command --testing heyo
command --slice 2,3 <test_string>...

::HEADER::FLAGS:::HEADER-END::
::FLAG::--help::FLAG-END::                  display help for command
::FLAG::--testing::FLAG-END::     string    this is a testing flag

::HEADER::ARGUMENTS:::HEADER-END::
One or more test strings can be supplied.
`
		subject := NewCommand("command", false)
		subject.RunSyntax = "--testing heyo\n{{path}} --slice 2,3 <test_string>..."
		subject.Description = "The test command is for testing."
		subject.ArgInfo = "One or more test strings can be supplied."
		subject.Err = b
		subject.Flags.String("testing", "", "this is a testing flag")
		subject.Flags.String("hideme", "", "i should not show up")
		subject.SecretFlag("hideme")
		subject.Flags.Usage = func() {
			fmt.Fprint(subject.Err, DefaultCommandUsageMarkedFunc(subject))
		}
		subject.Flags.Usage()

		got, err := ioutil.ReadAll(b)
		if err != nil {
			t.Errorf("[err] want nil, got %s", err)
		}
		if string(got) != want {
			t.Errorf("want: %s, got %s", want, string(got))
		}
	})

	t.Run("validate default marked usage - empty flags & secret command", func(t *testing.T) {
		want := `::DESCRIPTION::The test command is for testing.::DESCRIPTION-END::

::HEADER::USAGE:::HEADER-END::
command

::HEADER::ALIASES:::HEADER-END::
cmd

::HEADER::FLAGS:::HEADER-END::
::FLAG::--help::FLAG-END::        display help for command
`
		subject := NewCommand("command", false)
		subject.Aliases = []string{"cmd"}
		subject.Description = "The test command is for testing."
		subject.SubCommands = []*Command{{Name: "supersecret", Description: "I'm not known.", Secret: true}}
		subject.Usage = DefaultCommandUsageMarkedFunc
		got := subject.ShowUsage()
		if string(got) != want {
			t.Errorf("want: %s, got %s", want, string(got))
		}
	})

	t.Run("validate default marked usage - on root command no flags", func(t *testing.T) {
		want := `::DESCRIPTION::The test command is for testing.::DESCRIPTION-END::

::HEADER::USAGE:::HEADER-END::
command

::HEADER::FLAGS:::HEADER-END::
::FLAG::--help::FLAG-END::           display help for command
::FLAG::--version::FLAG-END::        display version information
`
		subject := NewCommand("command", false)
		subject.Description = "The test command is for testing."
		subject.root = true
		subject.Usage = DefaultCommandUsageMarkedFunc
		got := subject.ShowUsage()
		if string(got) != want {
			t.Errorf("want: %s, got %s", want, string(got))
		}
	})

	t.Run("validate default marked usage - on root command with flags", func(t *testing.T) {
		want := `::DESCRIPTION::The test command is for testing.::DESCRIPTION-END::

::HEADER::USAGE:::HEADER-END::
command

::HEADER::FLAGS:::HEADER-END::
::FLAG::--help::FLAG-END::                     display help for command
::FLAG::--testing -t::FLAG-END::     string    this is a testing flag
::FLAG::--version::FLAG-END::                  display version information
`
		subject := NewCommand("command", false)
		subject.Description = "The test command is for testing."
		subject.Flags.String("testing", "", "this is a testing flag")
		subject.Flags.String("t", "", "this is a testing flag")
		subject.Flags.String("hideme", "", "i should not show up")
		subject.Flags.String("m", "", "i should not show up")
		subject.SecretFlag("hideme")
		subject.SecretFlag("m")
		subject.MergeFlagUsage = true
		subject.root = true
		subject.Usage = DefaultCommandUsageMarkedFunc
		got := subject.ShowUsage()
		if string(got) != want {
			t.Errorf("want: %s, got %s", want, string(got))
		}
	})
}
