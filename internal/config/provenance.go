// Package config: provenance for effective config (which layer supplied each value).
// Used by ralph show config --provenance (O002/R007, T6.3, T7.3). Layers: default, global, workspace, explicit, env, prompt; "cli" only when a CLIOverlay is passed (e.g. from run); show config passes nil.

package config

import (
	"path/filepath"
	"strings"
)

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
	AICmd            string
	AICmdAlias       string
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
		AICmd:            ProvenanceDefault,
		AICmdAlias:       ProvenanceDefault,
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
		globalPath := GlobalPath(getenv)
		workspacePath := WorkspacePath(cwd)
		global, workspace, err := LoadGlobalAndWorkspace(globalPath, workspacePath)
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
	if section.AiCmd != "" {
		out.AICmd = section.AiCmd
		prov.AICmd = layer
	}
	if section.AiCmdAlias != "" {
		out.AICmdAlias = section.AiCmdAlias
		prov.AICmdAlias = layer
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
	if overlay.AICmd != nil {
		out.AICmd = *overlay.AICmd
		prov.AICmd = ProvenanceEnv
	}
	if overlay.AICmdAlias != nil {
		out.AICmdAlias = *overlay.AICmdAlias
		prov.AICmdAlias = ProvenanceEnv
	}
	return out, prov
}

// CLIOverlay holds optional CLI flag overrides for loop settings (T7.3, O002/R007).
// Used when computing effective loop with run-style overrides (e.g. from ralph run). The show config
// command does not accept run-style flags and calls LoopWithProvenance with nil.
// Semantics match run: MaxIterations 0 = use config; FailureThreshold/IterationTimeout -1 = not set.
// Unlimited sets max iterations to a large cap and provenance to cli.
type CLIOverlay struct {
	MaxIterations    int  // 0 = use config
	Unlimited        bool // overrides MaxIterations when true
	FailureThreshold int  // -1 = not set
	IterationTimeout int  // -1 = not set
	NoPreamble       bool
	SignalSuccess    string
	SignalFailure    string
	SignalPrecedence string
	Context          []string
	Verbose          bool
	Quiet            bool
	LogLevel         string
	NoStream         bool
}

// applyCLIOverlayWithProvenance applies CLI overlay and sets provenance to ProvenanceCLI for overridden fields.
func applyCLIOverlayWithProvenance(loop LoopSettings, o *CLIOverlay, prov LoopProvenance) (LoopSettings, LoopProvenance) {
	if o == nil {
		return loop, prov
	}
	out := loop
	if o.Unlimited {
		const unlimitedCap = 1<<31 - 1
		out.MaxIterations = unlimitedCap
		prov.MaxIterations = ProvenanceCLI
	} else if o.MaxIterations > 0 {
		out.MaxIterations = o.MaxIterations
		prov.MaxIterations = ProvenanceCLI
	}
	if o.FailureThreshold >= 0 {
		out.FailureThreshold = o.FailureThreshold
		prov.FailureThreshold = ProvenanceCLI
	}
	if o.IterationTimeout >= 0 {
		out.TimeoutSeconds = o.IterationTimeout
		prov.TimeoutSeconds = ProvenanceCLI
	}
	if o.SignalSuccess != "" {
		out.SuccessSignal = o.SignalSuccess
		prov.SuccessSignal = ProvenanceCLI
	}
	if o.SignalFailure != "" {
		out.FailureSignal = o.SignalFailure
		prov.FailureSignal = ProvenanceCLI
	}
	if o.SignalPrecedence != "" {
		out.SignalPrecedence = o.SignalPrecedence
		prov.SignalPrecedence = ProvenanceCLI
	}
	if o.NoPreamble {
		out.Preamble = ""
		prov.Preamble = ProvenanceCLI
	} else if len(o.Context) > 0 {
		contextBlock := "CONTEXT\n" + strings.Join(o.Context, "\n")
		if out.Preamble != "" {
			out.Preamble = out.Preamble + "\n" + contextBlock
		} else {
			out.Preamble = contextBlock
		}
		prov.Preamble = ProvenanceCLI
	}
	if o.Quiet && !o.Verbose {
		out.LogLevel = "error"
		out.Streaming = false
		prov.LogLevel = ProvenanceCLI
		prov.Streaming = ProvenanceCLI
	}
	if o.Verbose {
		out.LogLevel = "debug"
		out.Streaming = true
		prov.LogLevel = ProvenanceCLI
		prov.Streaming = ProvenanceCLI
	}
	if o.LogLevel != "" {
		out.LogLevel = o.LogLevel
		prov.LogLevel = ProvenanceCLI
	}
	if o.NoStream {
		out.Streaming = false
		prov.Streaming = ProvenanceCLI
	}
	return out, prov
}

// LoopWithProvenance returns the effective loop and provenance for the given context, including
// optional prompt overrides and CLI overlay (T7.3, O002/R007). Order: root (defaults → layers → env)
// then prompt overrides (if promptName is set and the prompt has loop overrides), then CLI overlay.
// When promptName is empty or the prompt has no loop section, prompt layer is skipped.
// When cli is nil, CLI layer is skipped.
func LoopWithProvenance(getenv func(string) string, cwd, configPath, promptName string, cli *CLIOverlay) (LoopSettings, LoopProvenance, error) {
	rootLoop, prov, err := RootLoopWithProvenance(getenv, cwd, configPath)
	if err != nil {
		return LoopSettings{}, LoopProvenance{}, err
	}
	loop := rootLoop

	if promptName != "" {
		resolved, _, err := loadLayersAndRootLoop(getenv, cwd, configPath)
		if err != nil {
			return LoopSettings{}, LoopProvenance{}, err
		}
		prompt, ok := resolved.Prompts[promptName]
		if ok && prompt.Loop != nil {
			loop = applySectionWithProvenance(loop, prompt.Loop, ProvenancePrompt, &prov)
		}
	}

	loop, prov = applyCLIOverlayWithProvenance(loop, cli, prov)
	return loop, prov, nil
}
