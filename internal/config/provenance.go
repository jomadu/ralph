// Package config: provenance for effective config (which layer supplied each value).
// Used by ralph show config --provenance (O002/R007, T6.3).

package config

import "path/filepath"

// Layer names for provenance (cli.md: default, global, workspace, explicit file, env, cli, prompt).
const (
	ProvenanceDefault   = "default"
	ProvenanceGlobal    = "global"
	ProvenanceWorkspace = "workspace"
	ProvenanceExplicit  = "explicit"
	ProvenanceEnv       = "env"
	ProvenanceCLI       = "cli"
	ProvenancePrompt    = "prompt"
)

// LoopProvenance records which layer supplied each loop setting (for show config --provenance).
type LoopProvenance struct {
	MaxIterations    string
	FailureThreshold string
	TimeoutSeconds   string
	SuccessSignal    string
	FailureSignal    string
	SignalPrecedence string
	Preamble         string
	Streaming        string
	LogLevel         string
}

// RootLoopWithProvenance computes root loop (defaults → global → workspace or explicit → env)
// and per-field provenance. When configPath is non-empty, only that file is used (explicit);
// otherwise global and workspace are used. Used by CLI show config --provenance.
func RootLoopWithProvenance(getenv func(string) string, cwd, configPath string) (LoopSettings, LoopProvenance, error) {
	provenance := LoopProvenance{
		MaxIterations:    ProvenanceDefault,
		FailureThreshold: ProvenanceDefault,
		TimeoutSeconds:   ProvenanceDefault,
		SuccessSignal:    ProvenanceDefault,
		FailureSignal:    ProvenanceDefault,
		SignalPrecedence: ProvenanceDefault,
		Preamble:         ProvenanceDefault,
		Streaming:        ProvenanceDefault,
		LogLevel:         ProvenanceDefault,
	}
	loop := DefaultLoopSettings()

	if configPath != "" {
		path := configPath
		if !filepath.IsAbs(path) {
			path = filepath.Join(cwd, path)
		}
		layer, err := LoadExplicit(path)
		if err != nil {
			return LoopSettings{}, LoopProvenance{}, err
		}
		loop = applySectionWithProvenance(loop, layer.Loop, ProvenanceExplicit, &provenance)
	} else {
		global, workspace, err := LoadGlobalAndWorkspace(getenv, cwd)
		if err != nil {
			return LoopSettings{}, LoopProvenance{}, err
		}
		if global != nil && global.Loop != nil {
			loop = applySectionWithProvenance(loop, global.Loop, ProvenanceGlobal, &provenance)
		}
		if workspace != nil && workspace.Loop != nil {
			loop = applySectionWithProvenance(loop, workspace.Loop, ProvenanceWorkspace, &provenance)
		}
	}

	overlay, err := ParseEnvOverlay(getenv)
	if err != nil {
		return LoopSettings{}, LoopProvenance{}, err
	}
	loop, provenance = applyEnvOverlayWithProvenance(loop, overlay, provenance)
	return loop, provenance, nil
}

// applySectionWithProvenance applies a file layer's loop section and updates provenance for overridden fields.
func applySectionWithProvenance(base LoopSettings, section *LoopSection, layer string, prov *LoopProvenance) LoopSettings {
	if section == nil {
		return base
	}
	out := base
	if section.MaxIterations != nil {
		out.MaxIterations = *section.MaxIterations
		prov.MaxIterations = layer
	}
	if section.FailureThreshold != nil {
		out.FailureThreshold = *section.FailureThreshold
		prov.FailureThreshold = layer
	}
	if section.TimeoutSeconds != nil {
		out.TimeoutSeconds = *section.TimeoutSeconds
		prov.TimeoutSeconds = layer
	}
	if section.SuccessSignal != "" {
		out.SuccessSignal = section.SuccessSignal
		prov.SuccessSignal = layer
	}
	if section.FailureSignal != "" {
		out.FailureSignal = section.FailureSignal
		prov.FailureSignal = layer
	}
	if section.SignalPrecedence != "" {
		out.SignalPrecedence = section.SignalPrecedence
		prov.SignalPrecedence = layer
	}
	if section.Streaming != nil {
		out.Streaming = *section.Streaming
		prov.Streaming = layer
	}
	if section.LogLevel != "" {
		out.LogLevel = section.LogLevel
		prov.LogLevel = layer
	}
	if s, ok := section.Preamble.(string); ok && s != "" {
		out.Preamble = s
		prov.Preamble = layer
	}
	if b, ok := section.Preamble.(bool); ok && !b {
		out.Preamble = ""
		prov.Preamble = layer
	}
	return out
}

// applyEnvOverlayWithProvenance applies env overlay and updates provenance for overridden fields.
func applyEnvOverlayWithProvenance(loop LoopSettings, overlay *EnvOverlay, prov LoopProvenance) (LoopSettings, LoopProvenance) {
	if overlay == nil {
		return loop, prov
	}
	out := loop
	if overlay.MaxIterations != nil {
		out.MaxIterations = *overlay.MaxIterations
		prov.MaxIterations = ProvenanceEnv
	}
	if overlay.FailureThreshold != nil {
		out.FailureThreshold = *overlay.FailureThreshold
		prov.FailureThreshold = ProvenanceEnv
	}
	if overlay.IterationTimeout != nil {
		out.TimeoutSeconds = *overlay.IterationTimeout
		prov.TimeoutSeconds = ProvenanceEnv
	}
	if overlay.LogLevel != nil {
		out.LogLevel = *overlay.LogLevel
		prov.LogLevel = ProvenanceEnv
	}
	if overlay.Streaming != nil {
		out.Streaming = *overlay.Streaming
		prov.Streaming = ProvenanceEnv
	}
	if overlay.Preamble != nil && !*overlay.Preamble {
		out.Preamble = ""
		prov.Preamble = ProvenanceEnv
	}
	return out, prov
}
