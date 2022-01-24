package geenee

import "testing"

func Test_GenerateBashCompletion(t *testing.T) {
	t.Run("validate bash completion", func(t *testing.T) {
		want := `#!/bin/bash
function _test () {
  COMPREPLY=($(test compreply $COMP_LINE));
};
complete -F _test test
`

		subject := &CommandInterface{Name: "test"}
		got := GenerateBashCompletion(subject)
		if got != want {
			t.Errorf("want %s, got %s", want, got)
		}
	})
}

func TestCommandInterface_CompletionReply(t *testing.T) {
	testCases := []struct {
		name    string
		line    string
		subject *CommandInterface
		want    string
	}{
		{
			name:    "empty",
			line:    "",
			subject: &CommandInterface{Name: "subject", RootCommand: &Command{Name: "subject", SubCommands: []*Command{&Command{Name: "heyo"}}}},
			want:    "",
		},
		{
			name:    "no commands",
			line:    "subject",
			subject: &CommandInterface{Name: "subject"},
			want:    "",
		},
		{
			name:    "no commands at all",
			line:    "subject nope",
			subject: &CommandInterface{Name: "subject", RootCommand: &Command{Name: "subject"}},
			want:    "",
		},
		{
			name:    "command not found",
			line:    "subject nope",
			subject: &CommandInterface{Name: "subject", RootCommand: &Command{Name: "subject", SubCommands: []*Command{&Command{Name: "heyo"}}}},
			want:    "",
		},
		{
			name:    "commands on root returned",
			line:    "subject",
			subject: &CommandInterface{Name: "subject", RootCommand: &Command{Name: "subject", SubCommands: []*Command{&Command{Name: "heyo"}, &Command{Name: "playo"}}}},
			want:    "heyo playo",
		},
		{
			name:    "commands on root returned - trailing space",
			line:    "subject ",
			subject: &CommandInterface{Name: "subject", RootCommand: &Command{Name: "subject", SubCommands: []*Command{&Command{Name: "heyo"}, &Command{Name: "playo"}}}},
			want:    "heyo playo",
		},
		{
			name:    "commands on root returned - partial",
			line:    "subject h",
			subject: &CommandInterface{Name: "subject", RootCommand: &Command{Name: "subject", SubCommands: []*Command{&Command{Name: "heyo"}, &Command{Name: "playo"}}}},
			want:    "heyo",
		},
		{
			name:    "subcommands returned - full path no trailing space",
			line:    "subject heyo",
			subject: &CommandInterface{Name: "subject", RootCommand: &Command{Name: "subject", SubCommands: []*Command{&Command{Name: "heyo", SubCommands: []*Command{&Command{Name: "cool"}, &Command{Name: "bro"}}}, &Command{Name: "playo", SubCommands: []*Command{&Command{Name: "dang"}, &Command{Name: "it"}}}}}},
			want:    "",
		},
		{
			name:    "subcommands returned - full path trailing space",
			line:    "subject heyo ",
			subject: &CommandInterface{Name: "subject", RootCommand: &Command{Name: "subject", SubCommands: []*Command{&Command{Name: "heyo", SubCommands: []*Command{&Command{Name: "cool"}, &Command{Name: "bro"}}}, &Command{Name: "playo", SubCommands: []*Command{&Command{Name: "dang"}, &Command{Name: "it"}}}}}},
			want:    "cool bro",
		},
		{
			name:    "subcommands returned - partial path",
			line:    "subject heyo c",
			subject: &CommandInterface{Name: "subject", RootCommand: &Command{Name: "subject", SubCommands: []*Command{&Command{Name: "heyo", SubCommands: []*Command{&Command{Name: "cool"}, &Command{Name: "bro"}, &Command{Name: "crazy"}}}, &Command{Name: "playo", SubCommands: []*Command{&Command{Name: "dang"}, &Command{Name: "it"}}}}}},
			want:    "cool crazy",
		},
		{
			name:    "subcommands returned - partial path not found",
			line:    "subject heyo z",
			subject: &CommandInterface{Name: "subject", RootCommand: &Command{Name: "subject", SubCommands: []*Command{&Command{Name: "heyo", SubCommands: []*Command{&Command{Name: "cool"}, &Command{Name: "bro"}, &Command{Name: "crazy"}}}, &Command{Name: "playo", SubCommands: []*Command{&Command{Name: "dang"}, &Command{Name: "it"}}}}}},
			want:    "",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got := tc.subject.CompletionReply(tc.line)
			if got != tc.want {
				t.Errorf("want %s, got %s", tc.want, got)
			}
		})
	}
}
