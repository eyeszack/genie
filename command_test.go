package geenee

import (
	"bytes"
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"reflect"
	"testing"
)

type Mock struct {
	Name string
}

func (m *Mock) Check(command *Command) error {
	if m.Name == "" {
		return errors.New("need a name")
	}
	return nil
}

func (m *Mock) Run(command *Command) error {
	command.Out.Write([]byte(fmt.Sprintf("%s ran", m.Name)))
	return nil
}

func Test_NoopUsage(t *testing.T) {
	t.Run("validate noopusage does nothing", func(t *testing.T) {
		NoopUsage()
	})
}

func TestNoopWriter_Write(t *testing.T) {
	t.Run("validate noop does nothing", func(t *testing.T) {
		subject := NoopWriter{}
		got, err := subject.Write([]byte("hello everybody"))
		if err != nil {
			t.Errorf("want nil, got %s", err)
		}
		if got != 0 {
			t.Errorf("want 0, got %d", got)
		}
	})
}

func Test_NewCommand(t *testing.T) {
	t.Run("validate new - silence false", func(t *testing.T) {
		wantName := "tester"
		got := NewCommand(wantName, false)
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
		if got.Flags.ErrorHandling() != flag.ExitOnError {
			t.Errorf("want %d, got %d", flag.ExitOnError, got.Flags.ErrorHandling())
		}
	})

	t.Run("validate new - silence true", func(t *testing.T) {
		wantName := "tester"
		got := NewCommand(wantName, true)
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
		if got.Flags.ErrorHandling() != flag.ContinueOnError {
			t.Errorf("want %d, got %d", flag.ContinueOnError, got.Flags.ErrorHandling())
		}
		if reflect.TypeOf(got.Flags.Output()) != reflect.TypeOf(&NoopWriter{}) {
			t.Errorf("want %T, got %T", &NoopWriter{}, got.Flags.Output())
		}
	})
}

func TestCommand_AnchorPaths(t *testing.T) {
	t.Run("validate paths after anchoring a command", func(t *testing.T) {
		want := "test command command2"
		subject := &Command{
			Name: "test",
			SubCommands: []*Command{
				{
					Name: "command",
					SubCommands: []*Command{
						{
							Name: "command2",
						},
					},
				},
			},
		}

		subject.AnchorPaths()
		got := subject.SubCommands[0].SubCommands[0].Path()
		if got != want {
			t.Errorf("want %s, got %s", want, got)
		}
	})
}

func TestCommand_adjustPath(t *testing.T) {
	t.Run("validate paths after anchoring a command", func(t *testing.T) {
		want := "test command command2"
		subject := &Command{
			Name: "test",
			SubCommands: []*Command{
				{
					Name: "command",
					SubCommands: []*Command{
						{
							Name: "command2",
						},
					},
				},
			},
		}

		subject.SubCommands[0].adjustPath(subject.Name)
		got := subject.SubCommands[0].SubCommands[0].Path()
		if got != want {
			t.Errorf("want %s, got %s", want, got)
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

func TestCommand_flagIsSecret(t *testing.T) {
	t.Run("validate secret flag is secret", func(t *testing.T) {
		subject := Command{Name: "test"}
		subject.SecretFlag("shhh")

		if !subject.flagIsSecret("shhh") {
			t.Errorf("want true, got false")
		}

		if subject.flagIsSecret("nope") {
			t.Errorf("want false, got true")
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

	t.Run("validate command checks and runs - struct with methods", func(t *testing.T) {
		b := bytes.NewBufferString("")
		mock := Mock{}
		want := "mock ran"
		subject := &Command{
			Name:  "command",
			Out:   b,
			Err:   b,
			Check: mock.Check,
			Run:   mock.Run,
		}
		subject.Flags = flag.NewFlagSet("command", flag.ContinueOnError)
		subject.Flags.StringVar(&mock.Name, "name", "", "give me a name")

		err := subject.run([]string{"--name", "mock"})
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

	t.Run("validate command returns help error if it is not runnable, but asked for help", func(t *testing.T) {
		want := flag.ErrHelp
		subject := &Command{
			Name: "command",
		}

		got := subject.run([]string{"-flag", "-help"})
		if got != want {
			t.Errorf("want %s, got %s", want, got)
		}
	})

	t.Run("validate command run returns flag errors when silence flags is true", func(t *testing.T) {
		want := errors.New("flag provided but not defined: -flag")
		subject := NewCommand("silenced", true)
		subject.Run = func(command *Command) error {
			return nil
		}
		got := subject.run([]string{"-flag", "value"})
		_, ok := got.(Error)
		if !ok {
			t.Errorf("want geenee.Error, got %T", got)
		}
		if got.Error() != want.Error() {
			t.Errorf("want %s, got %s", want.Error(), got.Error())
		}
	})

	t.Run("validate command run returns nil when silence flags is true and -help is undefined but passed in", func(t *testing.T) {
		b := bytes.NewBufferString("")
		want := `
USAGE:
silenced

FLAGS:
--help        display help for command
`
		subject := NewCommand("silenced", true)
		subject.Out = b
		subject.Run = func(command *Command) error {
			return nil
		}
		err := subject.run([]string{"-help"})
		if err != nil {
			t.Errorf("want nil, got %s", err)
		}
		got, err := ioutil.ReadAll(b)
		if err != nil {
			t.Errorf("[err] want nil, got %s", err)
		}
		if string(got) != want {
			t.Errorf("want: %s, got %s", want, string(got))
		}
	})
}

func Test_DefaultCommandRunner(t *testing.T) {
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

		err := DefaultCommandRunner(subject, []string{"-flag", "value"})
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
}

func Test_askedForHelp(t *testing.T) {
	t.Run("validate --help is found", func(t *testing.T) {
		args := []string{"interface", "command", "subcommand", "-flag", "value", "--help", "-d"}
		hasFlag := askedForHelp(args)
		if !hasFlag {
			t.Errorf("want true, got %t", hasFlag)
		}
	})

	t.Run("validate -help is found", func(t *testing.T) {
		args := []string{"interface", "command", "subcommand", "-flag", "value", "-help", "-d"}
		hasFlag := askedForHelp(args)
		if !hasFlag {
			t.Errorf("want true, got %t", hasFlag)
		}
	})

	t.Run("validate help is not found", func(t *testing.T) {
		args := []string{"interface", "command", "subcommand", "-flag", "value", "-d"}
		hasFlag := askedForHelp(args)
		if hasFlag {
			t.Errorf("want false, got %t", hasFlag)
		}
	})
}
