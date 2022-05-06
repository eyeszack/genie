# genie
Genie is a simple, flexible, module for building CLIs in Go using only the standard library.

## Getting Started

```go
package main

import (
	"fmt"
	"os"

	"github.com/eyeszack/genie"
)

func main() {
	cmd := newHeyoCommand()
	err := genie.DefaultCommandRunner(cmd, os.Args[1:])
	if err != nil {
		fmt.Printf(cmd.ShowUsage())
		os.Exit(1)
	}
}

func newHeyoCommand() *genie.Command {
	test := ""
	cmd := genie.NewCommand("heyo", true)
	cmd.Description = "A simple heyo."
	cmd.Check = func(command *genie.Command) error {
		//you can validate flag/arg values here if desired,
		//or anything else you want to check on before running command
		return nil
	}
	cmd.Run = func(command *genie.Command) error {
		_, err := fmt.Fprintf(command.Out, "I ran the heyo with [%s]\n", test)
		if err != nil {
			return err
        }
		return nil
	}

	cmd.MergeFlagUsage = true
	cmd.Flags.StringVar(&test, "test", "", "the test flag")
	cmd.Flags.StringVar(&test, "t", "", "the test flag")

	return cmd
}
```

