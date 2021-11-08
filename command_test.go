package geenee

import (
	"bytes"
	"errors"
	"flag"
	"io/ioutil"
	"os"
	"reflect"
	"testing"
	"time"
)

func Test_NewCommand(t *testing.T) {
	t.Run("validate new", func(t *testing.T) {
		wantName := "tester"
		got := NewCommand(wantName)
		if got == nil {
			t.Fatal("want command, got nil")
		}

		if got.Name != wantName {
			t.Errorf("want %s, got %s", wantName, got.Name)
		}
		if got.Flags.Name() != wantName {
			t.Errorf("want %s, got %s", wantName, got.Flags.Name())
		}
		if got.Out != os.Stdout {
			t.Errorf("want %v, got %v", os.Stdout, got.Out)
		}
		if got.Err != os.Stderr {
			t.Errorf("want %v, got %v", os.Stderr, got.Err)
		}
		if got.Usage == nil {
			t.Error("want DefaultCommandUsageFunc, got nil")
		}
	})
}

func TestCommand_FlagWasProvided(t *testing.T) {
	t.Run("validate flag is found", func(t *testing.T) {
		want := "I was provided."
		got := ""
		subject := &Command{
			Name:  "command",
			Flags: flag.NewFlagSet("command", flag.ExitOnError),
			Run: func(command *Command) error {
				if !command.FlagWasProvided("test") {
					t.Errorf("want true, got false")
				}

				if got != want {
					t.Errorf("want %s, got %s", want, got)
				}

				return nil
			},
		}
		subject.Flags.StringVar(&got, "test", "", "testing flags")

		subject.run([]string{"-test", "I was provided."})
	})

	t.Run("validate flag is not found", func(t *testing.T) {
		want := "default"
		got := ""
		subject := &Command{
			Name:  "command",
			Flags: flag.NewFlagSet("command", flag.ExitOnError),
			Run: func(command *Command) error {
				if command.FlagWasProvided("nope") {
					t.Errorf("want false, got true")
				}

				if got != want {
					t.Errorf("want %s, got %s", want, got)
				}

				return nil
			},
		}
		subject.Flags.String("test", "", "testing flags")
		subject.Flags.StringVar(&got, "nope", "default", "testing flags 2")

		subject.run([]string{"-test", "I was provided."})
	})

	t.Run("validate flag is not found when flagset is not used", func(t *testing.T) {
		want := ""
		got := ""
		subject := &Command{
			Name: "command",
			Run: func(command *Command) error {
				if command.FlagWasProvided("test") {
					t.Errorf("want false, got true")
				}

				if got != want {
					t.Errorf("want %s, got %s", want, got)
				}

				return nil
			},
		}

		subject.run([]string{"-test", "I was provided."})
	})
}

func TestCommand_findSubCommand(t *testing.T) {
	t.Run("validate subcommand is found", func(t *testing.T) {
		want := &Command{Name: "subcommand"}
		subject := &Command{
			Name: "command",
			SubCommands: []*Command{
				{
					Name: "subcommand",
				},
				{
					Name: "subcommandTwo",
				},
			},
		}

		got, found := subject.findSubCommand("subcommand")
		if !found {
			t.Error("want true, got false")
		}

		if !reflect.DeepEqual(got, want) {
			t.Errorf("want %v, got %v", want, got)
		}
	})

	t.Run("validate subcommand is not found", func(t *testing.T) {
		subject := &Command{
			Name: "command",
			SubCommands: []*Command{
				{
					Name: "subcommand",
				},
				{
					Name: "subcommandTwo",
				},
			},
		}

		got, found := subject.findSubCommand("nope")
		if found {
			t.Error("want false, got true")
		}

		if got != nil {
			t.Errorf("want nil, got %v", got)
		}
	})
}

func TestCommand_run(t *testing.T) {
	t.Run("validate command runs", func(t *testing.T) {
		b := bytes.NewBufferString("")
		want := "command ran"
		subject := &Command{
			Name: "command",
			Out:  b,
			Err:  b,
			Run: func(command *Command) error {
				command.Out.Write([]byte("command ran"))
				return nil
			},
		}

		err := subject.run([]string{"-flag", "value"})
		if err != nil {
			t.Errorf("[err] want nil, got %s", err)
		}
		got, err := ioutil.ReadAll(b)
		if err != nil {
			t.Errorf("[err] want nil, got %s", err)
		}
		if string(got) != want {
			t.Errorf("want: %s, got %s", want, string(got))
		}
	})

	t.Run("validate command checks and runs", func(t *testing.T) {
		b := bytes.NewBufferString("")
		want := "check rancommand ran"
		subject := &Command{
			Name: "command",
			Out:  b,
			Err:  b,
			Check: func(command *Command) error {
				command.Out.Write([]byte("check ran"))
				return nil
			},
			Run: func(command *Command) error {
				command.Out.Write([]byte("command ran"))
				return nil
			},
		}

		err := subject.run([]string{"-flag", "value"})
		if err != nil {
			t.Errorf("[err] want nil, got %s", err)
		}
		got, err := ioutil.ReadAll(b)
		if err != nil {
			t.Errorf("[err] want nil, got %s", err)
		}
		if string(got) != want {
			t.Errorf("want: %s, got %s", want, string(got))
		}
	})

	t.Run("validate command returns error if check fails", func(t *testing.T) {
		b := bytes.NewBufferString("")
		want := "check ran"
		subject := &Command{
			Name: "command",
			Out:  b,
			Err:  b,
			Check: func(command *Command) error {
				command.Out.Write([]byte("check ran"))
				return errors.New("check failed")
			},
			Run: func(command *Command) error {
				command.Out.Write([]byte("command ran"))
				return nil
			},
		}

		err := subject.run([]string{"-flag", "value"})
		if err == nil {
			t.Error("wanted an error, got nil")
		}
		got, err := ioutil.ReadAll(b)
		if err != nil {
			t.Errorf("[err] want nil, got %s", err)
		}
		if string(got) != want {
			t.Errorf("want: %s, got %s", want, string(got))
		}
	})

	t.Run("validate command returns error if it is not runnable", func(t *testing.T) {
		want := ErrCommandNotRunnable
		subject := &Command{
			Name: "command",
		}

		got := subject.run([]string{"-flag", "value"})
		if got != want {
			t.Errorf("want %s, got %s", want, got)
		}
	})
}

func TestCommand_DefaultUsage(t *testing.T) {
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
			Flags: flag.NewFlagSet("command", flag.ContinueOnError),
			SubCommands: []*Command{
				{
					Name:        "subcommand",
					RunSyntax:   "test command subcommand [-flags...] [args...]",
					Description: "The test command subcommand.",
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
}
