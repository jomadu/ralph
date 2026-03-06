package config

import (
	"testing"
)

func TestResolveAICommand(t *testing.T) {
	tests := []struct {
		name        string
		cfg         ConfigWithProvenance
		wantCommand string
		wantSource  string
		wantErr     bool
	}{
		{
			name: "direct command takes precedence",
			cfg: ConfigWithProvenance{
				Loop: LoopConfigWithProvenance{
					AICmd:      ValueWithProvenance[string]{Value: "my-cli --flag", Provenance: ProvenanceCLI},
					AICmdAlias: ValueWithProvenance[string]{Value: "claude", Provenance: ProvenanceCLI},
				},
				AICmdAliases: map[string]ValueWithProvenance[string]{
					"claude": {Value: "claude -p --dangerously-skip-permissions", Provenance: ProvenanceDefault},
				},
			},
			wantCommand: "my-cli --flag",
			wantSource:  "direct command",
			wantErr:     false,
		},
		{
			name: "alias resolution when no direct command",
			cfg: ConfigWithProvenance{
				Loop: LoopConfigWithProvenance{
					AICmd:      ValueWithProvenance[string]{Value: "", Provenance: ProvenanceDefault},
					AICmdAlias: ValueWithProvenance[string]{Value: "kiro", Provenance: ProvenanceEnv},
				},
				AICmdAliases: map[string]ValueWithProvenance[string]{
					"kiro": {Value: "kiro-cli chat --no-interactive", Provenance: ProvenanceDefault},
				},
			},
			wantCommand: "kiro-cli chat --no-interactive",
			wantSource:  "alias kiro",
			wantErr:     false,
		},
		{
			name: "no command configured",
			cfg: ConfigWithProvenance{
				Loop: LoopConfigWithProvenance{
					AICmd:      ValueWithProvenance[string]{Value: "", Provenance: ProvenanceDefault},
					AICmdAlias: ValueWithProvenance[string]{Value: "", Provenance: ProvenanceDefault},
				},
				AICmdAliases: map[string]ValueWithProvenance[string]{},
			},
			wantErr: true,
		},
		{
			name: "unknown alias",
			cfg: ConfigWithProvenance{
				Loop: LoopConfigWithProvenance{
					AICmd:      ValueWithProvenance[string]{Value: "", Provenance: ProvenanceDefault},
					AICmdAlias: ValueWithProvenance[string]{Value: "unknown", Provenance: ProvenanceEnv},
				},
				AICmdAliases: map[string]ValueWithProvenance[string]{},
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			res, err := ResolveAICommand(tt.cfg)
			if (err != nil) != tt.wantErr {
				t.Errorf("ResolveAICommand() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				if res.Command != tt.wantCommand {
					t.Errorf("ResolveAICommand() command = %v, want %v", res.Command, tt.wantCommand)
				}
				if res.Source != tt.wantSource {
					t.Errorf("ResolveAICommand() source = %v, want %v", res.Source, tt.wantSource)
				}
			}
		})
	}
}

func TestResolveEffectiveConfigForPrompt(t *testing.T) {
	tests := []struct {
		name        string
		cfg         ConfigWithProvenance
		promptAlias string
		wantValue   int
		wantProv    Provenance
		wantErr     bool
	}{
		{
			name: "prompt override takes precedence over workspace",
			cfg: ConfigWithProvenance{
				Loop: LoopConfigWithProvenance{
					FailureThreshold: ValueWithProvenance[int]{Value: 3, Provenance: ProvenanceWorkspace},
				},
				Prompts: map[string]PromptConfigWithProvenance{
					"build": {
						Loop: &LoopConfig{
							FailureThreshold: 5,
						},
						LoopRawMap: map[string]interface{}{
							"failure_threshold": 5,
						},
						Provenance: ProvenanceWorkspace,
					},
				},
			},
			promptAlias: "build",
			wantValue:   5,
			wantProv:    ProvenancePrompt,
			wantErr:     false,
		},
		{
			name: "prompt override does not affect other fields",
			cfg: ConfigWithProvenance{
				Loop: LoopConfigWithProvenance{
					FailureThreshold:     ValueWithProvenance[int]{Value: 3, Provenance: ProvenanceWorkspace},
					DefaultMaxIterations: ValueWithProvenance[int]{Value: 10, Provenance: ProvenanceWorkspace},
				},
				Prompts: map[string]PromptConfigWithProvenance{
					"build": {
						Loop: &LoopConfig{
							FailureThreshold: 5,
						},
						LoopRawMap: map[string]interface{}{
							"failure_threshold": 5,
						},
						Provenance: ProvenanceWorkspace,
					},
				},
			},
			promptAlias: "build",
			wantValue:   10,
			wantProv:    ProvenanceWorkspace,
			wantErr:     false,
		},
		{
			name: "no prompt overrides returns base config",
			cfg: ConfigWithProvenance{
				Loop: LoopConfigWithProvenance{
					FailureThreshold: ValueWithProvenance[int]{Value: 3, Provenance: ProvenanceWorkspace},
				},
				Prompts: map[string]PromptConfigWithProvenance{
					"build": {
						Loop:       nil,
						Provenance: ProvenanceWorkspace,
					},
				},
			},
			promptAlias: "build",
			wantValue:   3,
			wantProv:    ProvenanceWorkspace,
			wantErr:     false,
		},
		{
			name: "unknown prompt alias",
			cfg: ConfigWithProvenance{
				Loop: LoopConfigWithProvenance{
					FailureThreshold: ValueWithProvenance[int]{Value: 3, Provenance: ProvenanceWorkspace},
				},
				Prompts: map[string]PromptConfigWithProvenance{},
			},
			promptAlias: "unknown",
			wantErr:     true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			effective, err := ResolveEffectiveConfigForPrompt(tt.cfg, tt.promptAlias)
			if (err != nil) != tt.wantErr {
				t.Errorf("ResolveEffectiveConfigForPrompt() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				// Check failure_threshold for first two tests
				if tt.name == "prompt override takes precedence over workspace" {
					if effective.Loop.FailureThreshold.Value != tt.wantValue {
						t.Errorf("FailureThreshold value = %v, want %v", effective.Loop.FailureThreshold.Value, tt.wantValue)
					}
					if effective.Loop.FailureThreshold.Provenance != tt.wantProv {
						t.Errorf("FailureThreshold provenance = %v, want %v", effective.Loop.FailureThreshold.Provenance, tt.wantProv)
					}
				}
				// Check default_max_iterations for second test
				if tt.name == "prompt override does not affect other fields" {
					if effective.Loop.DefaultMaxIterations.Value != tt.wantValue {
						t.Errorf("DefaultMaxIterations value = %v, want %v", effective.Loop.DefaultMaxIterations.Value, tt.wantValue)
					}
					if effective.Loop.DefaultMaxIterations.Provenance != tt.wantProv {
						t.Errorf("DefaultMaxIterations provenance = %v, want %v", effective.Loop.DefaultMaxIterations.Provenance, tt.wantProv)
					}
				}
				// Check failure_threshold for third test
				if tt.name == "no prompt overrides returns base config" {
					if effective.Loop.FailureThreshold.Value != tt.wantValue {
						t.Errorf("FailureThreshold value = %v, want %v", effective.Loop.FailureThreshold.Value, tt.wantValue)
					}
					if effective.Loop.FailureThreshold.Provenance != tt.wantProv {
						t.Errorf("FailureThreshold provenance = %v, want %v", effective.Loop.FailureThreshold.Provenance, tt.wantProv)
					}
				}
			}
		})
	}
}
