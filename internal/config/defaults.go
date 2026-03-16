// Package config: built-in defaults so the tool works without a config file (O002/R001, R002).
// See docs/engineering/components/config.md.

package config

// Default success/failure signal strings (match build procedure and O001 signal detection).
const (
	DefaultSuccessSignal = "<promise>SUCCESS</promise>"
	DefaultFailureSignal = "<promise>FAILURE</promise>"
)

// DefaultLoopSettings returns built-in loop settings. Used when no config file
// is present or when a setting is omitted from all layers.
func DefaultLoopSettings() LoopSettings {
	return LoopSettings{
		MaxIterations:    10,
		FailureThreshold: 3,
		TimeoutSeconds:   0, // no per-iteration timeout
		SuccessSignal:    DefaultSuccessSignal,
		FailureSignal:    DefaultFailureSignal,
		SignalPrecedence: "static",
		Preamble:         "",
		Streaming:        true,
		LogLevel:         "info",
		AICmd:            "",
		AICmdAlias:       "",
	}
}

// DefaultEffective returns an effective config with default loop settings and
// built-in aliases only (no prompts). Ensures the tool works without a config file.
func DefaultEffective() *Effective {
	e := &Effective{
		Loop:    DefaultLoopSettings(),
		Prompts: make(map[string]Prompt),
		Aliases: make(map[string]Alias),
	}
	for k, v := range BuiltinAliases() {
		e.Aliases[k] = v
	}
	return e
}
