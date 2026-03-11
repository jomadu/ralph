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

// AssemblePrompt builds the assembled prompt sent to the AI: optional preamble
// (e.g. iteration count, context) followed by the buffered prompt content.
// Preamble is configurable via config/CLI; no preamble or prompt content is
// written to any user file (O009/R001). Implements T3.3, O002/R002.
//
// preamble may be empty (no injection). When non-empty, it is prepended to
// bufferedContent with a single newline separator so the AI receives one
// coherent prompt. Callers may build preamble from config and optionally
// inject iteration number or other per-iteration context.
func AssemblePrompt(preamble string, bufferedContent []byte) []byte {
	if preamble == "" {
		return bufferedContent
	}
	// Prepend preamble + newline + content. No trailing newline added to preamble
	// so we control exactly one separator.
	sep := []byte("\n")
	out := make([]byte, 0, len(preamble)+len(sep)+len(bufferedContent))
	out = append(out, preamble...)
	out = append(out, sep...)
	out = append(out, bufferedContent...)
	return out
}
