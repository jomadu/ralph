package config

import (
	"os"
	"path/filepath"
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

// TestResolveAICommand_configAliasName documents that the caller passes config-derived
// alias (e.g. eff.Loop.AICmdAlias) into ResolveAICommand as aliasName.
func TestResolveAICommand_configAliasName(t *testing.T) {
	resolved := &Resolved{
		Prompts: map[string]Prompt{},
		Aliases: map[string]Alias{"my-alias": {Command: "custom-ai --batch"}},
	}
	rootLoop := DefaultLoopSettings()
	rootLoop.AICmdAlias = "my-alias"
	eff := EffectiveWithBuiltins(RootEffective(resolved, rootLoop))
	aliasName := eff.Loop.AICmdAlias
	cmd, ok := ResolveAICommand(eff, "", aliasName)
	if !ok {
		t.Fatalf("ResolveAICommand(eff, \"\", eff.Loop.AICmdAlias) ok = false")
	}
	if cmd != "custom-ai --batch" {
		t.Errorf("ResolveAICommand(..., aliasName) = %q; want \"custom-ai --batch\"", cmd)
	}
}

// TestResolveAICommand_fromConfigFile asserts that when config is resolved from a file
// with loop.ai_cmd_alias set and no flags or env, the resolved command is the alias expansion
// (run/review use this path: Resolve → eff.Loop.AICmdAlias → ResolveAICommand).
func TestResolveAICommand_fromConfigFile(t *testing.T) {
	dir := t.TempDir()
	configPath := filepath.Join(dir, "ralph-config.yml")
	cfg := "loop:\n  ai_cmd_alias: cursor-agent\n"
	if err := os.WriteFile(configPath, []byte(cfg), 0644); err != nil {
		t.Fatal(err)
	}
	getenv := func(string) string { return "" }
	eff, ok, err := Resolve(getenv, dir, configPath, "")
	if err != nil {
		t.Fatalf("Resolve() err = %v", err)
	}
	if !ok || eff == nil {
		t.Fatalf("Resolve() = ok=%v eff=%v, want effective config", ok, eff)
	}
	if eff.Loop.AICmdAlias != "cursor-agent" {
		t.Fatalf("eff.Loop.AICmdAlias = %q, want cursor-agent (from config file)", eff.Loop.AICmdAlias)
	}
	directCmd := ""
	aliasName := eff.Loop.AICmdAlias
	command, ok := ResolveAICommand(eff, directCmd, aliasName)
	if !ok {
		t.Fatalf("ResolveAICommand(eff, \"\", eff.Loop.AICmdAlias) ok = false")
	}
	want := "agent -p --force --output-format stream-json --stream-partial-output"
	if command != want {
		t.Errorf("ResolveAICommand(from config file) = %q; want %q", command, want)
	}
}
