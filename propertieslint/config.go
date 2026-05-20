package propertieslint

import (
	"encoding/json"
	"fmt"
	"os"
)

const DefaultConfigFile = "propertieslint.json"

type Config struct {
	MissingSeparator         bool
	MissingValue             bool
	DuplicateKey             bool
	InvalidEscape            bool
	UnterminatedContinuation bool
	UntrimmedEntry           bool
	DuplicateBlankLine       bool
	NoLeadingBlankLine       bool
	TrailingNewline          bool
}

func DefaultConfig() Config {
	return Config{
		MissingSeparator:         true,
		MissingValue:             true,
		DuplicateKey:             true,
		InvalidEscape:            true,
		UnterminatedContinuation: true,
		UntrimmedEntry:           true,
		DuplicateBlankLine:       true,
		NoLeadingBlankLine:       true,
		TrailingNewline:          true,
	}
}

func LoadConfig(path string) (Config, error) {
	config := DefaultConfig()
	if path == "" {
		return config, nil
	}

	data, err := os.ReadFile(path)
	if err != nil {
		return Config{}, err
	}

	rules := map[string]bool{}
	if err := json.Unmarshal(data, &rules); err != nil {
		return Config{}, fmt.Errorf("parse config %q: %w", path, err)
	}

	for name, enabled := range rules {
		switch name {
		case "missing-key-value-separator":
			config.MissingSeparator = enabled
		case "missing-value":
			config.MissingValue = enabled
		case "duplicate-key":
			config.DuplicateKey = enabled
		case "invalid-escape":
			config.InvalidEscape = enabled
		case "unterminated-line-continuation":
			config.UnterminatedContinuation = enabled
		case "untrimmed-entry":
			config.UntrimmedEntry = enabled
		case "duplicate-blank-line":
			config.DuplicateBlankLine = enabled
		case "no-leading-blank-line":
			config.NoLeadingBlankLine = enabled
		case "trailing-newline":
			config.TrailingNewline = enabled
		default:
			return Config{}, fmt.Errorf("unknown rule %q", name)
		}
	}
	return config, nil
}

func ResolveConfigPath(explicit string) string {
	if explicit != "" {
		return explicit
	}
	if _, err := os.Stat(DefaultConfigFile); err == nil {
		return DefaultConfigFile
	}
	return ""
}
