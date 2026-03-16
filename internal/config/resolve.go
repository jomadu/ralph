// Package config: merge file layers into effective prompts and aliases for list/show.
// Uses same resolution as review: explicit file only, or global + workspace (workspace overrides global).
// Prompt paths in config are resolved relative to the directory of the config file that defined them.

package config

import "path/filepath"

// resolvePromptPath makes a prompt path absolute when it is relative, using the
// directory of the config file that defined the prompt. Absolute paths are returned cleaned.
// configFilePath is the full path to the config file; empty is not used (path is left as-is).
func resolvePromptPath(configFilePath, promptPath string) string {
	if promptPath == "" {
		return ""
	}
	if filepath.IsAbs(promptPath) {
		return filepath.Clean(promptPath)
	}
	if configFilePath == "" {
		return promptPath
	}
	return filepath.Join(filepath.Dir(configFilePath), promptPath)
}

// Effective is the full effective config for a run: loop settings, prompts, and aliases.
// Used by run-loop, CLI, and show config when a single merged config is needed.
type Effective struct {
	Loop    LoopSettings
	Prompts map[string]Prompt
	Aliases map[string]Alias
}

// Resolved holds merged prompts and aliases from the resolved config (one explicit layer
// or global + workspace with workspace overriding). Used by list and show.
type Resolved struct {
	Prompts map[string]Prompt
	Aliases map[string]Alias
}

// MergeLayers merges global and workspace layers into a single Resolved config.
// Workspace keys override global for the same name. Either or both may be nil.
// globalPath and workspacePath are the full paths to the config files; they are used
// to resolve relative prompt paths relative to the defining config file's directory.
func MergeLayers(global *FileLayer, globalPath string, workspace *FileLayer, workspacePath string) *Resolved {
	r := &Resolved{
		Prompts: make(map[string]Prompt),
		Aliases: make(map[string]Alias),
	}
	if global != nil {
		for k, v := range global.Prompts {
			p := v
			if p.Path != "" {
				p.Path = resolvePromptPath(globalPath, p.Path)
			}
			r.Prompts[k] = p
		}
		for k, v := range global.Aliases {
			r.Aliases[k] = v
		}
	}
	if workspace != nil {
		for k, v := range workspace.Prompts {
			p := v
			if p.Path != "" {
				p.Path = resolvePromptPath(workspacePath, p.Path)
			}
			r.Prompts[k] = p
		}
		for k, v := range workspace.Aliases {
			r.Aliases[k] = v
		}
	}
	return r
}

// BuiltinAliases returns default AI command aliases (R004). Merged with user config
// so list/show include both. Commands per docs/engineering/components/backend.md.
func BuiltinAliases() map[string]Alias {
	return map[string]Alias{
		"claude":       {Command: "claude -p --dangerously-skip-permissions"},
		"kiro":         {Command: "kiro-cli chat --no-interactive --trust-all-tools"},
		"copilot":      {Command: "copilot --yolo"},
		"cursor-agent": {Command: "agent -p --force --output-format stream-json --stream-partial-output"},
	}
}

// ResolvedWithBuiltins returns a copy of resolved with built-in aliases merged in.
// User aliases override built-ins for the same name.
func ResolvedWithBuiltins(r *Resolved) *Resolved {
	out := &Resolved{
		Prompts: r.Prompts,
		Aliases: make(map[string]Alias),
	}
	for k, v := range BuiltinAliases() {
		out.Aliases[k] = v
	}
	for k, v := range r.Aliases {
		out.Aliases[k] = v
	}
	return out
}

// EffectiveWithBuiltins returns a copy of e with built-in aliases merged into Aliases.
// User aliases override built-ins for the same name. Use when the Effective is the
// single resolved config for run-loop, review, list, or show (O002/R004, R007).
func EffectiveWithBuiltins(e *Effective) *Effective {
	if e == nil {
		return nil
	}
	out := &Effective{
		Loop:    e.Loop,
		Prompts: make(map[string]Prompt),
		Aliases: make(map[string]Alias),
	}
	for k, v := range e.Prompts {
		out.Prompts[k] = v
	}
	for k, v := range BuiltinAliases() {
		out.Aliases[k] = v
	}
	for k, v := range e.Aliases {
		out.Aliases[k] = v
	}
	return out
}

// EffectiveForPrompt builds an Effective config for a named prompt with merge order:
// defaults → global → workspace → explicit file → env → prompt overrides (CLI in Phase 4).
// rootLoop must already be MergeRootLoop(global, workspace) then ApplyEnvOverlayToLoop(_, env).
// If promptName is empty or not found, returns nil, false. Otherwise returns a copy of
// resolved with Effective.Loop set to the prompt's effective loop (root + prompt overrides).
func EffectiveForPrompt(resolved *Resolved, promptName string, rootLoop LoopSettings) (*Effective, bool) {
	if resolved == nil || promptName == "" {
		return nil, false
	}
	prompt, ok := resolved.Prompts[promptName]
	if !ok {
		return nil, false
	}
	loop := EffectiveLoopForPrompt(rootLoop, &prompt)
	// Copy maps so caller cannot mutate resolved
	prompts := make(map[string]Prompt)
	for k, v := range resolved.Prompts {
		prompts[k] = v
	}
	aliases := make(map[string]Alias)
	for k, v := range resolved.Aliases {
		aliases[k] = v
	}
	return &Effective{Loop: loop, Prompts: prompts, Aliases: aliases}, true
}

// RootEffective builds an Effective with no prompt selected: root loop (defaults → layers → env),
// and merged prompts/aliases from resolved. rootLoop must already have env applied if desired.
func RootEffective(resolved *Resolved, rootLoop LoopSettings) *Effective {
	if resolved == nil {
		return nil
	}
	loop := rootLoop
	prompts := make(map[string]Prompt)
	for k, v := range resolved.Prompts {
		prompts[k] = v
	}
	aliases := make(map[string]Alias)
	for k, v := range resolved.Aliases {
		aliases[k] = v
	}
	return &Effective{Loop: loop, Prompts: prompts, Aliases: aliases}
}

// ResolveAICommand resolves the AI command string from effective config and optional
// overrides. Used by run-loop and review so the backend receives the resolved command
// (T2.3, O002/R004). Precedence: directCmd (if non-empty) > aliasName (looked up in
// eff.Aliases). Returns the command string and true, or empty string and false if
// alias name is missing or unknown. eff must include built-in aliases (e.g. from
// EffectiveWithBuiltins).
func ResolveAICommand(eff *Effective, directCmd, aliasName string) (command string, ok bool) {
	if eff == nil {
		return "", false
	}
	if directCmd != "" {
		return directCmd, true
	}
	if aliasName == "" {
		return "", false
	}
	a, ok := eff.Aliases[aliasName]
	if !ok || a.Command == "" {
		return "", false
	}
	return a.Command, true
}
