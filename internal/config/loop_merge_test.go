package config

import (
	"testing"
)

func TestApplyLoopSection_nil(t *testing.T) {
	base := DefaultLoopSettings()
	got := ApplyLoopSection(base, nil)
	if got.MaxIterations != base.MaxIterations {
		t.Errorf("ApplyLoopSection(base, nil) modified base")
	}
}

func TestApplyLoopSection_overrides(t *testing.T) {
	base := DefaultLoopSettings()
	section := &LoopSection{
		MaxIterations:    intPtr(5),
		FailureThreshold: intPtr(2),
		LogLevel:         "debug",
		MaxOutputBuffer:  intPtr(32768),
	}
	got := ApplyLoopSection(base, section)
	if got.MaxIterations != 5 {
		t.Errorf("MaxIterations = %d, want 5", got.MaxIterations)
	}
	if got.FailureThreshold != 2 {
		t.Errorf("FailureThreshold = %d, want 2", got.FailureThreshold)
	}
	if got.LogLevel != "debug" {
		t.Errorf("LogLevel = %q, want debug", got.LogLevel)
	}
	if got.MaxOutputBuffer != 32768 {
		t.Errorf("MaxOutputBuffer = %d, want 32768", got.MaxOutputBuffer)
	}
	if got.SuccessSignal != base.SuccessSignal {
		t.Errorf("SuccessSignal changed unexpectedly")
	}
}

func TestMergeRootLoop_defaultsOnly(t *testing.T) {
	got := MergeRootLoop(nil, nil)
	d := DefaultLoopSettings()
	if got.MaxIterations != d.MaxIterations || got.FailureThreshold != d.FailureThreshold {
		t.Errorf("MergeRootLoop(nil,nil) = %+v, want defaults", got)
	}
}

func TestMergeRootLoop_workspaceOverridesGlobal(t *testing.T) {
	global := &FileLayer{
		Loop: &LoopSection{MaxIterations: intPtr(20), LogLevel: "info"},
	}
	workspace := &FileLayer{
		Loop: &LoopSection{MaxIterations: intPtr(3), LogLevel: "debug"},
	}
	got := MergeRootLoop(global, workspace)
	if got.MaxIterations != 3 {
		t.Errorf("MaxIterations = %d, want 3 (workspace overrides global)", got.MaxIterations)
	}
	if got.LogLevel != "debug" {
		t.Errorf("LogLevel = %q, want debug", got.LogLevel)
	}
}

func TestApplyEnvOverlayToLoop_nil(t *testing.T) {
	loop := DefaultLoopSettings()
	got := ApplyEnvOverlayToLoop(loop, nil)
	if got.MaxIterations != loop.MaxIterations {
		t.Errorf("overlay nil should not change loop")
	}
}

func TestApplyEnvOverlayToLoop_overrides(t *testing.T) {
	loop := LoopSettings{MaxIterations: 10, FailureThreshold: 3, LogLevel: "info"}
	warn := "warn"
	overlay := &EnvOverlay{
		MaxIterations:   intPtr(7),
		LogLevel:        &warn,
		MaxOutputBuffer: intPtr(8192),
	}
	got := ApplyEnvOverlayToLoop(loop, overlay)
	if got.MaxIterations != 7 {
		t.Errorf("MaxIterations = %d, want 7", got.MaxIterations)
	}
	if got.LogLevel != "warn" {
		t.Errorf("LogLevel = %q, want warn", got.LogLevel)
	}
	if got.MaxOutputBuffer != 8192 {
		t.Errorf("MaxOutputBuffer = %d, want 8192", got.MaxOutputBuffer)
	}
	if got.FailureThreshold != 3 {
		t.Errorf("FailureThreshold = %d, want 3 (unchanged)", got.FailureThreshold)
	}
}

func TestEffectiveLoopForPrompt_noOverride(t *testing.T) {
	root := DefaultLoopSettings()
	root.MaxIterations = 8
	prompt := &Prompt{Path: "p.md"}
	got := EffectiveLoopForPrompt(root, prompt)
	if got.MaxIterations != 8 {
		t.Errorf("MaxIterations = %d, want 8", got.MaxIterations)
	}
}

func TestEffectiveLoopForPrompt_withOverride(t *testing.T) {
	root := DefaultLoopSettings()
	root.MaxIterations = 10
	root.LogLevel = "info"
	prompt := &Prompt{
		Path: "p.md",
		Loop: &LoopSection{MaxIterations: intPtr(2), LogLevel: "debug"},
	}
	got := EffectiveLoopForPrompt(root, prompt)
	if got.MaxIterations != 2 {
		t.Errorf("MaxIterations = %d, want 2 (prompt override)", got.MaxIterations)
	}
	if got.LogLevel != "debug" {
		t.Errorf("LogLevel = %q, want debug", got.LogLevel)
	}
}

func TestEffectiveForPrompt_notFound(t *testing.T) {
	r := &Resolved{Prompts: map[string]Prompt{"a": {Path: "a.md"}}, Aliases: map[string]Alias{}}
	root := DefaultLoopSettings()
	eff, ok := EffectiveForPrompt(r, "b", root)
	if ok || eff != nil {
		t.Errorf("EffectiveForPrompt(_, \"b\", _) = %v, %v; want nil, false", eff, ok)
	}
}

func TestEffectiveForPrompt_found(t *testing.T) {
	r := &Resolved{
		Prompts: map[string]Prompt{
			"p1": {Path: "p1.md", Loop: &LoopSection{MaxIterations: intPtr(1)}},
		},
		Aliases: map[string]Alias{},
	}
	root := DefaultLoopSettings()
	root.MaxIterations = 10
	eff, ok := EffectiveForPrompt(r, "p1", root)
	if !ok || eff == nil {
		t.Fatalf("EffectiveForPrompt(_, \"p1\", _) = %v, %v; want non-nil, true", eff, ok)
	}
	if eff.Loop.MaxIterations != 1 {
		t.Errorf("eff.Loop.MaxIterations = %d, want 1 (prompt override)", eff.Loop.MaxIterations)
	}
	if eff.Prompts["p1"].Path != "p1.md" {
		t.Errorf("Prompts[\"p1\"].Path = %q", eff.Prompts["p1"].Path)
	}
}

func TestRootEffective(t *testing.T) {
	r := &Resolved{Prompts: map[string]Prompt{"x": {Path: "x.md"}}, Aliases: map[string]Alias{}}
	root := DefaultLoopSettings()
	eff := RootEffective(r, root)
	if eff == nil {
		t.Fatal("RootEffective = nil")
	}
	if eff.Loop.MaxIterations != root.MaxIterations {
		t.Errorf("Loop = %+v, want %+v", eff.Loop, root)
	}
	if eff.Prompts["x"].Path != "x.md" {
		t.Errorf("Prompts[\"x\"].Path = %q", eff.Prompts["x"].Path)
	}
}
