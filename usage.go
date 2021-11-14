package geenee

import (
	"flag"
	"fmt"
	"sort"
	"strings"
	"text/tabwriter"
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
			builder.WriteString(mergeFlagsUsage(command.Flags, command.root))
		} else {
			builder.WriteString(flagsUsage(command.Flags, command.root))
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

func mergeFlagsUsage(flags *flag.FlagSet, isRoot bool) string {
	usages := make(map[string][]string)

	var builder strings.Builder
	tabWriter := tabwriter.NewWriter(&builder, 0, 0, 4, ' ', tabwriter.DiscardEmptyColumns)
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
				temp := []string{dashedFlag, "\t" + typeOf}
				usages[usage] = temp
			}
		})
	}

	if isRoot {
		//we add the version flag to root since we support it automatically
		usages["display version information"] = []string{"--version", "\t"}
	}

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
		tabWriter.Write([]byte(s))
	}

	tabWriter.Flush()
	return builder.String()
}

func flagsUsage(flags *flag.FlagSet, isRoot bool) string {
	flagCount := 0
	var builder strings.Builder
	var usages []string
	tabWriter := tabwriter.NewWriter(&builder, 0, 0, 4, ' ', tabwriter.DiscardEmptyColumns)
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
			usages = append(usages, fmt.Sprintf("%s%s\t%s\t%s%s\n", dashes, f.Name, typeOf, f.Usage, defaultVal))
		})
	}

	if isRoot {
		//we add the version flag to root since we support it automatically
		flagCount++
		usages = append(usages, fmt.Sprintf("%s\t%s\t%s\n", "--version", "", "display version information"))
	}

	if flagCount == 0 {
		return ""
	}

	sort.Slice(usages, func(i, j int) bool {
		return usages[i] < usages[j]
	})

	for _, s := range usages {
		tabWriter.Write([]byte(s))
	}

	tabWriter.Flush()
	return builder.String()
}
