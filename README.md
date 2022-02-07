# geenee
geenee is a simple, flexible, zero dependency module for building CLIs with Go.

## Getting Started

```go
package main

import (
	"fmt"
	"os"

	"github.com/eyeszack/geenee"
)

func main() {
	cmd := newHeyoCommand()
	err := geenee.DefaultCommandRunner(cmd, os.Args[1:])
	if err != nil {
		fmt.Printf(cmd.ShowUsage())
		os.Exit(1)
	}
}

func newHeyoCommand() *geenee.Command {
	test := ""
	cmd := geenee.NewCommand("heyo", true)
	cmd.Description = "A simple heyo."
	cmd.Check = func(command *geenee.Command) error {
		//you can validate flag/arg values here if desired,
		//or anything else you want to check on before running command
		return nil
	}
	cmd.Run = func(command *geenee.Command) error {
		fmt.Fprintf(command.Out, "I ran the heyo with [%s]\n", test)
		return nil
	}

	cmd.MergeFlagUsage = true
	cmd.Flags.StringVar(&test, "test", "", "the test flag")
	cmd.Flags.StringVar(&test, "t", "", "the test flag")

	return cmd
}
```