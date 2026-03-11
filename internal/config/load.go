package config

// LoadGlobalAndWorkspace loads the global and workspace config file layers.
// getenv is typically os.Getenv; cwd is the current working directory.
// Missing files are skipped without error (returns nil for that layer).
// Returns (globalLayer, workspaceLayer, error). Error is only for read/parse failure.
func LoadGlobalAndWorkspace(getenv func(string) string, cwd string) (global, workspace *FileLayer, err error) {
	globalPath := GlobalPath(getenv)
	global, err = ReadLayer(globalPath)
	if err != nil {
		return nil, nil, err
	}
	workspacePath := WorkspacePath(cwd)
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
