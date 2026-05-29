package cmd

import (
	"fmt"
	"os"
	"runtime/debug"
	"strings"

	"github.com/hanggrian/propertieslint/linter"
)

func Execute() error {
	// parse arguments
	args := os.Args[1:]
	quiet := false
	configPath := ""
	targets := make([]string, 0, len(args))
	for index := 0; index < len(args); index++ {
		arg := args[index]
		if arg == "-h" ||
			arg == "--help" {
			fmt.Printf("Lint Java Properties files\n\n")
			fmt.Printf("%s\n", Bold(Cyan("Usage:")))
			fmt.Printf(
				"   %s %s %s\n\n",
				Cyan("propertieslint"),
				Magenta("[PATHS]"),
				Blue("[OPTIONS]"),
			)
			fmt.Printf("%s\n", Bold(Magenta("Paths:")))
			fmt.Printf(
				"   %s      Supports %s file\n",
				Magenta("file"),
				Italic(".properties"),
			)
			fmt.Printf("   %s       Recursively find files in this directory\n", Magenta("dir"))
			fmt.Printf(
				"   %s   For example, %s for all properties files in this\n",
				Magenta("pattern"),
				Italic("*.properties"),
			)
			fmt.Printf("             directory, %s for all files\n\n", Italic("**/*"))
			fmt.Printf("%s\n", Bold(Blue("Options:")))
			fmt.Printf(
				"   %s, %s %s   Configuration file\n",
				Blue("-c"),
				Blue("--config"),
				Faint(Blue("<PATH>")),
			)
			fmt.Printf(
				"   %s, %s           Disable informational messages\n",
				Blue("-q"),
				Blue("--quiet"),
			)
			fmt.Printf(
				"   %s, %s            Display this message\n",
				Blue("-h"),
				Blue("--help"),
			)
			fmt.Printf(
				"   %s, %s         Show app version\n",
				Blue("-v"),
				Blue("--version"),
			)
			return nil
		}
		if arg == "-q" ||
			arg == "--quiet" {
			quiet = true
		}
		if arg == "-v" ||
			arg == "--version" {
			var version = "dev"
			if info, ok := debug.ReadBuildInfo(); ok {
				if info.Main.Version != "(devel)" {
					version = info.Main.Version
				}
			}
			fmt.Printf("propertieslint %s\n", Bold(version))
			return nil
		}
		switch {
		case arg == "-q" || arg == "--quiet":
			quiet = true

		case arg == "-c" || arg == "--config":
			index++
			if index >= len(args) {
				return fmt.Errorf("Missing value for %s", arg)
			}
			configPath = args[index]

		case strings.HasPrefix(arg, "--config="):
			configPath = strings.TrimPrefix(arg, "--config=")

		default:
			targets = append(targets, arg)
		}
	}

	// load config
	config, err := linter.LoadConfig(linter.ResolveConfigPath(configPath))
	if err != nil {
		return err
	}

	// run linter
	result, err := linter.Targets(targets, config)
	if err != nil {
		fmt.Fprintln(os.Stderr, Red(err.Error()))
		os.Exit(2)
	}

	// report results
	issues := len(result.Issues)
	if issues > 0 {
		fmt.Fprintf(os.Stderr, "%s\n", fmt.Sprintf("%d issue(s) found:", issues))
		for _, issue := range result.Issues {
			fmt.Printf("%s:%d:%d: %s\n", issue.Path, issue.Line, issue.Column, issue.Message)
		}
		os.Exit(1)
	}
	if result.CheckedFiles == 0 {
		if !quiet {
			fmt.Printf("%s\n", Yellow("No properties files found."))
		}
		return nil
	}
	if !quiet {
		fmt.Printf("%s\n", Green("All checks passed."))
	}
	return nil
}
