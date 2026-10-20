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
						"comment-spacing": false,
						"comment-trimming": false,
						"no-empty-comment": false
					},
					"format": false,
					"pair": {
						"key-order": false,
						"no-duplicate-key": false,
						"no-empty-key": false,
						"no-empty-value": true
					},
					"whitespace": {
						"trailing-newline": false,
						"entry-trimming": false
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
	if config.CommentSpacing {
		t.Fatalf("expected comment-spacing to be disabled")
	}
	if config.CommentTrimming {
		t.Fatalf("expected comment-trimming to be disabled")
	}
	if config.NoEmptyComment {
		t.Fatalf("expected no-empty-comment to be disabled")
	}
	if config.ValidEscape ||
		config.NoDanglingContinuation ||
		config.SeparatorStyle {
		t.Fatalf("expected format group to be disabled, got %#v", config)
	}
	if config.NoDuplicateKey {
		t.Fatalf("expected no-duplicate-key to be disabled")
	}
	if config.KeyOrder {
		t.Fatalf("expected key-order to be disabled")
	}
	if !config.KeyName {
		t.Fatalf("expected key-name to remain enabled")
	}
	if config.NoEmptyKey {
		t.Fatalf("expected no-empty-key to be disabled")
	}
	if !config.NoEmptyValue {
		t.Fatalf("expected no-empty-value to remain enabled")
	}
	if config.TrailingNewline {
		t.Fatalf("expected trailing-newline to be disabled")
	}
	if config.EntryTrimming {
		t.Fatalf("expected entry-trimming to be disabled")
	}
	if !config.NoDuplicateBlankLine || !config.NoLeadingBlankLine {
		t.Fatalf("expected other whitespace rules to remain enabled")
	}
}

func TestLoadConfigParsesSeparatorStyle(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, DefaultConfigFile)

	if err :=
		os.WriteFile(
			path,
			[]byte(`{"format": {"separator-style": false}}`),
			0o644,
		); err != nil {
		t.Fatalf("write config: %v", err)
	}

	config, err := LoadConfig(path)
	if err != nil {
		t.Fatalf("LoadConfig returned error: %v", err)
	}

	if config.SeparatorStyle {
		t.Fatalf("expected separator-style to be disabled")
	}
	if !config.ValidEscape || !config.NoDanglingContinuation {
		t.Fatalf("expected other format rules to remain enabled, got %#v", config)
	}
}

func TestLoadConfigParsesNoDanglingContinuation(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, DefaultConfigFile)

	if err :=
		os.WriteFile(
			path,
			[]byte(`{"format": {"no-dangling-continuation": false}}`),
			0o644,
		); err != nil {
		t.Fatalf("write config: %v", err)
	}

	config, err := LoadConfig(path)
	if err != nil {
		t.Fatalf("LoadConfig returned error: %v", err)
	}

	if config.NoDanglingContinuation {
		t.Fatalf("expected no-dangling-continuation to be disabled")
	}
	if !config.ValidEscape || !config.SeparatorStyle {
		t.Fatalf("expected other format rules to remain enabled, got %#v", config)
	}
}

func TestLoadConfigParsesValidEscape(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, DefaultConfigFile)

	if err :=
		os.WriteFile(
			path,
			[]byte(`{"format": {"valid-escape": false}}`),
			0o644,
		); err != nil {
		t.Fatalf("write config: %v", err)
	}

	config, err := LoadConfig(path)
	if err != nil {
		t.Fatalf("LoadConfig returned error: %v", err)
	}

	if config.ValidEscape {
		t.Fatalf("expected valid-escape to be disabled")
	}
	if !config.NoDanglingContinuation || !config.SeparatorStyle {
		t.Fatalf("expected other format rules to remain enabled, got %#v", config)
	}
}
