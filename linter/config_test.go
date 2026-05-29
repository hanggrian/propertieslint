package linter

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLoadConfigParsesGroupedRules(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, DefaultConfigFile)

	if err :=
		os.WriteFile(
			path,
			[]byte(
				`{
					"comment": {
						"comment-spaces": false
					},
					"format": false,
					"pair": {
						"duplicate-key": false,
						"missing-key": false
					},
					"whitespace": {
						"trailing-newline": false,
						"untrimmed-entry": false
					}
				}`,
			),
			0o644,
		); err != nil {
		t.Fatalf("write config: %v", err)
	}

	config, err := LoadConfig(path)
	if err != nil {
		t.Fatalf("LoadConfig returned error: %v", err)
	}

	if !config.CommentStyle {
		t.Fatalf("expected comment-style to remain enabled")
	}
	if config.CommentSpaces {
		t.Fatalf("expected comment-spaces to be disabled")
	}
	if config.InvalidEscape ||
		config.MissingSeparator ||
		config.UnterminatedContinuation {
		t.Fatalf("expected format group to be disabled, got %#v", config)
	}
	if config.DuplicateKey {
		t.Fatalf("expected duplicate-key to be disabled")
	}
	if !config.KeyName {
		t.Fatalf("expected key-name to remain enabled")
	}
	if config.MissingKey {
		t.Fatalf("expected missing-key to be disabled")
	}
	if !config.MissingValue {
		t.Fatalf("expected missing-value to remain enabled")
	}
	if config.TrailingNewline {
		t.Fatalf("expected trailing-newline to be disabled")
	}
	if config.UntrimmedEntry {
		t.Fatalf("expected untrimmed-entry to be disabled")
	}
	if !config.DuplicateBlankLine || !config.NoLeadingBlankLine {
		t.Fatalf("expected other whitespace rules to remain enabled")
	}
}
