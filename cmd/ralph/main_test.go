package main

import (
	"bytes"
	"strings"
	"testing"
)

// TestShowPromptGuide_success verifies that "ralph show prompt-guide" exits 0 and stdout
// contains the full guide: all four dimension names, the summary table, and the closing
// line that references ralph show prompt-guide (PLAN T4, ralph-uc5).
func TestShowPromptGuide_success(t *testing.T) {
	root := newRoot()
	out := new(bytes.Buffer)
	root.SetOut(out)
	root.SetArgs([]string{"show", "prompt-guide"})
	if err := root.Execute(); err != nil {
		t.Fatalf("ralph show prompt-guide: %v", err)
	}
	got := out.String()

	// Four dimension names (from docs/writing-ralph-prompts.md).
	for _, phrase := range []string{
		"Signal and state",
		"Iteration awareness",
		"Scope and convergence",
		"Subjective completion",
	} {
		if !strings.Contains(got, phrase) {
			t.Errorf("stdout missing dimension %q; got %d bytes", phrase, len(got))
		}
	}
	// Summary table.
	if !strings.Contains(got, "| Dimension |") || !strings.Contains(got, "What to address") {
		t.Errorf("stdout missing summary table; got %d bytes", len(got))
	}
	// Closing line that references ralph show prompt-guide.
	if !strings.Contains(got, "ralph show prompt-guide") {
		t.Errorf("stdout missing closing line referencing ralph show prompt-guide; got %d bytes", len(got))
	}
}

// TestShowPromptGuide_markdown verifies that "ralph show prompt-guide --markdown" outputs
// the full guide and contains markdown (e.g. ## or **) (PLAN T4, ralph-uc5).
func TestShowPromptGuide_markdown(t *testing.T) {
	root := newRoot()
	out := new(bytes.Buffer)
	root.SetOut(out)
	root.SetArgs([]string{"show", "prompt-guide", "--markdown"})
	if err := root.Execute(); err != nil {
		t.Fatalf("ralph show prompt-guide --markdown: %v", err)
	}
	got := out.String()

	// Same full-guide content: four dimensions.
	for _, phrase := range []string{
		"Signal and state",
		"Iteration awareness",
		"Scope and convergence",
		"Subjective completion",
	} {
		if !strings.Contains(got, phrase) {
			t.Errorf("stdout with --markdown missing dimension %q", phrase)
		}
	}
	// Contains markdown.
	if !strings.Contains(got, "##") && !strings.Contains(got, "**") {
		t.Errorf("stdout with --markdown should contain markdown (## or **); got %d bytes", len(got))
	}
}

// TestShowPromptGuide_invalidUsage verifies that an unexpected positional argument
// causes a non-zero exit (PLAN T4, ralph-uc5).
func TestShowPromptGuide_invalidUsage(t *testing.T) {
	root := newRoot()
	errOut := new(bytes.Buffer)
	root.SetErr(errOut)
	root.SetArgs([]string{"show", "prompt-guide", "foo"})
	err := root.Execute()
	if err == nil {
		t.Fatal("expected error for unexpected positional arg; got nil")
	}
	if !strings.Contains(err.Error(), "unexpected argument") {
		t.Errorf("error should mention unexpected argument; got %q", err.Error())
	}
}
