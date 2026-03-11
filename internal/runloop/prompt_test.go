package runloop

import (
	"errors"
	"os"
	"path/filepath"
	"testing"

	"github.com/maxdunn/ralph/internal/config"
	"github.com/maxdunn/ralph/internal/review"
)

func TestLoadPromptOnce_noSource(t *testing.T) {
	provider := &review.FileLayerProvider{}
	cwd := t.TempDir()
	_, err := LoadPromptOnce(provider, cwd, review.ResolveOptions{})
	if err == nil {
		t.Fatal("LoadPromptOnce(no source) err = nil, want ErrMissingSource")
	}
	if !errors.Is(err, review.ErrMissingSource) {
		t.Errorf("err = %v, want ErrMissingSource", err)
	}
	if !review.IsExit2(err) {
		t.Error("IsExit2(err) = false, want true")
	}
}

func TestLoadPromptOnce_multipleSources(t *testing.T) {
	provider := &review.FileLayerProvider{}
	cwd := t.TempDir()
	_, err := LoadPromptOnce(provider, cwd, review.ResolveOptions{
		Alias:    "default",
		FilePath: "prompt.md",
	})
	if err == nil {
		t.Fatal("LoadPromptOnce(alias+file) err = nil, want ErrMultipleSources")
	}
	if !errors.Is(err, review.ErrMultipleSources) {
		t.Errorf("err = %v, want ErrMultipleSources", err)
	}
}

func TestLoadPromptOnce_stdin(t *testing.T) {
	provider := &review.FileLayerProvider{}
	cwd := t.TempDir()
	content := []byte("# My prompt\nDo the thing.")
	out, err := LoadPromptOnce(provider, cwd, review.ResolveOptions{Stdin: content})
	if err != nil {
		t.Fatalf("LoadPromptOnce(stdin) err = %v", err)
	}
	if string(out) != string(content) {
		t.Errorf("got %q, want %q", out, content)
	}
}

func TestLoadPromptOnce_fileNotFound(t *testing.T) {
	provider := &review.FileLayerProvider{}
	cwd := t.TempDir()
	_, err := LoadPromptOnce(provider, cwd, review.ResolveOptions{FilePath: "nonexistent.md"})
	if err == nil {
		t.Fatal("LoadPromptOnce(missing file) err = nil, want ErrFileNotFound")
	}
	if !errors.Is(err, review.ErrFileNotFound) {
		t.Errorf("err = %v, want ErrFileNotFound", err)
	}
}

func TestLoadPromptOnce_fileExists(t *testing.T) {
	provider := &review.FileLayerProvider{}
	dir := t.TempDir()
	p := filepath.Join(dir, "prompt.md")
	content := []byte("# File prompt")
	if err := os.WriteFile(p, content, 0644); err != nil {
		t.Fatal(err)
	}
	out, err := LoadPromptOnce(provider, dir, review.ResolveOptions{FilePath: "prompt.md"})
	if err != nil {
		t.Fatalf("LoadPromptOnce(file) err = %v", err)
	}
	if string(out) != string(content) {
		t.Errorf("got %q, want %q", out, content)
	}
}

func TestLoadPromptOnce_invalidAlias(t *testing.T) {
	provider := &review.FileLayerProvider{}
	cwd := t.TempDir()
	_, err := LoadPromptOnce(provider, cwd, review.ResolveOptions{Alias: "nonexistent"})
	if err == nil {
		t.Fatal("LoadPromptOnce(unknown alias) err = nil, want ErrInvalidAlias")
	}
	if !errors.Is(err, review.ErrInvalidAlias) {
		t.Errorf("err = %v, want ErrInvalidAlias", err)
	}
}

