#!/usr/bin/env sh
# install.sh - Install Ralph binary to a directory on PATH and record state for uninstall.
# O6 T2. Spec: R2, R3. Plan: PLAN_INSTALL_UNINSTALL.md.
# Usage: ./scripts/install.sh [--dir <path>]
# Env: RALPH_INSTALL_DIR (override default), RALPH_BINARY (path to pre-built binary).
# Default install dir: $HOME/bin. State file: ~/.config/ralph/install-state.

set -e

# Resolve directory to absolute path (portable).
resolve_dir() {
  _d="$1"
  if command -v realpath >/dev/null 2>&1; then
    realpath "$_d"
  else
    (cd "$_d" && pwd -P)
  fi
}

# Find repo root (directory containing go.mod).
find_repo_root() {
  _p="$1"
  while [ -n "$_p" ] && [ "$_p" != "/" ]; do
    [ -f "$_p/go.mod" ] && echo "$_p" && return 0
    _p="$(dirname "$_p")"
  done
  return 1
}

# --- Resolve install directory (--dir overrides env overrides default)
INSTALL_DIR=""
if [ "$1" = "--dir" ] && [ -n "$2" ]; then
  INSTALL_DIR="$2"
  shift 2
elif [ -n "$RALPH_INSTALL_DIR" ]; then
  INSTALL_DIR="$RALPH_INSTALL_DIR"
fi
if [ -z "$INSTALL_DIR" ]; then
  INSTALL_DIR="${HOME:?}/bin"
fi
# Expand ~ if present
case "$INSTALL_DIR" in
  ~*) INSTALL_DIR="$HOME${INSTALL_DIR#\~}" ;;
esac
mkdir -p "$INSTALL_DIR"
INSTALL_DIR_ABS="$(resolve_dir "$INSTALL_DIR")"

# --- Obtain binary
RALPH_SRC=""
if [ -n "$RALPH_BINARY" ] && [ -f "$RALPH_BINARY" ]; then
  RALPH_SRC="$RALPH_BINARY"
elif [ -n "$RALPH_BINARY" ]; then
  echo "Error: RALPH_BINARY is set but file not found: $RALPH_BINARY" >&2
  exit 1
fi

if [ -z "$RALPH_SRC" ]; then
  # Try to build from repo (script may be run from repo root or from scripts/)
  SCRIPT_DIR="$(cd "$(dirname "$0")" && pwd -P)"
  REPO_ROOT=""
  if [ -f "$SCRIPT_DIR/../go.mod" ]; then
    REPO_ROOT="$(resolve_dir "$SCRIPT_DIR/..")"
  fi
  if [ -z "$REPO_ROOT" ]; then
    echo "Error: No Ralph binary available. Either:" >&2
    echo "  1. Run this script from the Ralph repo root (with Go and cmd/ralph) so it can build, or" >&2
    echo "  2. Set RALPH_BINARY to the path of a pre-built ralph binary." >&2
    exit 1
  fi
  BUILD_OUT="/tmp/ralph-install-$$"
  if ! (cd "$REPO_ROOT" && go build -o "$BUILD_OUT" ./cmd/ralph 2>/dev/null) || [ ! -f "$BUILD_OUT" ]; then
    rm -f "$BUILD_OUT"
    echo "Error: Could not build ralph from $REPO_ROOT (missing cmd/ralph or go build failed)." >&2
    echo "  Set RALPH_BINARY to a pre-built binary, or add Go source under cmd/ralph." >&2
    exit 1
  fi
  RALPH_SRC="$BUILD_OUT"
fi

# --- Copy binary and set executable
cp -f "$RALPH_SRC" "$INSTALL_DIR_ABS/ralph"
chmod +x "$INSTALL_DIR_ABS/ralph"
# Remove temp build artifact if we built from source
case "$RALPH_SRC" in /tmp/ralph-install-*) rm -f "$RALPH_SRC" ;; esac

# --- Write state file (single line: absolute install directory)
STATE_DIR="${HOME}/.config/ralph"
STATE_FILE="$STATE_DIR/install-state"
mkdir -p "$STATE_DIR"
echo "$INSTALL_DIR_ABS" > "$STATE_FILE"

# --- Success message
echo "Installed ralph to $INSTALL_DIR_ABS."
echo "Run 'ralph version' in a new terminal to verify. Ensure that directory is on your PATH (e.g. add $INSTALL_DIR_ABS to PATH if needed)."
