package rules_test

import "testing"

func keyOrderConfig() Config {
	config := DefaultConfig()
	config.KeyOrder = true
	return config
}

func TestKeyOrderAllowsSortedKeys(t *testing.T) {
	issues := lintSampleWithConfig(t, "alpha=1\nbeta=2\ngamma=3\n", keyOrderConfig())

	if len(issues) != 0 {
		t.Fatalf("expected no issues, got %#v", issues)
	}
}

func TestKeyOrderFlagsOutOfOrderKey(t *testing.T) {
	issues := lintSampleWithConfig(t, "beta=1\nalpha=2\n", keyOrderConfig())

	assertIssueAt(
		t,
		issues,
		2,
		1,
		`Arrange 'alpha' before 'beta'`,
	)
}

func TestKeyOrderIsCaseInsensitive(t *testing.T) {
	issues := lintSampleWithConfig(t, "Alpha=1\nbeta=2\nGamma=3\n", keyOrderConfig())

	if len(issues) != 0 {
		t.Fatalf("expected no issues, got %#v", issues)
	}
}

func TestKeyOrderAllowsCaseVariants(t *testing.T) {
	issues := lintSampleWithConfig(t, "item=1\nItem=2\n", keyOrderConfig())

	assertNoIssueContaining(t, issues, "Arrange")
}

func TestKeyOrderComparesUnescapedKeys(t *testing.T) {
	issues := lintSampleWithConfig(t, "aB=1\na\\u0041=2\n", keyOrderConfig())

	assertIssueAt(
		t,
		issues,
		2,
		1,
		`Arrange 'aA' before 'aB'`,
	)
}

func TestKeyOrderWithoutEscapeNormalization(t *testing.T) {
	config := keyOrderConfig()
	config.ValidEscape = false

	issues := lintSampleWithConfig(t, "a=1\n\\u0062=2\n", config)

	assertIssueAt(
		t,
		issues,
		2,
		1,
		`Arrange '\u0062' before 'a'`,
	)
}

func TestKeyOrderIgnoresCommentsAndBlankLines(t *testing.T) {
	issues := lintSampleWithConfig(t, "alpha=1\n# comment\n\nbeta=2\n", keyOrderConfig())

	if len(issues) != 0 {
		t.Fatalf("expected no issues, got %#v", issues)
	}
}

func TestKeyOrderReportsKeyStartColumn(t *testing.T) {
	issues := lintSampleWithConfig(t, "beta=1\n  alpha=2\n", keyOrderConfig())

	assertIssueAt(t, issues, 2, 3, "Arrange")
}

func TestKeyOrderAllowsDuplicateKeys(t *testing.T) {
	issues := lintSampleWithConfig(t, "alpha=1\nalpha=2\n", keyOrderConfig())

	assertNoIssueContaining(t, issues, "Arrange")
}

func TestKeyOrderIgnoresEmptyKeys(t *testing.T) {
	issues := lintSampleWithConfig(t, "=1\nalpha=2\n", keyOrderConfig())

	assertNoIssueContaining(t, issues, "Arrange")
}

func TestKeyOrderComparesContinuationKeys(t *testing.T) {
	issues := lintSampleWithConfig(t, "beta=1\nalpha\\\n=2\n", keyOrderConfig())

	assertIssueAt(t, issues, 2, 1, "Arrange")
}

func TestKeyOrderReportsFirstInversionOnly(t *testing.T) {
	issues := lintSampleWithConfig(t, "beta=1\nalpha=2\ngamma=3\n", keyOrderConfig())

	if len(issues) != 1 {
		t.Fatalf("expected one issue, got %#v", issues)
	}
	assertIssueAt(t, issues, 2, 1, "Arrange")
}

func TestKeyOrderAllowsPrefixesAndPunctuation(t *testing.T) {
	issues := lintSampleWithConfig(t, "foo=1\nfoo-bar=2\nfoo.bar=3\nfoo2=4\nfooBar=5\n", keyOrderConfig())

	if len(issues) != 0 {
		t.Fatalf("expected no issues, got %#v", issues)
	}
}

func TestKeyOrderCanBeDisabled(t *testing.T) {
	config := keyOrderConfig()
	config.KeyOrder = false

	issues := lintSampleWithConfig(t, "beta=1\nalpha=2\n", config)

	assertNoIssueContaining(t, issues, "Arrange")
}
