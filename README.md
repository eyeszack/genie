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
	lamp := genie.NewLamp("magic", "1.0.0", true)
	lamp.RootCommand.SubCommands = []*genie.Command{
		newWishCommand(),
	}
	
	if cmd, err := lamp.Execute(); err != nil {
		//you can show cmd usage on error, or do something else
		if cmd != nil {
			_, _ = fmt.Fprint(lamp.Err, cmd.ShowUsage())
			os.Exit(1)
		}
	}
}

func newWishCommand() *genie.Command {
	test := ""
	cmd := genie.NewCommand("wish", true)
	cmd.Description = "A simple wish."
	cmd.Check = func(command *genie.Command) error {
		//you can validate flag/arg values here if desired,
		//or anything else you want to check on before running command
		return nil
	}
	cmd.Run = func(command *genie.Command) error {
		_, err := fmt.Fprintf(command.Out, "I ran the wish with [%s]\n", test)
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

