package cmd

import "github.com/fatih/color"

func Bold(text string) string {
	return color.New(color.Bold).Sprint(text)
}

func Dim(text string) string {
	return color.New(color.Faint).Sprint(text)
}

func Italic(text string) string {
	return color.New(color.Italic).Sprint(text)
}

func Red(text string) string {
	return color.New(color.FgRed).Sprint(text)
}

func Green(text string) string {
	return color.New(color.FgGreen).Sprint(text)
}

func Yellow(text string) string {
	return color.New(color.FgYellow).Sprint(text)
}

func Blue(text string) string {
	return color.New(color.FgBlue).Sprint(text)
}

func Cyan(text string) string {
	return color.New(color.FgCyan).Sprint(text)
}
