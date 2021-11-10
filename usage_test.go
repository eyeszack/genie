package geenee

import (
	"bytes"
	"flag"
	"io/ioutil"
	"testing"
	"time"
)

func Test_DefaultUsage(t *testing.T) {
	t.Run("validate default usage", func(t *testing.T) {
		want := `The test command is for testing.

USAGE:
test command [-flags...] [args...]

HEADING:
This can be used to provide all kinds of extra usage info.

FLAGS:
  --count int
		this is an int count (default 100)
  -n uint
		this is a uint flag (default 0)
  --price float
		this is a float flag (default 1.5)
  --testing string
		this is a testing flag
  --time duration
		this is a duration flag (default 1h0m0s)
  -v
		this is another testing flag (default false)

SUBCOMMANDS:
subcommand	The test command subcommand.
`
		subject := &Command{
			Name:        "command",
			RunSyntax:   "test command [-flags...] [args...]",
			Description: "The test command is for testing.",
			ExtraInfo: `HEADING:
This can be used to provide all kinds of extra usage info.`,
			Flags: flag.NewFlagSet("command", flag.ExitOnError),
			SubCommands: []*Command{
				{
					Name:        "subcommand",
					RunSyntax:   "test command subcommand [-flags...] [args...]",
					Description: "The test command subcommand.",
				},
				{
					Name:        "shhh",
					RunSyntax:   "test command shhh [-flags...] [args...]",
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

		got := subject.Usage(subject)
		if got != want {
			t.Errorf("want: %s got: %s", want, got)
		}
	})

	t.Run("validate default usage - merged flags", func(t *testing.T) {
		want := `The test command is for testing.

USAGE:
test command [-flags...] [args...]

HEADING:
This can be used to provide all kinds of extra usage info.

FLAGS:
  --count -c int
		this is an int count (default 100)
  --number -n uint
		this is a uint flag (default 0)
  --price -p float
		this is a float flag (default 1.5)
  --testing -t string
		this is a testing flag
  --time -d duration
		this is a duration flag (default 1h0m0s)
  --verbose -v
		this is another testing flag (default false)

SUBCOMMANDS:
subcommand	The test command subcommand.
`
		subject := &Command{
			Name:        "command",
			RunSyntax:   "test command [-flags...] [args...]",
			Description: "The test command is for testing.",
			ExtraInfo: `HEADING:
This can be used to provide all kinds of extra usage info.`,
			Flags: flag.NewFlagSet("command", flag.ExitOnError),
			SubCommands: []*Command{
				{
					Name:        "subcommand",
					RunSyntax:   "test command subcommand [-flags...] [args...]",
					Description: "The test command subcommand.",
				},
			},
			Usage:          DefaultCommandUsageFunc,
			MergeFlagUsage: true,
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

		got := subject.Usage(subject)
		if got != want {
			t.Errorf("want: %s got: %s", want, got)
		}
	})

	t.Run("validate default usage minimal", func(t *testing.T) {
		want := `The test command is for testing.

USAGE:
command
`
		subject := &Command{
			Name:        "command",
			Description: "The test command is for testing.",
			Usage:       DefaultCommandUsageFunc,
		}

		got := subject.Usage(subject)
		if got != want {
			t.Errorf("want: %s got: %s", want, got)
		}
	})

	t.Run("validate default usage - flags error", func(t *testing.T) {
		b := bytes.NewBufferString("")
		want := `The test command is for testing.

USAGE:
command

FLAGS:
  --testing string
		this is a testing flag
`
		subject := NewCommand("command")
		subject.Description = "The test command is for testing."
		subject.Err = b
		subject.Flags.String("testing", "", "this is a testing flag")
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
`
		subject := NewCommand("command")
		subject.Description = "The test command is for testing."
		subject.SubCommands = []*Command{{Name: "supersecret", Description: "I'm not known.", Secret: true}}
		got := subject.Usage(subject)
		if string(got) != want {
			t.Errorf("want: %s, got %s", want, string(got))
		}
	})
}
