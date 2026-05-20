package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/hanggrian/propertieslint/propertieslint"
)

var Version = "0.1.0"

func Execute() error {
	args := os.Args[1:]
	for _, arg := range args {
		if arg == "-h" ||
			arg == "--help" {
			fmt.Printf("Lint Java Properties files\n\n")
			fmt.Printf("\U0001f680 %s\n", Bold("Usage:"))
			fmt.Printf("   propertieslint %s %s\n\n", Cyan("<paths>"), Blue("[options]"))
			fmt.Printf("\U0001f4c4 %s\n", Bold(Cyan("Paths:")))
			fmt.Printf("   file      Supports %s file\n", Italic(".properties"))
			fmt.Printf("   dir       Recursively find files in this directory\n")
			fmt.Printf(
				"   pattern   For example, %s for all properties files in this\n",
				Italic("*.properties"),
			)
			fmt.Printf("             directory, %s for all files\n\n", Italic("**/*"))
			fmt.Printf("\u2699\ufe0f  %s\n", Bold(Blue("Options:")))
			fmt.Printf("   -c  [ --config ] arg   Configuration file\n")
			fmt.Printf("   -h  [ --help ]         Display this message\n")
			fmt.Printf("   -v  [ --version ]      Show app version\n")
			return nil
		}
		if arg == "-v" ||
			arg == "--version" {
			fmt.Printf("propertieslint %s\n", Bold(Version))
			return nil
		}
	}

	configPath, targets, err := parseArgs(args)
	if err != nil {
		return err
	}

	configPath = propertieslint.ResolveConfigPath(configPath)
	config, err := propertieslint.LoadConfig(configPath)
	if err != nil {
		return err
	}

	result, err := propertieslint.Targets(targets, config)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(2)
	}

	for _, issue := range result.Issues {
		fmt.Printf("%s:%d:%d: %s\n", issue.Path, issue.Line, issue.Column, issue.Message)
	}

	if len(result.Issues) > 0 {
		os.Exit(1)
	}

	if result.CheckedFiles == 0 {
		fmt.Fprintln(os.Stderr, "no properties files found")
		os.Exit(1)
	}

	fmt.Printf("lint ok: %d file(s) checked\n", result.CheckedFiles)
	return nil
}

func parseArgs(args []string) (string, []string, error) {
	configPath := ""
	targets := make([]string, 0, len(args))

	for index := 0; index < len(args); index++ {
		arg := args[index]
		switch {
		case arg == "-c" || arg == "--config":
			index++
			if index >= len(args) {
				return "", nil, fmt.Errorf("missing value for %s", arg)
			}
			configPath = args[index]
		case strings.HasPrefix(arg, "--config="):
			configPath = strings.TrimPrefix(arg, "--config=")
		default:
			targets = append(targets, arg)
		}
	}

	return configPath, targets, nil
}
