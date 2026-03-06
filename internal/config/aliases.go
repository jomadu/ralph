package config

// BuiltinAliases returns the built-in AI command aliases per O3/R1.
func BuiltinAliases() map[string]string {
	return map[string]string{
		"claude":       "claude -p --dangerously-skip-permissions",
		"kiro":         "kiro-cli chat --no-interactive --trust-all-tools",
		"copilot":      "copilot --yolo",
		"cursor-agent": "scripts/cursor-wrapper.sh",
	}
}

// MergedAliases returns built-in aliases merged with user-defined aliases.
// User-defined aliases override built-in aliases with the same name per O3/R3.
func MergedAliases(cfg Config) map[string]string {
	merged := make(map[string]string)
	
	// Start with built-ins
	for k, v := range BuiltinAliases() {
		merged[k] = v
	}
	
	// Overlay user-defined (overrides built-in for same key)
	for k, v := range cfg.AICmdAliases {
		merged[k] = v
	}
	
	return merged
}
