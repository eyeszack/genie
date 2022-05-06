package genie

import "testing"

func Test_GenerateBashCompletion(t *testing.T) {
	t.Run("validate bash completion", func(t *testing.T) {
		want := `#!/bin/bash
function _test () {
  COMPREPLY=($(test compreply "$COMP_LINE"));
};
complete -F _test test
`

		subject := &Lamp{Name: "test"}
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
		subject *Lamp
		want    string
	}{
		{
			name:    "empty",
			line:    "",
			subject: &Lamp{Name: "subject", RootCommand: &Command{Name: "subject", SubCommands: []*Command{&Command{Name: "heyo"}, &Command{Name: "mayo", Secret: true}}}},
			want:    "",
		},
		{
			name:    "no commands",
			line:    "subject",
			subject: &Lamp{Name: "subject"},
			want:    "",
		},
		{
			name:    "no commands at all",
			line:    "subject nope",
			subject: &Lamp{Name: "subject", RootCommand: &Command{Name: "subject"}},
			want:    "",
		},
		{
			name:    "command not found",
			line:    "subject nope",
			subject: &Lamp{Name: "subject", RootCommand: &Command{Name: "subject", SubCommands: []*Command{&Command{Name: "heyo"}, &Command{Name: "mayo", Secret: true}}}},
			want:    "",
		},
		{
			name:    "commands on root returned",
			line:    "subject",
			subject: &Lamp{Name: "subject", RootCommand: &Command{Name: "subject", SubCommands: []*Command{&Command{Name: "heyo"}, &Command{Name: "playo"}, &Command{Name: "mayo", Secret: true}}}},
			want:    "heyo playo",
		},
		{
			name:    "commands on root returned - trailing space",
			line:    "subject ",
			subject: &Lamp{Name: "subject", RootCommand: &Command{Name: "subject", SubCommands: []*Command{&Command{Name: "heyo"}, &Command{Name: "playo"}, &Command{Name: "mayo", Secret: true}}}},
			want:    "heyo playo",
		},
		{
			name:    "commands on root returned - trailing tab",
			line:    "subject\t",
			subject: &Lamp{Name: "subject", RootCommand: &Command{Name: "subject", SubCommands: []*Command{&Command{Name: "heyo"}, &Command{Name: "playo"}, &Command{Name: "mayo", Secret: true}}}},
			want:    "heyo playo",
		},
		{
			name:    "commands on root returned - partial",
			line:    "subject h",
			subject: &Lamp{Name: "subject", RootCommand: &Command{Name: "subject", SubCommands: []*Command{&Command{Name: "heyo"}, &Command{Name: "playo"}, &Command{Name: "mayo", Secret: true}}}},
			want:    "heyo",
		},
		{
			name:    "subcommands returned - full path no trailing space",
			line:    "subject heyo",
			subject: &Lamp{Name: "subject", RootCommand: &Command{Name: "subject", SubCommands: []*Command{&Command{Name: "heyo", SubCommands: []*Command{&Command{Name: "cool"}, &Command{Name: "bro"}}}, &Command{Name: "playo", SubCommands: []*Command{&Command{Name: "dang"}, &Command{Name: "it"}}}}}},
			want:    "",
		},
		{
			name:    "subcommands returned - full path trailing space",
			line:    "subject heyo ",
			subject: &Lamp{Name: "subject", RootCommand: &Command{Name: "subject", SubCommands: []*Command{&Command{Name: "heyo", SubCommands: []*Command{&Command{Name: "cool"}, &Command{Name: "bro"}}}, &Command{Name: "playo", SubCommands: []*Command{&Command{Name: "dang"}, &Command{Name: "it"}}}}}},
			want:    "cool bro",
		},
		{
			name:    "subcommands returned - full path trailing tab",
			line:    "subject heyo\t",
			subject: &Lamp{Name: "subject", RootCommand: &Command{Name: "subject", SubCommands: []*Command{&Command{Name: "heyo", SubCommands: []*Command{&Command{Name: "cool"}, &Command{Name: "bro"}, &Command{Name: "man", Secret: true}}}, &Command{Name: "playo", SubCommands: []*Command{&Command{Name: "dang"}, &Command{Name: "it"}}}}}},
			want:    "cool bro",
		},
		{
			name:    "subcommands returned - full path space then tab",
			line:    "subject heyo \t",
			subject: &Lamp{Name: "subject", RootCommand: &Command{Name: "subject", SubCommands: []*Command{&Command{Name: "heyo", SubCommands: []*Command{&Command{Name: "cool"}, &Command{Name: "bro"}}}, &Command{Name: "playo", SubCommands: []*Command{&Command{Name: "dang"}, &Command{Name: "it"}}}}}},
			want:    "cool bro",
		},
		{
			name:    "subcommands returned - partial path",
			line:    "subject heyo c",
			subject: &Lamp{Name: "subject", RootCommand: &Command{Name: "subject", SubCommands: []*Command{&Command{Name: "heyo", SubCommands: []*Command{&Command{Name: "cool"}, &Command{Name: "bro"}, &Command{Name: "crazy"}}}, &Command{Name: "playo", SubCommands: []*Command{&Command{Name: "dang"}, &Command{Name: "it"}}}}}},
			want:    "cool crazy",
		},
		{
			name:    "subcommands returned - partial path not found",
			line:    "subject heyo z",
			subject: &Lamp{Name: "subject", RootCommand: &Command{Name: "subject", SubCommands: []*Command{&Command{Name: "heyo", SubCommands: []*Command{&Command{Name: "cool"}, &Command{Name: "bro"}, &Command{Name: "crazy"}}}, &Command{Name: "playo", SubCommands: []*Command{&Command{Name: "dang"}, &Command{Name: "it"}}}}}},
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
