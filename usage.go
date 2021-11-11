package geenee

import (
	"flag"
	"fmt"
	"sort"
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
		if command.MergeFlagUsage {
			builder.WriteString(mergeFlagsUsage(command.Flags))
		} else {
			builder.WriteString(flagsUsage(command.Flags))
		}
	}

	if len(command.SubCommands) > 0 {
		wroteHeader := false
		for _, subcommand := range command.SubCommands {
			if !subcommand.Secret {
				if !wroteHeader {
					builder.WriteString("\nCOMMANDS:\n")
					wroteHeader = true
				}
				builder.WriteString(fmt.Sprintf("%s\t%s\n", subcommand.Name, subcommand.Description))
			}
		}
	}

	//flagset automatically handles -help or -h so it's basically like having a help command for free
	if command.Flags != nil {
		builder.WriteString("\nUse \"--help\" with any command for more information.\n")
	}

	return builder.String()
}

func mergeFlagsUsage(flags *flag.FlagSet) string {
	usages := make(map[string][]string)

	var builder strings.Builder
	if flags != nil {
		flags.VisitAll(func(f *flag.Flag) {
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
			case "*flag.intValue", "int64Value":
				typeOf = "int"
			case "*flag.stringValue":
				typeOf = "string"
			case "*flag.uintValue", "uint64Value":
				typeOf = "uint"
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
				temp := []string{dashedFlag, typeOf}
				usages[usage] = temp
			}
		})
	}

	//now we sort the flag usage to match flag package
	sorted := make([]string, len(usages))
	for k, v := range usages {
		sorted = append(sorted, fmt.Sprintf("  %s\n\t\t%s\n", strings.TrimSuffix(strings.Join(v, " "), " "), k))
	}

	sort.Slice(sorted, func(i, j int) bool {
		return sorted[i] < sorted[j]
	})

	if len(sorted) > 0 {
		builder.WriteString("\nFLAGS:\n")
	}

	for _, s := range sorted {
		builder.WriteString(s)
	}

	return builder.String()
}

func flagsUsage(flags *flag.FlagSet) string {
	flagCount := 0
	var builder strings.Builder
	if flags != nil {
		builder.WriteString("\nFLAGS:\n")
		flags.VisitAll(func(f *flag.Flag) {
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
			case "*flag.intValue", "int64Value":
				typeOf = " int"
			case "*flag.stringValue":
				typeOf = " string"
			case "*flag.uintValue", "uint64Value":
				typeOf = " uint"
			}

			dashes := "--"
			if len(f.Name) == 1 {
				dashes = "-"
			}
			builder.WriteString(fmt.Sprintf("  %s%s%s\n", dashes, f.Name, typeOf))
			builder.WriteString(fmt.Sprintf("\t\t%s%s\n", f.Usage, defaultVal))
		})
	}

	if flagCount == 0 {
		return ""
	}

	return builder.String()
}