func TestLoadPromptOnce_aliasWithContent(t *testing.T) {
	layer, err := config.ParseLayer([]byte(`
prompts:
  inline:
    content: "# Inline prompt"
`))
	if err != nil || layer == nil {
		t.Fatalf("ParseLayer failed: %v", err)
	}
	provider := &review.FileLayerProvider{Layer: layer}
	cwd := t.TempDir()
	out, err := LoadPromptOnce(provider, cwd, review.ResolveOptions{Alias: "inline"})
	if err != nil {
		t.Fatalf("LoadPromptOnce(alias content) err = %v", err)
	}
	if string(out) != "# Inline prompt" {
		t.Errorf("got %q, want %q", out, "# Inline prompt")
	}
}

func TestLoadPromptOnce_aliasWithPath(t *testing.T) {
	dir := t.TempDir()
	p := filepath.Join(dir, "main.md")
	content := []byte("# Alias prompt")
	if err := os.WriteFile(p, content, 0644); err != nil {
		t.Fatal(err)
	}
	layer, err := config.ParseLayer([]byte("prompts:\n  default:\n    path: main.md"))
	if err != nil {
		t.Fatal(err)
	}
	if layer == nil {
		t.Fatal("ParseLayer failed")
	}
	provider := &review.FileLayerProvider{Layer: layer}
	out, err := LoadPromptOnce(provider, dir, review.ResolveOptions{Alias: "default"})
	if err != nil {
		t.Fatalf("LoadPromptOnce(alias path) err = %v", err)
	}
	if string(out) != string(content) {
		t.Errorf("got %q, want %q", out, content)
	}
}

func TestLoadPromptOnce_aliasPathMissing(t *testing.T) {
	dir := t.TempDir()
	layer, err := config.ParseLayer([]byte("prompts:\n  default:\n    path: missing.md"))
	if err != nil {
		t.Fatal(err)
	}
	if layer == nil {
		t.Fatal("ParseLayer failed")
	}
	provider := &review.FileLayerProvider{Layer: layer}
	_, err = LoadPromptOnce(provider, dir, review.ResolveOptions{Alias: "default"})
	if err == nil {
		t.Fatal("LoadPromptOnce(alias missing file) err = nil, want ErrAliasSourceMissing")
	}
	if !errors.Is(err, review.ErrAliasSourceMissing) {
		t.Errorf("err = %v, want ErrAliasSourceMissing", err)
	}
}

func TestAssemblePrompt_noPreamble(t *testing.T) {
	content := []byte("# Prompt\nDo the thing.")
	got := AssemblePrompt("", content)
	if string(got) != string(content) {
		t.Errorf("AssemblePrompt(empty preamble) = %q, want %q", got, content)
	}
	// Must return same slice when no preamble (no copy)
	if len(content) > 0 && &got[0] != &content[0] {
		t.Error("AssemblePrompt(empty preamble) should return bufferedContent as-is, not a copy")
	}
}

func TestAssemblePrompt_withPreamble(t *testing.T) {
	preamble := "CONTEXT:\niteration 1 of 5"
	content := []byte("# Prompt\nDo the thing.")
	got := AssemblePrompt(preamble, content)
	want := preamble + "\n" + string(content)
	if string(got) != want {
		t.Errorf("AssemblePrompt(%q, content) = %q, want %q", preamble, got, want)
	}
}

func TestAssemblePrompt_preambleWithNewline(t *testing.T) {
	// Preamble "line" + separator "\n" + content "body" => "line\nbody"
	content := []byte("body")
	got := AssemblePrompt("line", content)
	want := "line\nbody"
	if string(got) != want {
		t.Errorf("AssemblePrompt(\"line\", body) = %q, want %q", got, want)
	}
}

func TestAssemblePrompt_emptyContent(t *testing.T) {
	got := AssemblePrompt("preamble", nil)
	want := "preamble\n"
	if string(got) != want {
		t.Errorf("AssemblePrompt(preamble, nil) = %q, want %q", got, want)
	}
	got2 := AssemblePrompt("preamble", []byte{})
	if string(got2) != want {
		t.Errorf("AssemblePrompt(preamble, []) = %q, want %q", got2, want)
	}
}
