package review

import (
	"errors"
	"os"
	"path/filepath"
	"testing"

	"github.com/maxdunn/ralph/internal/config"
)

func TestResolvePromptSource_noSource(t *testing.T) {
	provider := &FileLayerProvider{}
	cwd := t.TempDir()
	_, err := ResolvePromptSource(provider, cwd, ResolveOptions{})
	if err == nil {
		t.Fatal("ResolvePromptSource(no source) err = nil, want ErrMissingSource")
	}
	if !errors.Is(err, ErrMissingSource) {
		t.Errorf("err = %v, want ErrMissingSource", err)
	}
	if !IsExit2(err) {
		t.Error("IsExit2(err) = false, want true")
	}
}

func TestResolvePromptSource_multipleSources(t *testing.T) {
	provider := &FileLayerProvider{}
	cwd := t.TempDir()
	_, err := ResolvePromptSource(provider, cwd, ResolveOptions{
		Alias:    "default",
		FilePath: "prompt.md",
	})
	if err == nil {
		t.Fatal("ResolvePromptSource(alias+file) err = nil, want ErrMultipleSources")
	}
	if !errors.Is(err, ErrMultipleSources) {
		t.Errorf("err = %v, want ErrMultipleSources", err)
	}
	if !IsExit2(err) {
		t.Error("IsExit2(err) = false, want true")
	}
}

func TestResolvePromptSource_stdin(t *testing.T) {
	provider := &FileLayerProvider{}
	cwd := t.TempDir()
	content := []byte("# My prompt\nDo the thing.")
	out, err := ResolvePromptSource(provider, cwd, ResolveOptions{Stdin: content})
	if err != nil {
		t.Fatalf("ResolvePromptSource(stdin) err = %v", err)
	}
	if string(out) != string(content) {
		t.Errorf("got %q, want %q", out, content)
	}
}

func TestResolvePromptSource_fileNotFound(t *testing.T) {
	provider := &FileLayerProvider{}
	cwd := t.TempDir()
	_, err := ResolvePromptSource(provider, cwd, ResolveOptions{FilePath: "nonexistent.md"})
	if err == nil {
		t.Fatal("ResolvePromptSource(missing file) err = nil, want ErrFileNotFound")
	}
	if !errors.Is(err, ErrFileNotFound) {
		t.Errorf("err = %v, want ErrFileNotFound", err)
	}
	if !IsExit2(err) {
		t.Error("IsExit2(err) = false, want true")
	}
}

func TestResolvePromptSource_fileExists(t *testing.T) {
	provider := &FileLayerProvider{}
	dir := t.TempDir()
	p := filepath.Join(dir, "prompt.md")
	content := []byte("# File prompt")
	if err := os.WriteFile(p, content, 0644); err != nil {
		t.Fatal(err)
	}
	out, err := ResolvePromptSource(provider, dir, ResolveOptions{FilePath: "prompt.md"})
	if err != nil {
		t.Fatalf("ResolvePromptSource(file) err = %v", err)
	}
	if string(out) != string(content) {
		t.Errorf("got %q, want %q", out, content)
	}
}

func TestResolvePromptSource_invalidAlias(t *testing.T) {
	provider := &FileLayerProvider{}
	cwd := t.TempDir()
	_, err := ResolvePromptSource(provider, cwd, ResolveOptions{Alias: "nonexistent"})
	if err == nil {
		t.Fatal("ResolvePromptSource(unknown alias) err = nil, want ErrInvalidAlias")
	}
	if !errors.Is(err, ErrInvalidAlias) {
		t.Errorf("err = %v, want ErrInvalidAlias", err)
	}
	if !IsExit2(err) {
		t.Error("IsExit2(err) = false, want true")
	}
}

func TestResolvePromptSource_aliasWithContent(t *testing.T) {
	layer, err := config.ParseLayer([]byte(`
prompts:
  inline:
    content: "# Inline prompt"
`))
	if err != nil || layer == nil {
		t.Fatalf("ParseLayer failed: %v", err)
	}
	provider := &FileLayerProvider{Layer: layer}
	cwd := t.TempDir()
	out, err := ResolvePromptSource(provider, cwd, ResolveOptions{Alias: "inline"})
	if err != nil {
		t.Fatalf("ResolvePromptSource(alias content) err = %v", err)
	}
	if string(out) != "# Inline prompt" {
		t.Errorf("got %q, want %q", out, "# Inline prompt")
	}
}

func TestResolvePromptSource_aliasWithPath(t *testing.T) {
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
	provider := &FileLayerProvider{Layer: layer}
	out, err := ResolvePromptSource(provider, dir, ResolveOptions{Alias: "default"})
	if err != nil {
		t.Fatalf("ResolvePromptSource(alias path) err = %v", err)
	}
	if string(out) != string(content) {
		t.Errorf("got %q, want %q", out, content)
	}
}

func TestResolvePromptSource_aliasPathMissing(t *testing.T) {
	dir := t.TempDir()
	layer, err := config.ParseLayer([]byte("prompts:\n  default:\n    path: missing.md"))
	if err != nil {
		t.Fatal(err)
	}
	if layer == nil {
		t.Fatal("ParseLayer failed")
	}
	provider := &FileLayerProvider{Layer: layer}
	_, err = ResolvePromptSource(provider, dir, ResolveOptions{Alias: "default"})
	if err == nil {
		t.Fatal("ResolvePromptSource(alias missing file) err = nil, want ErrAliasSourceMissing")
	}
	if !errors.Is(err, ErrAliasSourceMissing) {
		t.Errorf("err = %v, want ErrAliasSourceMissing", err)
	}
	if !IsExit2(err) {
		t.Error("IsExit2(err) = false, want true")
	}
}
