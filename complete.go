package genie

import (
	"fmt"
	"strings"
)

//CompletionReply provides very basic completion support for CLIs using geenee. Currently, only subcommand completion is
//supported. Flag and argument completion is not supported.
func (l *Lamp) CompletionReply(line string) string {
	reply := ""
	if l.RootCommand == nil || len(l.RootCommand.SubCommands) == 0 {
		return reply
	}

	path := strings.Split(line, " ")
	if len(path) == 1 {
		if path[0] == "" {
			return reply
		}
		//this means we should pass the subcommands on the root, if available we won't check if path[0] actually
		//matched the root/interface name in the off chance that someone registered completion under alias or some such
		i := 0
		for _, sc := range l.RootCommand.SubCommands {
			if sc.Secret { //don't show secret commands
				continue
			}

			if i == 0 {
				reply += sc.Name
				i++
				continue
			}
			reply = fmt.Sprintf("%s %s", reply, sc.Name)
		}
		return reply
	}

	if strings.HasSuffix(path[len(path)-1], "\t") {
		path[len(path)-1] = strings.TrimSuffix(path[len(path)-1], "\t")
		path = append(path, "")
	}

	cmd, found, pos := l.searchPathForCommand(path[1:], true)
	if !found {
		//let's double check in case the completion request is on root
		if len(path) == 2 {
			i := 0
			for _, sc := range l.RootCommand.SubCommands {
				if strings.HasPrefix(sc.Name, path[1]) {
					if sc.Secret { //don't show secret commands
						continue
					}

					if i == 0 {
						reply += sc.Name
						i++
						continue
					}
					reply = fmt.Sprintf("%s %s", reply, sc.Name)
				}
			}
		}
		return reply
	}
	//we add 1 to the find position to account for the first element we exclude in the searching
	pos += 1

	if len(path)-1 == pos { //this means the completion request is just for the next set of subcommands
		return reply
	} else { //we have values after the command that was found, let's take next path part and narrow down subcommands
		i := 0
		for _, sc := range cmd.SubCommands {
			if strings.HasPrefix(sc.Name, path[pos+1]) {
				if sc.Secret { //don't show secret commands
					continue
				}

				if i == 0 {
					reply += sc.Name
					i++
					continue
				}
				reply = fmt.Sprintf("%s %s", reply, sc.Name)
			}
		}
	}

	return reply
}

//GenerateBashCompletion will return a bash script that can be sourced to provide a hook into your completion logic.
//If you use this with your CLI you'll need to reply to the compreply with the appropriate values to show the user.
//You can use the simple completion support provided by the CompletionReply function or roll your own.
func GenerateBashCompletion(cli *Lamp) string {
	complete := `#!/bin/bash
function _%s () {
  COMPREPLY=($(%s compreply "$COMP_LINE"));
};
complete -F _%s %s
`

	return fmt.Sprintf(complete, cli.Name, cli.Name, cli.Name, cli.Name)
}
