package geenee

import (
	"flag"
	"fmt"
	"strings"
)

var DefaultCommandUsageFunc = func(command *Command) string {
	var builder strings.Builder

	if command.Description != "" {
		builder.WriteString(fmt.Sprintf("%s\n", command.Description))
	}

	builder.WriteString("\nUSAGE:\n")
	syntax := command.RunSyntax
	if syntax == "" {
		syntax = command.Name
	}
	builder.WriteString(fmt.Sprintf("%s\n", syntax))

	if command.ExtraInfo != "" {
		builder.WriteString(fmt.Sprintf("\n%s\n", command.ExtraInfo))
	}

	if command.Flags != nil {
		builder.WriteString("\nFLAGS:\n")
		command.Flags.VisitAll(func(f *flag.Flag) {
			defaultVal := ""
			if f.DefValue != "" {
				defaultVal = fmt.Sprintf(" (default %s)", f.DefValue)
			}

			typeOf := ""
			switch fmt.Sprintf("%T", f.Value) {
			case "*flag.boolValue":
				typeOf = "\tboolean"
			case "*flag.durationValue":
				typeOf = "\tduration"
			case "*flag.float64Value":
				typeOf = "\tfloat"
			case "*flag.intValue", "int64Value":
				typeOf = "\tint"
			case "*flag.stringValue":
				typeOf = "\tstring"
			case "*flag.uintValue", "uint64Value":
				typeOf = "\tuint"
			}

			dashes := "--"
			if len(f.Name) == 1 {
				dashes = "-"
			}
			builder.WriteString(fmt.Sprintf("  %s%s%s\n", dashes, f.Name, typeOf))
			builder.WriteString(fmt.Sprintf("\t\t%s%s\n", f.Usage, defaultVal))
		})
	}

	if len(command.SubCommands) > 0 {
		builder.WriteString("\nSUBCOMMANDS:\n")
		for _, subcommand := range command.SubCommands {
			builder.WriteString(fmt.Sprintf("%s\t%s\n", subcommand.Name, subcommand.Description))
		}
		//add this back if default --help logic is added
		//builder.WriteString("\nUse \"--help\" with any command or subcommand for more information.\n")
	}
	return builder.String()
}
