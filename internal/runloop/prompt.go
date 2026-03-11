package runloop

import (
	"github.com/maxdunn/ralph/internal/review"
)

// LoadPromptOnce loads the prompt once from the resolved source (alias → file path,
// file path, or stdin) and returns the buffered content for use by all iterations.
// Implements O001/R002: load once, buffer in memory, fail before loop if source
// is unavailable. The caller must not re-read the source between iterations.
//
// provider supplies prompt definitions for alias resolution; cwd is used to
// resolve relative paths. opts must specify exactly one of Alias, FilePath, or
// Stdin (Stdin is the already-read bytes when prompt source is stdin).
func LoadPromptOnce(provider review.PromptProvider, cwd string, opts review.ResolveOptions) ([]byte, error) {
	return review.ResolvePromptSource(provider, cwd, opts)
}
