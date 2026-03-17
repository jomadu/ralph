// Package config: provenance for effective config (which layer supplied each value).
// Used by ralph show config --provenance (O002/R007, T6.3, T7.3). Layers: default, global, workspace, explicit, env, prompt; "cli" only when a CLIOverlay is passed (e.g. from run); show config passes nil.

package config

import "strings"

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

// RootLoopInput holds all inputs needed to compute root loop and provenance. The caller
// constructs this (e.g. via NewRootLoopInput) so that RootLoopWithProvenance does no I/O
// or env reads—making it a pure merge. When Explicit is non-nil, only Explicit and
// EnvOverlay are used; otherwise Global and Workspace are merged, then EnvOverlay applied.
type RootLoopInput struct {
	Explicit   *FileLayer  // when set, only this layer is used (no global/workspace)
	Global     *FileLayer  // used when Explicit is nil
	Workspace  *FileLayer  // used when Explicit is nil
	EnvOverlay *EnvOverlay // applied after file layers; may be nil
}

// LoopProvenance records which layer supplied each loop setting (for show config --provenance).
type LoopProvenance struct {
	MaxIterations    string
	FailureThreshold string
	TimeoutSeconds   string
	SuccessSignal    string
	FailureSignal    string
	Preamble         string
	Context          string
	Streaming        string
	LogLevel         string
	MaxOutputBuffer  string
	AICmd            string
	AICmdAlias       string
}

// RootLoopWithProvenance computes root loop and provenance from pre-loaded inputs. It does
// no I/O or env reads. Order: defaults → explicit (if set) or global → workspace → env overlay.
// Used by CLI show config --provenance. Call NewRootLoopInput to build the input from getenv/cwd/configPath.
func RootLoopWithProvenance(input RootLoopInput) (LoopSettings, LoopProvenance) {
	provenance := LoopProvenance{
		MaxIterations:    ProvenanceDefault,
		FailureThreshold: ProvenanceDefault,
		TimeoutSeconds:   ProvenanceDefault,
		SuccessSignal:    ProvenanceDefault,
		FailureSignal:    ProvenanceDefault,
		Preamble:         ProvenanceDefault,
		Context:          ProvenanceDefault,
		Streaming:        ProvenanceDefault,
		LogLevel:         ProvenanceDefault,
		MaxOutputBuffer:  ProvenanceDefault,
		AICmd:            ProvenanceDefault,
		AICmdAlias:       ProvenanceDefault,
	}
	loop := DefaultLoopSettings()

	if input.Explicit != nil {
		loop = applySectionWithProvenance(loop, input.Explicit.Loop, ProvenanceExplicit, &provenance)
	} else {
		if input.Global != nil && input.Global.Loop != nil {
			loop = applySectionWithProvenance(loop, input.Global.Loop, ProvenanceGlobal, &provenance)
		}
		if input.Workspace != nil && input.Workspace.Loop != nil {
			loop = applySectionWithProvenance(loop, input.Workspace.Loop, ProvenanceWorkspace, &provenance)
		}
	}

	loop, provenance = applyEnvOverlayWithProvenance(loop, input.EnvOverlay, provenance)
	return loop, provenance
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
	if section.Streaming != nil {
		out.Streaming = *section.Streaming
		prov.Streaming = layer
	}
	if section.LogLevel != "" {
		out.LogLevel = section.LogLevel
		prov.LogLevel = layer
	}
	if section.MaxOutputBuffer != nil {
		out.MaxOutputBuffer = *section.MaxOutputBuffer
		prov.MaxOutputBuffer = layer
	}
	if section.Preamble != nil {
		out.Preamble = *section.Preamble
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
	if overlay.Preamble != nil {
		out.Preamble = *overlay.Preamble
		prov.Preamble = ProvenanceEnv
	}
	if overlay.MaxOutputBuffer != nil {
		out.MaxOutputBuffer = *overlay.MaxOutputBuffer
		prov.MaxOutputBuffer = ProvenanceEnv
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
	if o.NoPreamble {
		out.Preamble = false
		prov.Preamble = ProvenanceCLI
	}
	if len(o.Context) > 0 {
		out.Context = strings.Join(o.Context, "\n")
		prov.Context = ProvenanceCLI
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

// LoopWithProvenanceInput holds all inputs for LoopWithProvenance. The caller builds this
// (Root from NewRootLoopInput; Resolved from the same call for prompt lookup). Pure options.
type LoopWithProvenanceInput struct {
	Root       RootLoopInput
	Resolved   *Resolved // used when PromptName is set to apply prompt loop overrides; may be nil
	PromptName string
	CLI        *CLIOverlay
}

// LoopWithProvenance returns the effective loop and provenance from the given options. It does
// no I/O or env reads. Order: root (from opts.Root) → prompt overrides (if opts.PromptName
// and opts.Resolved) → CLI overlay (if opts.CLI). Call NewRootLoopInput to build opts.Root and opts.Resolved.
func LoopWithProvenance(opts LoopWithProvenanceInput) (LoopSettings, LoopProvenance) {
	loop, prov := RootLoopWithProvenance(opts.Root)

	if opts.PromptName != "" && opts.Resolved != nil {
		if prompt, ok := opts.Resolved.Prompts[opts.PromptName]; ok && prompt.Loop != nil {
			loop = applySectionWithProvenance(loop, prompt.Loop, ProvenancePrompt, &prov)
		}
	}

	loop, prov = applyCLIOverlayWithProvenance(loop, opts.CLI, prov)
	return loop, prov
}
