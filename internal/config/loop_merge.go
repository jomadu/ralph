// Package config: merge loop settings per layer order (O002/R003).
// Order: defaults → global → workspace → explicit file → env → prompt overrides → CLI (CLI in Phase 4).

package config

// ApplyLoopSection applies a file layer's loop section on top of base settings.
// Only non-nil (pointers) or non-empty (strings) fields in section override base.
func ApplyLoopSection(base LoopSettings, section *LoopSection) LoopSettings {
	if section == nil {
		return base
	}
	out := base
	if section.MaxIterations != nil {
		out.MaxIterations = *section.MaxIterations
	}
	if section.FailureThreshold != nil {
		out.FailureThreshold = *section.FailureThreshold
	}
	if section.TimeoutSeconds != nil {
		out.TimeoutSeconds = *section.TimeoutSeconds
	}
	if section.SuccessSignal != "" {
		out.SuccessSignal = section.SuccessSignal
	}
	if section.FailureSignal != "" {
		out.FailureSignal = section.FailureSignal
	}
	if section.SignalPrecedence != "" {
		out.SignalPrecedence = section.SignalPrecedence
	}
	if section.Streaming != nil {
		out.Streaming = *section.Streaming
	}
	if section.LogLevel != "" {
		out.LogLevel = section.LogLevel
	}
	// Preamble: string or bool in YAML; we store string. If section has a string, use it.
	if s, ok := section.Preamble.(string); ok && s != "" {
		out.Preamble = s
	}
	if b, ok := section.Preamble.(bool); ok && !b {
		out.Preamble = ""
	}
	if section.AiCmd != "" {
		out.AICmd = section.AiCmd
	}
	if section.AiCmdAlias != "" {
		out.AICmdAlias = section.AiCmdAlias
	}
	return out
}

// MergeRootLoop merges root loop from defaults and file layers.
// Order: defaults → global → workspace. Use when building effective config from global+workspace
// or from a single explicit layer (pass that layer as workspace and nil global).
func MergeRootLoop(global, workspace *FileLayer) LoopSettings {
	loop := DefaultLoopSettings()
	if global != nil && global.Loop != nil {
		loop = ApplyLoopSection(loop, global.Loop)
	}
	if workspace != nil && workspace.Loop != nil {
		loop = ApplyLoopSection(loop, workspace.Loop)
	}
	return loop
}

// ApplyEnvOverlayToLoop applies RALPH_LOOP_* env overlay onto loop settings.
// Only non-nil overlay fields override (O010/R004).
func ApplyEnvOverlayToLoop(loop LoopSettings, overlay *EnvOverlay) LoopSettings {
	if overlay == nil {
		return loop
	}
	out := loop
	if overlay.MaxIterations != nil {
		out.MaxIterations = *overlay.MaxIterations
	}
	if overlay.FailureThreshold != nil {
		out.FailureThreshold = *overlay.FailureThreshold
	}
	if overlay.IterationTimeout != nil {
		out.TimeoutSeconds = *overlay.IterationTimeout
	}
	if overlay.LogLevel != nil {
		out.LogLevel = *overlay.LogLevel
	}
	if overlay.Streaming != nil {
		out.Streaming = *overlay.Streaming
	}
	if overlay.Preamble != nil {
		if !*overlay.Preamble {
			out.Preamble = ""
		}
	}
	if overlay.AICmd != nil {
		out.AICmd = *overlay.AICmd
	}
	if overlay.AICmdAlias != nil {
		out.AICmdAlias = *overlay.AICmdAlias
	}
	return out
}

// EffectiveLoopForPrompt returns the effective loop settings for a given prompt.
// Order: rootLoop (already defaults → layers → env) then prompt-level loop overrides.
// If prompt is nil or has no Loop, returns rootLoop unchanged.
func EffectiveLoopForPrompt(rootLoop LoopSettings, prompt *Prompt) LoopSettings {
	if prompt == nil || prompt.Loop == nil {
		return rootLoop
	}
	return ApplyLoopSection(rootLoop, prompt.Loop)
}
