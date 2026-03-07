#!/usr/bin/env sh
# install.sh - Install Ralph binary from a GitHub release to a directory on PATH.
# O6. Spec: R1, R2, R3. Only installs from release artifacts (no build from source).
# Usage: ./scripts/install.sh [VERSION] [--dir <path>]
#   VERSION: optional; e.g. 1.0.0 or v1.0.0. Omit for latest release.
# Env: RALPH_INSTALL_DIR (override default dir), RALPH_REPO (default: maxdunn/ralph).
# Requires: curl. State file: ~/.config/ralph/install-state.

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

# Detect OS and arch for release asset name (ralph-<version>-<os>-<arch>[.exe]).
detect_os_arch() {
  _kernel=""
  _arch=""
  if command -v uname >/dev/null 2>&1; then
    _kernel="$(uname -s)"
    _arch="$(uname -m)"
  fi
  case "$_kernel" in
    Darwin)   _os="darwin" ;;
    Linux)    _os="linux" ;;
    MINGW*|MSYS*|CYGWIN*) _os="windows" ;;
    *)        _os="" ;;
  esac
  case "$_arch" in
    x86_64|amd64) _arch="amd64" ;;
    aarch64|arm64) _arch="arm64" ;;
    *) _arch="" ;;
  esac
  echo "${_os:-unknown}"
  echo "${_arch:-unknown}"
}

# --- Parse arguments: [VERSION] [--dir PATH] (version is first positional; --dir can appear anywhere)
VERSION_ARG=""
INSTALL_DIR=""
while [ $# -gt 0 ]; do
  if [ "$1" = "--dir" ] && [ -n "${2:-}" ]; then
    INSTALL_DIR="$2"
    shift 2
  elif [ -z "$VERSION_ARG" ] && [ "$1" != "--dir" ]; then
    VERSION_ARG="$1"
    shift
  else
    shift
  fi
done
if [ -n "$RALPH_INSTALL_DIR" ]; then
  INSTALL_DIR="${INSTALL_DIR:-$RALPH_INSTALL_DIR}"
fi
if [ -z "$INSTALL_DIR" ]; then
  INSTALL_DIR="${HOME:?}/bin"
fi
case "$INSTALL_DIR" in
  ~*) INSTALL_DIR="$HOME${INSTALL_DIR#\~}" ;;
esac
mkdir -p "$INSTALL_DIR"
INSTALL_DIR_ABS="$(resolve_dir "$INSTALL_DIR")"

# --- Require curl
if ! command -v curl >/dev/null 2>&1; then
  echo "Error: curl is required to download the release binary." >&2
  exit 1
fi

# --- Resolve tag (v-prefixed) and version (no v) for asset name
RALPH_REPO="${RALPH_REPO:-maxdunn/ralph}"
if [ -n "$VERSION_ARG" ]; then
  _ver="${VERSION_ARG#v}"
  _tag="v${_ver}"
else
  _api="https://api.github.com/repos/$RALPH_REPO/releases/latest"
  _tag=""
  _tag="$(curl -sSfL "$_api" 2>/dev/null | grep -o '"tag_name":"[^"]*"' | head -n1 | cut -d'"' -f4)" || true
  if [ -z "$_tag" ]; then
    echo "Error: Could not determine latest release (check network and $RALPH_REPO releases)." >&2
    exit 1
  fi
  _ver="${_tag#v}"
fi

# --- Detect OS/arch and download
_detected="$(detect_os_arch)"
_os="$(echo "$_detected" | head -n1)"
_arch="$(echo "$_detected" | tail -n1)"
if [ "$_os" = "unknown" ] || [ "$_arch" = "unknown" ]; then
  echo "Error: Unsupported OS/arch. Supported: Linux, macOS, Windows (amd64, arm64)." >&2
  exit 1
fi
case "$_os" in
  windows) _suf=".exe"; _out="ralph.exe" ;;
  *)       _suf=""; _out="ralph" ;;
esac
_asset="ralph-${_ver}-${_os}-${_arch}${_suf}"
_url="https://github.com/${RALPH_REPO}/releases/download/${_tag}/${_asset}"
_tmp="/tmp/ralph-install-$$"
if ! curl -sSfL -o "$_tmp" "$_url" 2>/dev/null || [ ! -f "$_tmp" ]; then
  rm -f "$_tmp"
  echo "Error: Download failed for $_asset (version $_ver). Check that the release exists: $RALPH_REPO/releases" >&2
  exit 1
fi
chmod +x "$_tmp"

# --- Install
cp -f "$_tmp" "$INSTALL_DIR_ABS/${_out}"
chmod +x "$INSTALL_DIR_ABS/${_out}"
rm -f "$_tmp"

# --- Write state file
STATE_DIR="${HOME}/.config/ralph"
STATE_FILE="$STATE_DIR/install-state"
mkdir -p "$STATE_DIR"
echo "$INSTALL_DIR_ABS" > "$STATE_FILE"

echo "Installed ralph ${_ver} to $INSTALL_DIR_ABS."
echo "Run 'ralph version' in a new terminal to verify. Ensure that directory is on your PATH (e.g. add $INSTALL_DIR_ABS to PATH if needed)."
