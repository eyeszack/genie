package geenee

import (
	"fmt"
	"strings"
)

func (ci *CommandInterface) CompletionReply(line string) string {
	reply := ""
	if ci.RootCommand == nil || len(ci.RootCommand.SubCommands) == 0 {
		return reply
	}

	path := strings.Split(line, " ")
	if len(path) == 1 {
		if path[0] == "" {
			return reply
		}
		//this means we should pass the subcommands on the root, if available we won't check if path[0] actually
		//matched the root/interface name in the off chance that someone registered completion under alias or some such
		for i, sc := range ci.RootCommand.SubCommands {
			if i == 0 {
				reply += sc.Name
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

	cmd, found, pos := ci.searchPathForCommand(path[1:], true)
	if !found {
		//let's double check in case the completion request is on root
		if len(path) == 2 {
			i := 0
			for _, sc := range ci.RootCommand.SubCommands {
				if strings.HasPrefix(sc.Name, path[1]) {
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

func GenerateBashCompletion(cli *CommandInterface) string {
	complete := `#!/bin/bash
function _%s () {
  COMPREPLY=($(%s compreply $COMP_LINE));
};
complete -F _%s %s
`

	return fmt.Sprintf(complete, cli.Name, cli.Name, cli.Name, cli.Name)
}
