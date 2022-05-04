package geenee

import (
	"flag"
	"fmt"
	"sort"
	"strings"
	"text/tabwriter"
)

type UsageAwareFlagValue interface {
	Type() string
}

var DefaultCommandUsageFunc = func(command *Command) string {
	var builder strings.Builder

	if command.Description != "" {
		builder.WriteString(fmt.Sprintf("%s\n", command.Description))
	}

	builder.WriteString("\nUSAGE:\n")
	if command.path == "" {
		command.path = command.Name
	}
	syntax := fmt.Sprintf("%s %s", command.path, strings.ReplaceAll(command.RunSyntax, "{{path}}", command.path))
	builder.WriteString(fmt.Sprintf("%s\n", strings.Trim(syntax, " ")))

	if len(command.Aliases) > 0 {
		builder.WriteString("\nALIASES:\n")
		for _, a := range command.Aliases {
			builder.WriteString(fmt.Sprintf("%s\n", a))
		}
	}

	if command.ExtraInfo != "" {
		builder.WriteString(fmt.Sprintf("\n%s\n", command.ExtraInfo))
	}

	if command.MergeFlagUsage {
		builder.WriteString(mergeFlagsUsage(command))
	} else {
		builder.WriteString(flagsUsage(command))
	}

	if command.ArgInfo != "" {
		builder.WriteString("\nARGUMENTS:\n")
		builder.WriteString(fmt.Sprintf("%s\n", command.ArgInfo))
	}

	tabWriter := tabwriter.NewWriter(&builder, 0, 0, 4, ' ', tabwriter.DiscardEmptyColumns)
	wroteCommandHeader := false
	if len(command.SubCommands) > 0 {
		for _, subcommand := range command.SubCommands {
			if !subcommand.Secret {
				if !wroteCommandHeader {
					builder.WriteString("\nCOMMANDS:\n")
					wroteCommandHeader = true
				}
				_, _ = tabWriter.Write([]byte(fmt.Sprintf("%s\t%s\n", subcommand.Name, subcommand.Description)))
			}
		}
	}
	_ = tabWriter.Flush()
	if wroteCommandHeader {
		helpMsg := "\nUse \"--help\" with any command for more information.\n"
		if command.path != "" {
			helpMsg = fmt.Sprintf("\nUse \"%s <command> --help\" for more information.\n", command.path)
		}

		builder.WriteString(helpMsg)
	}

	return builder.String()
}

func mergeFlagsUsage(command *Command) string {
	usages := make(map[string][]string)

	var builder strings.Builder
	tabWriter := tabwriter.NewWriter(&builder, 0, 0, 4, ' ', tabwriter.DiscardEmptyColumns)
	if command.Flags != nil {
		command.Flags.VisitAll(func(f *flag.Flag) {
			if command.flagIsSecret(f.Name) {
				return
			}
			defaultVal := ""
			if f.DefValue != "" {
				defaultVal = fmt.Sprintf(" (default %s)", f.DefValue)
			}

			typeOf := ""
			switch fmt.Sprintf("%T", f.Value) {
			case "*flag.boolValue":
				typeOf = ""
			case "*flag.durationValue":
				typeOf = "duration"
			case "*flag.float64Value":
				typeOf = "float"
			case "*flag.intValue", "*flag.int64Value":
				typeOf = "int"
			case "*flag.stringValue":
				typeOf = "string"
			case "*flag.uintValue", "*flag.uint64Value":
				typeOf = "uint"
			default:
				u, ok := f.Value.(UsageAwareFlagValue)
				if ok {
					typeOf = u.Type()
				}
			}

			dashes := "--"
			if len(f.Name) == 1 {
				dashes = "-"
			}

			usage := fmt.Sprintf("%s%s", f.Usage, defaultVal)
			dashedFlag := fmt.Sprintf("%s%s", dashes, f.Name)

			existingFlags, exists := usages[usage]
			if exists {
				existingFlags = append([]string{dashedFlag}, existingFlags...)
				usages[usage] = existingFlags
			} else {
				temp := []string{dashedFlag, "\t" + typeOf}
				usages[usage] = temp
			}
		})
	}

	if command.root {
		//we add the version flag to root since we support it automatically
		usages["display version information"] = []string{"--version", "\t"}
	}
	usages["display help for command"] = []string{"--help", "\t"}

	//now we sort the flag usage to match flag package
	sorted := make([]string, len(usages))
	for k, v := range usages {
		sorted = append(sorted, fmt.Sprintf("%s\t%s\n", strings.TrimSuffix(strings.Join(v, " "), " "), k))
	}

	sort.Slice(sorted, func(i, j int) bool {
		return sorted[i] < sorted[j]
	})

	if len(sorted) > 0 {
		builder.WriteString("\nFLAGS:\n")
	}

	for _, s := range sorted {
		_, _ = tabWriter.Write([]byte(s))
	}

	_ = tabWriter.Flush()
	return builder.String()
}

func flagsUsage(command *Command) string {
	flagCount := 0
	var builder strings.Builder
	var usages []string
	tabWriter := tabwriter.NewWriter(&builder, 0, 0, 4, ' ', tabwriter.DiscardEmptyColumns)
	builder.WriteString("\nFLAGS:\n") //all commands have at least --help
	if command.Flags != nil {
		command.Flags.VisitAll(func(f *flag.Flag) {
			if command.flagIsSecret(f.Name) {
				return
			}

			flagCount++
			defaultVal := ""
			if f.DefValue != "" {
				defaultVal = fmt.Sprintf(" (default %s)", f.DefValue)
			}

			typeOf := ""
			switch fmt.Sprintf("%T", f.Value) {
			case "*flag.boolValue":
				typeOf = ""
			case "*flag.durationValue":
				typeOf = " duration"
			case "*flag.float64Value":
				typeOf = " float"
			case "*flag.intValue", "*flag.int64Value":
				typeOf = " int"
			case "*flag.stringValue":
				typeOf = " string"
			case "*flag.uintValue", "*flag.uint64Value":
				typeOf = " uint"
			default:
				u, ok := f.Value.(UsageAwareFlagValue)
				if ok {
					typeOf = fmt.Sprintf(" %s", u.Type())
				}
			}

			dashes := "--"
			if len(f.Name) == 1 {
				dashes = "-"
			}
			usages = append(usages, fmt.Sprintf("%s%s\t%s\t%s%s\n", dashes, f.Name, typeOf, f.Usage, defaultVal))
		})
	}

	if command.root {
		//we add the version flag to root since we support it automatically
		flagCount++
		usages = append(usages, fmt.Sprintf("%s\t%s\t%s\n", "--version", "", "display version information"))
	}

	flagCount++
	usages = append(usages, fmt.Sprintf("%s\t%s\t%s\n", "--help", "", "display help for command"))

	sort.Slice(usages, func(i, j int) bool {
		return usages[i] < usages[j]
	})

	for _, s := range usages {
		_, _ = tabWriter.Write([]byte(s))
	}

	_ = tabWriter.Flush()
	return builder.String()
}
