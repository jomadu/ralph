// Package config resolves configuration for the current context (cwd, explicit path, env)
// into a single effective config used by run-loop, review, list, and show (O002/R007).
//
// Single entrypoint: Resolve — pass getenv (e.g. os.Getenv), cwd, optional config file path,
// and optional prompt name; get back (*Effective, ok, error) with built-in aliases included.
package config

import "path/filepath"

// LoadGlobalAndWorkspace loads the global and workspace config file layers from the
// given paths. Callers resolve paths via GlobalPath(getenv) and WorkspacePath(cwd).
// Missing files are skipped without error (returns nil for that layer).
// Returns (globalLayer, workspaceLayer, error). Error is only for read/parse failure.
func LoadGlobalAndWorkspace(globalPath, workspacePath string) (global, workspace *FileLayer, err error) {
	global, err = ReadLayer(globalPath)
	if err != nil {
		return nil, nil, err
	}
	workspace, err = ReadLayer(workspacePath)
	if err != nil {
		return nil, nil, err
	}
	return global, workspace, nil
}

// LoadExplicit loads only the config file at the given path. Use when the user
// supplies an explicit config path (e.g. CLI --config). Global and workspace
// config are not read. The file must exist and be readable; if it is missing,
// a directory, or unreadable, returns an error (O002/R005). Path may be
// relative to the current working directory or absolute.
func LoadExplicit(path string) (*FileLayer, error) {
	return ReadLayerRequired(path)
}

// loadLayersAndRootLoop loads file layers and applies env overlay, returning
// merged resolved config and root loop (defaults → layers → env). Used by
// ResolveEffective and ResolveEffectiveForPrompt.
func loadLayersAndRootLoop(getenv func(string) string, cwd, configPath string) (*Resolved, LoopSettings, error) {
	var resolved *Resolved
	var rootLoop LoopSettings
	if configPath != "" {
		path := configPath
		if !filepath.IsAbs(path) {
			path = filepath.Join(cwd, path)
		}
		layer, err := LoadExplicit(path)
		if err != nil {
			return nil, LoopSettings{}, err
		}
		resolved = MergeLayers(nil, "", layer, path)
		rootLoop = MergeRootLoop(nil, layer)
	} else {
		globalPath := GlobalPath(getenv)
		workspacePath := WorkspacePath(cwd)
		global, workspace, err := LoadGlobalAndWorkspace(globalPath, workspacePath)
		if err != nil {
			return nil, LoopSettings{}, err
		}
		resolved = MergeLayers(global, globalPath, workspace, workspacePath)
		rootLoop = MergeRootLoop(global, workspace)
	}
	overlay, err := ParseEnvOverlay(getenv)
	if err != nil {
		return nil, LoopSettings{}, err
	}
	rootLoop = ApplyEnvOverlayToLoop(rootLoop, overlay)
	return resolved, rootLoop, nil
}

// ResolveEffective loads config (explicit path or global+workspace), applies
// RALPH_LOOP_* environment overlay, and returns the effective config (root
// loop only; no prompt selected). RALPH_CONFIG_HOME is used when resolving
// global path (see GlobalPath). Invalid env values produce a clear error
// (O010/R004). getenv is typically os.Getenv. When configPath is non-empty,
// only that file is used for file-based config.
func ResolveEffective(getenv func(string) string, cwd, configPath string) (*Effective, error) {
	resolved, rootLoop, err := loadLayersAndRootLoop(getenv, cwd, configPath)
	if err != nil {
		return nil, err
	}
	return RootEffective(resolved, rootLoop), nil
}

// ResolveEffectiveForPrompt resolves config for a given prompt name with merge
// order: defaults → global → workspace → explicit file → env → prompt overrides
// (CLI in Phase 4). When promptName is non-empty and the prompt exists, the
// returned Effective has that prompt's loop overrides applied. When promptName
// is empty, returns root effective (same as ResolveEffective). When promptName
// is set but not found, returns (nil, false, nil). Errors from load or env
// parsing are returned as (nil, false, err).
func ResolveEffectiveForPrompt(getenv func(string) string, cwd, configPath, promptName string) (*Effective, bool, error) {
	resolved, rootLoop, err := loadLayersAndRootLoop(getenv, cwd, configPath)
	if err != nil {
		return nil, false, err
	}
	if promptName == "" {
		return RootEffective(resolved, rootLoop), true, nil
	}
	eff, ok := EffectiveForPrompt(resolved, promptName, rootLoop)
	if !ok {
		return nil, false, nil
	}
	return eff, true, nil
}

// Resolve is the single entrypoint to resolve effective config for the current context
// (cwd, explicit config path, env). It returns the Effective config used by run-loop,
// review, list, and show (O002/R007). Built-in aliases are included in the returned
// Effective; user aliases override built-ins for the same name.
//
// Parameters:
//   - getenv: typically os.Getenv; used for RALPH_CONFIG_HOME and RALPH_LOOP_*.
//   - cwd: current working directory (global/workspace paths when configPath is empty).
//   - configPath: explicit config file path; if non-empty, only that file is used (no global/workspace).
//   - promptName: optional prompt name; when non-empty and the prompt exists, Effective.Loop
//     includes that prompt's overrides; when empty, root loop is returned.
//
// Returns (*Effective, true, nil) on success; (nil, false, nil) when promptName is set but
// the prompt is not found; (nil, false, err) on load or env parse error.
func Resolve(getenv func(string) string, cwd, configPath, promptName string) (*Effective, bool, error) {
	eff, ok, err := ResolveEffectiveForPrompt(getenv, cwd, configPath, promptName)
	if err != nil || !ok || eff == nil {
		return eff, ok, err
	}
	return EffectiveWithBuiltins(eff), true, nil
}
