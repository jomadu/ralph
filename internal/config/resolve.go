// Package config: merge file layers into effective prompts and aliases for list/show.
// Uses same resolution as review: explicit file only, or global + workspace (workspace overrides global).

package config

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
func MergeLayers(global, workspace *FileLayer) *Resolved {
	r := &Resolved{
		Prompts: make(map[string]Prompt),
		Aliases: make(map[string]Alias),
	}
	if global != nil {
		for k, v := range global.Prompts {
			r.Prompts[k] = v
		}
		for k, v := range global.Aliases {
			r.Aliases[k] = v
		}
	}
	if workspace != nil {
		for k, v := range workspace.Prompts {
			r.Prompts[k] = v
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
