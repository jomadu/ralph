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
