package genie

import (
	"bytes"
	"flag"
	"io/ioutil"
	"testing"
	"time"
)

type testFlag struct {
	value string
}

func (t *testFlag) String() string {
	return t.value
}

func (t *testFlag) Set(s string) error {
	t.value = s
	return nil
}

func (t *testFlag) Type() string {
	return "testFlag"
}

func Test_DefaultUsage(t *testing.T) {
	t.Run("validate default usage", func(t *testing.T) {
		tf := testFlag{}
		want := `The test command is for testing.

USAGE:
command [-flags...] [args...]

HEADING:
This can be used to provide all kinds of extra usage info.

FLAGS:
--count       int         this is an int count (default 100)
--help                    display help for command
--price       float       this is a float flag (default 1.5)
--testing     string      this is a testing flag
--time        duration    this is a duration flag (default 1h0m0s)
-n            uint        this is a uint flag (default 0)
-v                        this is another testing flag (default false)
-z            testFlag    this is a custom flag

COMMANDS:
subcommand    The test command subcommand.

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
			Usage: DefaultCommandUsageFunc,
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

	t.Run("validate default usage - merged", func(t *testing.T) {
		tf := testFlag{}
		want := `The test command is for testing.

USAGE:
command [-flags...] [args...]

HEADING:
This can be used to provide all kinds of extra usage info.

FLAGS:
--count       int         this is an int count (default 100)
--help                    display help for command
--price       float       this is a float flag (default 1.5)
--testing     string      this is a testing flag
--time        duration    this is a duration flag (default 1h0m0s)
-n            uint        this is a uint flag (default 0)
-v                        this is another testing flag (default false)
-z            testFlag    this is a custom flag

ARGUMENTS:
One or more test strings can be supplied.

COMMANDS:
subcommand    The test command subcommand.

Use "command <command> --help" for more information.
`
		subject := &Command{
			Name:        "command",
			RunSyntax:   "[-flags...] [args...]",
			Description: "The test command is for testing.",
			ExtraInfo: `HEADING:
This can be used to provide all kinds of extra usage info.`,
			ArgInfo: "One or more test strings can be supplied.",
			Flags:   flag.NewFlagSet("command", flag.ExitOnError),
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
			Usage: DefaultCommandUsageFunc,
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

	t.Run("validate default usage - merged flags and path set", func(t *testing.T) {
		want := `The test command is for testing.

USAGE:
command [-flags...] [args...]

HEADING:
This can be used to provide all kinds of extra usage info.

FLAGS:
--count -c       int         this is an int count (default 100)
--help                       display help for command
--number -n      uint        this is a uint flag (default 0)
--price -p       float       this is a float flag (default 1.5)
--testing -t     string      this is a testing flag
--time -d        duration    this is a duration flag (default 1h0m0s)
--verbose -v                 this is another testing flag (default false)
--version                    display version information

COMMANDS:
subcommand    The test command subcommand.

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
			Usage:          DefaultCommandUsageFunc,
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

	t.Run("validate default usage minimal", func(t *testing.T) {
		want := `The test command is for testing.

USAGE:
command

FLAGS:
--help        display help for command
`
		subject := &Command{
			Name:        "command",
			Description: "The test command is for testing.",
			Usage:       DefaultCommandUsageFunc,
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

	t.Run("validate default usage - flags error non-silenced", func(t *testing.T) {
		b := bytes.NewBufferString("")
		want := `The test command is for testing.

USAGE:
command --testing heyo
command --slice 2,3 <test_string>...

FLAGS:
--help                  display help for command
--testing     string    this is a testing flag
`
		subject := NewCommand("command", false)
		subject.RunSyntax = "--testing heyo\n{{path}} --slice 2,3 <test_string>..."
		subject.Description = "The test command is for testing."
		subject.Err = b
		subject.Flags.String("testing", "", "this is a testing flag")
		subject.Flags.String("hideme", "", "i should not show up")
		subject.SecretFlag("hideme")
		subject.Flags.Usage()

		got, err := ioutil.ReadAll(b)
		if err != nil {
			t.Errorf("[err] want nil, got %s", err)
		}
		if string(got) != want {
			t.Errorf("want: %s, got %s", want, string(got))
		}
	})

	t.Run("validate default usage - empty flags & secret command", func(t *testing.T) {
		want := `The test command is for testing.

USAGE:
command

ALIASES:
cmd

FLAGS:
--help        display help for command
`
		subject := NewCommand("command", false)
		subject.Aliases = []string{"cmd"}
		subject.Description = "The test command is for testing."
		subject.SubCommands = []*Command{{Name: "supersecret", Description: "I'm not known.", Secret: true}}
		got := subject.ShowUsage()
		if string(got) != want {
			t.Errorf("want: %s, got %s", want, string(got))
		}
	})

	t.Run("validate default usage - on root command no flags", func(t *testing.T) {
		want := `The test command is for testing.

USAGE:
command

FLAGS:
--help           display help for command
--version        display version information
`
		subject := NewCommand("command", false)
		subject.Description = "The test command is for testing."
		subject.root = true
		got := subject.ShowUsage()
		if string(got) != want {
			t.Errorf("want: %s, got %s", want, string(got))
		}
	})

	t.Run("validate default usage - on root command with flags", func(t *testing.T) {
		want := `The test command is for testing.

USAGE:
command

FLAGS:
--help                     display help for command
--testing -t     string    this is a testing flag
--version                  display version information
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
		got := subject.ShowUsage()
		if string(got) != want {
			t.Errorf("want: %s, got %s", want, string(got))
		}
	})
}

func Test_DefaultFlagsUsageFunc(t *testing.T) {
	t.Run("validate default flag usage", func(t *testing.T) {
		want := `FLAGS:
--help                  display help for command
--testing     string    this is a testing flag
--version               display version information
-t            string    this is a testing flag
`
		subject := &Command{Name: "command", Flags: flag.NewFlagSet("command", flag.ExitOnError)}
		subject.Description = "The test command is for testing."
		subject.Flags.String("testing", "", "this is a testing flag")
		subject.Flags.String("t", "", "this is a testing flag")
		subject.Flags.String("hideme", "", "i should not show up")
		subject.Flags.String("m", "", "i should not show up")
		subject.SecretFlag("hideme")
		subject.SecretFlag("m")
		subject.root = true
		got := DefaultFlagsUsageFunc(subject)
		if got != want {
			t.Errorf("want: %s, got %s", want, got)
		}
	})

	t.Run("validate default flag usage - merged", func(t *testing.T) {
		want := `FLAGS:
--help                     display help for command
--testing -t     string    this is a testing flag
--version                  display version information
`
		subject := &Command{Name: "command", Flags: flag.NewFlagSet("command", flag.ExitOnError)}
		subject.Description = "The test command is for testing."
		subject.Flags.String("testing", "", "this is a testing flag")
		subject.Flags.String("t", "", "this is a testing flag")
		subject.Flags.String("hideme", "", "i should not show up")
		subject.Flags.String("m", "", "i should not show up")
		subject.SecretFlag("hideme")
		subject.SecretFlag("m")
		subject.root = true
		subject.MergeFlagUsage = true
		got := DefaultFlagsUsageFunc(subject)
		if got != want {
			t.Errorf("want: %s, got %s", want, got)
		}
	})
}

func Test_mergeFlagsUsage(t *testing.T) {
	t.Run("validate flag usage", func(t *testing.T) {
		want := `
FLAGS:
--help                     display help for command
--testing -t     string    this is a testing flag
--version                  display version information
`
		subject := &Command{Name: "command", Flags: flag.NewFlagSet("command", flag.ExitOnError)}
		subject.Description = "The test command is for testing."
		subject.Flags.String("testing", "", "this is a testing flag")
		subject.Flags.String("t", "", "this is a testing flag")
		subject.Flags.String("hideme", "", "i should not show up")
		subject.Flags.String("m", "", "i should not show up")
		subject.SecretFlag("hideme")
		subject.SecretFlag("m")
		subject.root = true
		got := mergeFlagsUsage(subject)
		if got != want {
			t.Errorf("want: %s, got %s", want, got)
		}
	})
}

func Test_flagsUsage(t *testing.T) {
	t.Run("validate flag usage", func(t *testing.T) {
		want := `
FLAGS:
--help                  display help for command
--testing     string    this is a testing flag
--version               display version information
-t            string    this is a testing flag
`
		subject := &Command{Name: "command", Flags: flag.NewFlagSet("command", flag.ExitOnError)}
		subject.Description = "The test command is for testing."
		subject.Flags.String("testing", "", "this is a testing flag")
		subject.Flags.String("t", "", "this is a testing flag")
		subject.Flags.String("hideme", "", "i should not show up")
		subject.Flags.String("m", "", "i should not show up")
		subject.SecretFlag("hideme")
		subject.SecretFlag("m")
		subject.root = true
		got := flagsUsage(subject)
		if got != want {
			t.Errorf("want: %s, got %s", want, got)
		}
	})
}
