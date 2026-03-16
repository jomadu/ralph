package review

import (
	"github.com/jomadu/ralph/internal/config"
)

// FileLayerProvider adapts a single config file layer to PromptProvider.
// Used when effective config is not yet available (T1.7); later, effective config
// can implement PromptProvider directly and merge multiple layers.
// When ConfigPath is set, prompt file paths are resolved relative to that config file's directory.
type FileLayerProvider struct {
	Layer      *config.FileLayer
	ConfigPath string // optional; used to resolve relative prompt paths
}

// PromptByName implements PromptProvider using the layer's Prompts map.
func (p *FileLayerProvider) PromptByName(name string) (path, content string, ok bool) {
	if p == nil || p.Layer == nil || p.Layer.Prompts == nil {
		return "", "", false
	}
	prompt, ok := p.Layer.Prompts[name]
	if !ok {
		return "", "", false
	}
	path = prompt.Path
	if path != "" && p.ConfigPath != "" {
		path = config.ResolvePromptPath(p.ConfigPath, path)
	}
	return path, prompt.Content, true
}
