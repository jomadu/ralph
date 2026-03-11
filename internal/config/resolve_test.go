package config

import (
	"testing"
)

func TestResolveAICommand_nilEffective(t *testing.T) {
	cmd, ok := ResolveAICommand(nil, "", "")
	if ok || cmd != "" {
		t.Errorf("ResolveAICommand(nil, ..., ...) = %q, %v; want \"\", false", cmd, ok)
	}
}

func TestResolveAICommand_directCmdOverrides(t *testing.T) {
	eff := DefaultEffective()
	cmd, ok := ResolveAICommand(eff, "my-custom-cmd --arg", "claude")
	if !ok || cmd != "my-custom-cmd --arg" {
		t.Errorf("ResolveAICommand(direct) = %q, %v; want \"my-custom-cmd --arg\", true", cmd, ok)
	}
}

func TestResolveAICommand_builtinAlias(t *testing.T) {
	eff := DefaultEffective()
	cmd, ok := ResolveAICommand(eff, "", "cursor-agent")
	if !ok {
		t.Fatalf("ResolveAICommand(eff, \"\", \"cursor-agent\") ok = false")
	}
	want := "agent -p --force --output-format stream-json --stream-partial-output"
	if cmd != want {
		t.Errorf("ResolveAICommand(..., \"cursor-agent\") = %q; want %q", cmd, want)
	}
}

func TestResolveAICommand_allBuiltinAliases(t *testing.T) {
	eff := DefaultEffective()
	aliases := []string{"claude", "kiro", "copilot", "cursor-agent"}
	for _, name := range aliases {
		cmd, ok := ResolveAICommand(eff, "", name)
		if !ok || cmd == "" {
			t.Errorf("ResolveAICommand(eff, \"\", %q) = %q, %v; want non-empty command, true", name, cmd, ok)
		}
	}
}

func TestResolveAICommand_emptyAliasName(t *testing.T) {
	eff := DefaultEffective()
	cmd, ok := ResolveAICommand(eff, "", "")
	if ok || cmd != "" {
		t.Errorf("ResolveAICommand(eff, \"\", \"\") = %q, %v; want \"\", false", cmd, ok)
	}
}

func TestResolveAICommand_unknownAlias(t *testing.T) {
	eff := DefaultEffective()
	cmd, ok := ResolveAICommand(eff, "", "nonexistent-alias")
	if ok || cmd != "" {
		t.Errorf("ResolveAICommand(eff, \"\", \"nonexistent-alias\") = %q, %v; want \"\", false", cmd, ok)
	}
}

func TestResolveAICommand_userAliasOverridesBuiltin(t *testing.T) {
	eff := &Effective{
		Loop:    DefaultLoopSettings(),
		Prompts: make(map[string]Prompt),
		Aliases: make(map[string]Alias),
	}
	for k, v := range BuiltinAliases() {
		eff.Aliases[k] = v
	}
	eff.Aliases["claude"] = Alias{Command: "claude-custom --my-flag"}
	cmd, ok := ResolveAICommand(eff, "", "claude")
	if !ok {
		t.Fatalf("ResolveAICommand(user override) ok = false")
	}
	if cmd != "claude-custom --my-flag" {
		t.Errorf("ResolveAICommand(user override) = %q; want \"claude-custom --my-flag\"", cmd)
	}
}
