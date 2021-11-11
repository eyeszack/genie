package geenee

import (
	"bytes"
	"io/ioutil"
	"os"
	"reflect"
	"testing"
)

func TestInterfaceError_Error(t *testing.T) {
	t.Run("validate error string", func(t *testing.T) {
		want := "heyo error"
		got := GeeneeError("heyo error")
		if want != got.Error() {
			t.Errorf("want %s, got %s", want, got)
		}
	})
}

func Test_NewInterface(t *testing.T) {
	t.Run("validate new", func(t *testing.T) {
		wantName := "tester"
		wantVersion := "0.0.0"
		got := NewInterface(wantName, wantVersion, false)
		if got == nil {
			t.Fatal("want interface, got nil")
		}

		if got.Name != wantName {
			t.Errorf("want %s, got %s", wantName, got.Name)
		}
		if got.Version != wantVersion {
			t.Errorf("want %s, got %s", wantVersion, got.Version)
		}
		if got.Out != os.Stdout {
			t.Errorf("want %v, got %v", os.Stdout, got.Out)
		}
		if got.Err != os.Stderr {
			t.Errorf("want %v, got %v", os.Stderr, got.Err)
		}

		if got.RootCommand == nil {
			t.Fatal("want root command, got nil")
		}

		if got.RootCommand.Name != wantName {
			t.Errorf("want %s, got %s", wantName, got.RootCommand.Name)
		}
		if got.RootCommand.Flags.Name() != wantName {
			t.Errorf("want %s, got %s", wantName, got.RootCommand.Flags.Name())
		}
		if got.RootCommand.Out != os.Stdout {
			t.Errorf("want %v, got %v", os.Stdout, got.RootCommand.Out)
		}
		if got.RootCommand.Err != os.Stderr {
			t.Errorf("want %v, got %v", os.Stderr, got.RootCommand.Err)
		}
		if got.RootCommand.Usage == nil {
			t.Error("want DefaultCommandUsageFunc, got nil")
		}
	})
}

func TestInterface_SetWriters(t *testing.T) {
	t.Run("validate new", func(t *testing.T) {
		b := bytes.NewBufferString("")
		got := NewInterface("tester", "0.0.0", false)
		got.RootCommand.SubCommands = []*Command{
			NewCommand("heyo", false),
		}

		if got.Out != os.Stdout {
			t.Errorf("want %v, got %v", os.Stdout, got.Out)
		}
		if got.Err != os.Stderr {
			t.Errorf("want %v, got %v", os.Stderr, got.Err)
		}

		if got.RootCommand.Out != os.Stdout {
			t.Errorf("want %v, got %v", os.Stdout, got.RootCommand.Out)
		}
		if got.RootCommand.Err != os.Stderr {
			t.Errorf("want %v, got %v", os.Stderr, got.RootCommand.Err)
		}

		if got.RootCommand.Flags.Output() != os.Stderr {
			t.Errorf("want %v, got %v", os.Stderr, got.RootCommand.Flags.Output())
		}

		if got.RootCommand.SubCommands[0].Out != os.Stdout {
			t.Errorf("want %v, got %v", os.Stdout, got.RootCommand.SubCommands[0].Out)
		}
		if got.RootCommand.SubCommands[0].Err != os.Stderr {
			t.Errorf("want %v, got %v", os.Stderr, got.RootCommand.SubCommands[0].Err)
		}

		if got.RootCommand.SubCommands[0].Flags.Output() != os.Stderr {
			t.Errorf("want %v, got %v", os.Stderr, got.RootCommand.SubCommands[0].Flags.Output())
		}

		got.SetWriters(b, b)

		if got.Out != b {
			t.Errorf("want %v, got %v", b, got.Out)
		}
		if got.Err != b {
			t.Errorf("want %v, got %v", b, got.Err)
		}

		if got.RootCommand.Out != b {
			t.Errorf("want %v, got %v", b, got.RootCommand.Out)
		}
		if got.RootCommand.Err != b {
			t.Errorf("want %v, got %v", b, got.RootCommand.Err)
		}

		if got.RootCommand.Flags.Output() != b {
			t.Errorf("want %v, got %v", b, got.RootCommand.Flags.Output())
		}

		if got.RootCommand.SubCommands[0].Out != b {
			t.Errorf("want %v, got %v", b, got.RootCommand.SubCommands[0].Out)
		}
		if got.RootCommand.SubCommands[0].Err != b {
			t.Errorf("want %v, got %v", b, got.RootCommand.SubCommands[0].Err)
		}

		if got.RootCommand.SubCommands[0].Flags.Output() != b {
			t.Errorf("want %v, got %v", b, got.RootCommand.SubCommands[0].Flags.Output())
		}

		got.SetWriters(nil, nil)
		if got.Out != b {
			t.Errorf("want %v, got %v", b, got.Out)
		}
		if got.Err != b {
			t.Errorf("want %v, got %v", b, got.Err)
		}

		if got.RootCommand.Out != b {
			t.Errorf("want %v, got %v", b, got.RootCommand.Out)
		}
		if got.RootCommand.Err != b {
			t.Errorf("want %v, got %v", b, got.RootCommand.Err)
		}

		if got.RootCommand.Flags.Output() != b {
			t.Errorf("want %v, got %v", b, got.RootCommand.Flags.Output())
		}

		if got.RootCommand.SubCommands[0].Out != b {
			t.Errorf("want %v, got %v", b, got.RootCommand.SubCommands[0].Out)
		}
		if got.RootCommand.SubCommands[0].Err != b {
			t.Errorf("want %v, got %v", b, got.RootCommand.SubCommands[0].Err)
		}

		if got.RootCommand.SubCommands[0].Flags.Output() != b {
			t.Errorf("want %v, got %v", b, got.RootCommand.SubCommands[0].Flags.Output())
		}
	})
}

