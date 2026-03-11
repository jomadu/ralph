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
