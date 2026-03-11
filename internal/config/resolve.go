// Package config: merge file layers into effective prompts and aliases for list/show.
// Uses same resolution as review: explicit file only, or global + workspace (workspace overrides global).

package config

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
// so list/show include both. Keys here are names; values are the command string.
func BuiltinAliases() map[string]Alias {
	return map[string]Alias{
		"claude":       {Command: "claude --non-interactive"},
		"cursor-agent": {Command: "cursor-agent"},
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