func TestInterface_Exec(t *testing.T) {
	t.Run("validate interface root command runs - no flags/args", func(t *testing.T) {
		b := bytes.NewBufferString("")
		want := "root command ran"
		subject := &Interface{
			Name: "test",
			RootCommand: &Command{
				Name: "test",
				Out:  b,
				Err:  b,
				Run: func(command *Command) error {
					command.Out.Write([]byte("root command ran"))
					return nil
				},
			},
			Out: b,
			Err: b,
		}
		err := subject.Exec([]string{"test"})
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

	t.Run("validate interface root command runs - no flags with args", func(t *testing.T) {
		b := bytes.NewBufferString("")
		want := "root command ran"
		subject := &Interface{
			Name: "test",
			RootCommand: &Command{
				Name: "test",
				Out:  b,
				Err:  b,
				Run: func(command *Command) error {
					command.Out.Write([]byte("root command ran"))
					return nil
				},
			},
			Out: b,
			Err: b,
		}
		err := subject.Exec([]string{"test", "notaflag", "notasubcommand"})
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

	t.Run("validate interface root command runs - no flags", func(t *testing.T) {
		b := bytes.NewBufferString("")
		want := "root command ran"
		subject := &Interface{
			Name: "test",
			RootCommand: &Command{
				Name: "test",
				Out:  b,
				Err:  b,
				Run: func(command *Command) error {
					command.Out.Write([]byte("root command ran"))
					return nil
				},
			},
			Out: b,
			Err: b,
		}
		err := subject.Exec([]string{"test"})
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

	t.Run("validate interface root command runs - flags", func(t *testing.T) {
		b := bytes.NewBufferString("")
		want := "root command ran"
		subject := &Interface{
			Name: "test",
			RootCommand: &Command{
				Name: "test",
				Out:  b,
				Err:  b,
				Run: func(command *Command) error {
					command.Out.Write([]byte("root command ran"))
					return nil
				},
			},
			Out: b,
			Err: b,
		}
		err := subject.Exec([]string{"test", "-flag", "value"})
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

	t.Run("validate command runs - no flags", func(t *testing.T) {
		b := bytes.NewBufferString("")
		want := "command ran"
		subject := &Interface{
			Name: "test",
			RootCommand: &Command{
				Name: "test",
				Out:  b,
				Err:  b,
				Run: func(command *Command) error {
					command.Out.Write([]byte("root command ran"))
					return nil
				},
				SubCommands: []*Command{
					{
						Name: "command",
						Out:  b,
						Err:  b,
						Run: func(command *Command) error {
							command.Out.Write([]byte("command ran"))
							return nil
						},
					},
				},
			},
			Out: b,
			Err: b,
		}
		err := subject.Exec([]string{"test", "command"})
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

	t.Run("validate command runs - flags", func(t *testing.T) {
		b := bytes.NewBufferString("")
		want := "command ran"
		subject := &Interface{
			Name: "test",
			RootCommand: &Command{
				Name: "test",
				Out:  b,
				Err:  b,
				Run: func(command *Command) error {
					command.Out.Write([]byte("root command ran"))
					return nil
				},
				SubCommands: []*Command{
					{
						Name: "command",
						Out:  b,
						Err:  b,
						Run: func(command *Command) error {
							command.Out.Write([]byte("command ran"))
							return nil
						},
					},
				},
			},
			Out: b,
			Err: b,
		}
		err := subject.Exec([]string{"test", "command", "-flag", "value"})
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

	t.Run("validate subcommand runs - no flags", func(t *testing.T) {
		b := bytes.NewBufferString("")
		want := "subcommand ran"
		subject := &Interface{
			Name: "test",
			RootCommand: &Command{
				Name: "test",
				Out:  b,
				Err:  b,
				Run: func(command *Command) error {
					command.Out.Write([]byte("root command ran"))
					return nil
				},
				SubCommands: []*Command{
					{
						Name: "command",
						Out:  b,
						Err:  b,
						Run: func(command *Command) error {
							command.Out.Write([]byte("command ran"))
							return nil
						},
						SubCommands: []*Command{
							{
								Name: "subcommand",
								Out:  b,
								Err:  b,
								Run: func(command *Command) error {
									command.Out.Write([]byte("subcommand ran"))
									return nil
								},
							},
						},
					},
				},
			},
			Out: b,
			Err: b,
		}
		err := subject.Exec([]string{"test", "command", "subcommand"})
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

	t.Run("validate subcommand runs - flags", func(t *testing.T) {
		b := bytes.NewBufferString("")
		want := "subcommand ran"
		subject := &Interface{
			Name: "test",
			RootCommand: &Command{
				Name: "test",
				Out:  b,
				Err:  b,
				Run: func(command *Command) error {
					command.Out.Write([]byte("root command ran"))
					return nil
				},
				SubCommands: []*Command{
					{
						Name: "command",
						Out:  b,
						Err:  b,
						Run: func(command *Command) error {
							command.Out.Write([]byte("command ran"))
							return nil
						},
						SubCommands: []*Command{
							{
								Name: "subcommand",
								Out:  b,
								Err:  b,
								Run: func(command *Command) error {
									command.Out.Write([]byte("subcommand ran"))
									return nil
								},
							},
						},
					},
				},
			},
			Out: b,
			Err: b,
		}
		err := subject.Exec([]string{"test", "command", "subcommand", "-flag", "value"})
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

func TestInterface_Exec_error(t *testing.T) {
	t.Run("validate interface returns error if no args provided", func(t *testing.T) {
		want := ErrNoArgs
		subject := &Interface{
			Name: "test",
			RootCommand: &Command{
				Name: "test",
			},
		}
		got := subject.Exec(nil)
		if got != want {
			t.Errorf("want %s, got %s", want, got)
		}
	})

	t.Run("validate interface returns error if noop state", func(t *testing.T) {
		want := ErrNoOp
		subject := &Interface{
			Name: "test",
		}
		got := subject.Exec([]string{"test"})
		if got != want {
			t.Errorf("want %s, got %s", want, got)
		}
	})

	t.Run("validate interface returns error if different interface called", func(t *testing.T) {
		want := ErrInvalidInterfaceName
		subject := &Interface{
			Name: "test",
			RootCommand: &Command{
				Name: "test",
			},
		}
		got := subject.Exec([]string{"nottest"})
		if got != want {
			t.Errorf("want %s, got %s", want, got)
		}
	})

	t.Run("validate interface returns error if invalid command called (noop) - no flag", func(t *testing.T) {
		want := ErrNoOp
		subject := &Interface{
			Name: "test",
		}
		got := subject.Exec([]string{"test", "nope"})
		if got != want {
			t.Errorf("want %s, got %s", want, got)
		}
	})

	t.Run("validate interface returns error if invalid command called and root not runnable - no flag", func(t *testing.T) {
		want := ErrCommandNotRunnable
		subject := &Interface{
			Name: "test",
			RootCommand: &Command{
				Name: "test",
				SubCommands: []*Command{
					{
						Name: "command",
						Run: func(command *Command) error {
							t.Error("unexpected command run")
							return nil
						},
					},
				},
			},
		}
		got := subject.Exec([]string{"test", "nope"})
		if got != want {
			t.Errorf("want %s, got %s", want, got)
		}
	})

	t.Run("validate interface returns error if invalid command called", func(t *testing.T) {
		want := ErrCommandNotFound
		subject := &Interface{
			Name: "test",
			RootCommand: &Command{
				Name: "test",
				SubCommands: []*Command{
					{
						Name: "command",
						Run: func(command *Command) error {
							t.Error("unexpected command run")
							return nil
						},
					},
				},
			},
		}
		got := subject.Exec([]string{"test", "nope", "-flag", "value"})
		if got != want {
			t.Errorf("want %s, got %s", want, got)
		}
	})

	t.Run("validate interface returns error if invalid subcommand called", func(t *testing.T) {
		want := ErrCommandNotFound
		subject := &Interface{
			Name: "test",
			RootCommand: &Command{
				Name: "test",
				SubCommands: []*Command{
					{
						Name: "command",
						Run: func(command *Command) error {
							t.Error("unexpected command run")
							return nil
						},
						SubCommands: []*Command{
							{
								Name: "subcommand",
								Run: func(command *Command) error {
									t.Error("unexpected subcommand run")
									return nil
								},
							},
						},
					},
				},
			},
		}
		got := subject.Exec([]string{"test", "command", "nope", "-flag", "value"})
		if got != want {
			t.Errorf("want %s, got %s", want, got)
		}
	})

	t.Run("validate interface returns error if valid subcommand called on an invalid command", func(t *testing.T) {
		want := ErrCommandNotFound
		subject := &Interface{
			Name: "test",
			RootCommand: &Command{
				Name: "test",
				SubCommands: []*Command{
					{
						Name: "command",
						Run: func(command *Command) error {
							t.Error("unexpected command run")
							return nil
						},
						SubCommands: []*Command{
							{
								Name: "subcommand",
								Run: func(command *Command) error {
									t.Error("unexpected subcommand run")
									return nil
								},
							},
						},
					},
				},
			},
		}
		got := subject.Exec([]string{"test", "nope", "subcommand", "-flag", "value"})
		if got != want {
			t.Errorf("want %s, got %s", want, got)
		}
	})

	t.Run("validate interface returns error if command depth invalid", func(t *testing.T) {
		want := ErrCommandDepthInvalid
		subject := &Interface{
			Name: "test",
			RootCommand: &Command{
				Name: "test",
				SubCommands: []*Command{
					{
						Name: "command",
						Run: func(command *Command) error {
							t.Error("unexpected command run")
							return nil
						},
						SubCommands: []*Command{
							{
								Name: "subcommand",
								Run: func(command *Command) error {
									t.Error("unexpected subcommand run")
									return nil
								},
							},
						},
					},
				},
			},
		}
		got := subject.Exec([]string{"test", "command", "subcommand", "subsubcommand", "-flag", "value"})
		if got != want {
			t.Errorf("want %s, got %s", want, got)
		}
	})
}

func TestInterface_searchPathForCommand(t *testing.T) {
	t.Run("validate command is found", func(t *testing.T) {
		want := &Command{Name: "command2"}
		subject := &Interface{
			Name: "test",
			RootCommand: &Command{
				Name: "test",
				SubCommands: []*Command{
					{
						Name: "command",
					},
					{
						Name: "command2",
					},
					{
						Name: "command3",
					},
				},
			},
		}
		got, found := subject.searchPathForCommand([]string{"command2"})
		if !found {
			t.Errorf("want true, got %t", found)
		}
		if !reflect.DeepEqual(got, want) {
			t.Errorf("want %v, got %v", want, got)
		}
	})

	t.Run("validate command is found - deep", func(t *testing.T) {
		want := &Command{Name: "dang", Description: "It's me you're looking for."}
		subject := &Interface{
			Name: "test",
			RootCommand: &Command{
				Name: "test",
				SubCommands: []*Command{
					{
						Name: "command",
						SubCommands: []*Command{
							{
								Name: "subcommand",
							},
						},
					},
					{
						Name: "command2",
						SubCommands: []*Command{
							{
								Name: "subcommand2",
								SubCommands: []*Command{
									{
										Name: "subsubcommand3",
										SubCommands: []*Command{
											{
												Name:        "dang",
												Description: "It's me you're looking for.",
											},
										},
									},
								},
							},
						},
					},
					{
						Name: "command3",
					},
				},
			},
		}
		got, found := subject.searchPathForCommand([]string{"command2", "subcommand2", "subsubcommand3", "dang"})
		if !found {
			t.Errorf("want true, got %t", found)
		}
		if !reflect.DeepEqual(got, want) {
			t.Errorf("want %v, got %v", want, got)
		}
	})

	t.Run("validate command is found with alias", func(t *testing.T) {
		want := &Command{Name: "command2", Aliases: []string{"heyo"}}
		subject := &Interface{
			Name: "test",
			RootCommand: &Command{
				Name: "test",
				SubCommands: []*Command{
					{
						Name:    "command",
						Aliases: []string{"nope"},
					},
					{
						Name:    "command2",
						Aliases: []string{"heyo"},
					},
					{
						Name: "command3",
					},
				},
			},
		}
		got, found := subject.searchPathForCommand([]string{"heyo"})
		if !found {
			t.Errorf("want true, got %t", found)
		}
		if !reflect.DeepEqual(got, want) {
			t.Errorf("want %v, got %v", want, got)
		}
	})

	t.Run("validate command is not found", func(t *testing.T) {
		subject := &Interface{
			Name: "test",
			RootCommand: &Command{
				Name: "test",
				SubCommands: []*Command{
					{
						Name: "command",
					},
					{
						Name: "command2",
					},
					{
						Name: "command3",
					},
				},
			},
		}
		got, found := subject.searchPathForCommand([]string{"nope"})
		if found {
			t.Errorf("want false, got %t", found)
		}
		if got != nil {
			t.Errorf("want nil, got %v", got)
		}
	})
}

func TestInterface_hasFlags(t *testing.T) {
	t.Run("validate flag is found in correct position", func(t *testing.T) {
		want := 3
		args := []string{"interface", "command", "subcommand", "-flag", "value", "heyo"}
		subject := &Interface{}
		got, hasFlag := subject.hasFlags(args)
		if !hasFlag {
			t.Errorf("want true, got %t", hasFlag)
		}
		if got != want {
			t.Errorf("want %d, got %d", want, got)
		}
	})

	t.Run("validate flag is found in correct position (--)", func(t *testing.T) {
		want := 2
		args := []string{"interface", "command", "--flag", "value", "heyo"}
		subject := &Interface{}
		got, hasFlag := subject.hasFlags(args)
		if !hasFlag {
			t.Errorf("want true, got %t", hasFlag)
		}
		if got != want {
			t.Errorf("want %d, got %d", want, got)
		}
	})

	t.Run("validate flag is not found", func(t *testing.T) {
		want := -1
		args := []string{"interface", "command", "subcommand", "notaflag", "value"}
		subject := &Interface{}
		got, hasFlag := subject.hasFlags(args)
		if hasFlag {
			t.Errorf("want false, got %t", hasFlag)
		}
		if got != want {
			t.Errorf("want %d, got %d", want, got)
		}
	})
}
