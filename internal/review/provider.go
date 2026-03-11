package review

import (
	"github.com/maxdunn/ralph/internal/config"
)

// FileLayerProvider adapts a single config file layer to PromptProvider.
// Used when effective config is not yet available (T1.7); later, effective config
// can implement PromptProvider directly and merge multiple layers.
type FileLayerProvider struct {
	Layer *config.FileLayer
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
	return prompt.Path, prompt.Content, true
}
