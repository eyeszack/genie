package geenee

import (
	"bytes"
	"errors"
	"flag"
	"io/ioutil"
	"os"
	"reflect"
	"testing"
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
		if got.Path != wantName {
			t.Errorf("want %s, got %s", wantName, got.Path)
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

func TestCommand_AddSubCommand(t *testing.T) {
	t.Run("validate add subcommand", func(t *testing.T) {
		wantPath := "root command"
		wantLen := 1
		rootCommand := &Command{
			Name: "root",
			Path: "root",
		}
		subject := &Command{
			Name: "command",
			Path: "command",
		}
		rootCommand.AddSubCommand(subject)
		if len(rootCommand.SubCommands) != wantLen {
			t.Errorf("want %d, got %d", wantLen, len(rootCommand.SubCommands))
		}
		if subject.Path != wantPath {
			t.Errorf("want %s, got %s", wantPath, subject.Path)
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

	t.Run("validate subcommand is found with alias", func(t *testing.T) {
		want := &Command{Name: "subcommand", Aliases: []string{"me"}}
		subject := &Command{
			Name: "command",
			SubCommands: []*Command{
				{
					Name:    "subcommand",
					Aliases: []string{"me"},
				},
				{
					Name:    "subcommandTwo",
					Aliases: []string{"nope"},
				},
			},
		}

		got, found := subject.findSubCommand("me")
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
