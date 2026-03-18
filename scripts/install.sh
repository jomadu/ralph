#!/usr/bin/env sh
# Install Ralph from a GitHub release. Usage: install.sh [RELEASE] [--dir DIR]
# RELEASE is the version without a leading v (e.g. 1.0.0-rc.14). Same form is used for the release tag and download URL.
# Default: latest release (GitHub /releases/latest).
# Install directory is platform- and privilege-dependent.
# Uninstall looks for the binary in standard locations (no state file).
# Requires: curl, jq. Supported: Linux, macOS, Windows (Git Bash); amd64, arm64.
#
# Default install dir: Linux: /usr/local/bin if writable else ~/.local/bin;
#   macOS: /usr/local/bin if writable else ~/bin; Windows: ~/bin.

set -e

# --- Setup ---
REPO="${RALPH_REPO:-jomadu/ralph}"
GITHUB_API="https://api.github.com/repos/${REPO}/releases"

# --- 1. Detect platform (sets GOOS, GOARCH, SUF for rest of script) ---
_os="$(uname -s | tr '[:upper:]' '[:lower:]')"
_arch="$(uname -m)"
case "$_os" in
  linux)  GOOS="linux" ;;
  darwin) GOOS="darwin" ;;
  mingw*|msys*|cygwin*) GOOS="windows" ;;
  *) echo "Unsupported OS: $_os" >&2; exit 1 ;;
esac
case "$_arch" in
  x86_64|amd64) GOARCH="amd64" ;;
  aarch64|arm64) GOARCH="arm64" ;;
  *) echo "Unsupported arch: $_arch" >&2; exit 1 ;;
esac
if [ "$GOOS" = "windows" ]; then
  SUF=".exe"
else
  SUF=""
fi

# --- Helpers ---
usage() {
  echo "Usage: $0 [RELEASE] [--dir DIR]"
  echo "  RELEASE  Optional. Release version without v (e.g. 1.0.0, 1.0.0-rc.14). A leading v is ignored."
  echo "  --dir    Optional. Install directory. Default: platform-dependent or RALPH_INSTALL_DIR."
  echo "Env: RALPH_REPO, RALPH_INSTALL_DIR"
  exit 0
}

default_install_dir() {
  case "$GOOS" in
    linux)
      if [ -d /usr/local/bin ] && [ -w /usr/local/bin ] 2>/dev/null; then
        echo "/usr/local/bin"
      else
        echo "${HOME}/.local/bin"
      fi
      ;;
    darwin)
      if [ -d /usr/local/bin ] && [ -w /usr/local/bin ] 2>/dev/null; then
        echo "/usr/local/bin"
      else
        echo "${HOME}/bin"
      fi
      ;;
    windows) echo "${HOME}/bin" ;;
    *)       echo "${HOME}/.local/bin" ;;
  esac
}

# Release version without leading v (matches release name; artifact uses this form).
release_version() {
  case "$1" in
    v*) printf '%s' "${1#v}" ;;
    *)  printf '%s' "$1" ;;
  esac
}

resolve_release() {
  if [ -n "$VERSION" ]; then
    release_version "$VERSION"
  else
    _api_tag="$(curl -sSfL "${GITHUB_API}/latest" | jq -r '.tag_name' | tr -d '\n')"
    release_version "$_api_tag"
  fi
}

artifact_name() {
  _rel="$1"
  echo "ralph-${_rel}-${GOOS}-${GOARCH}${SUF}"
}

download_url() {
  _rel="$1"
  _name="$(artifact_name "$_rel")"
  echo "https://github.com/${REPO}/releases/download/${_rel}/${_name}"
}

is_dir_in_path() {
  _dir="$1"
  _abs="$(cd "$_dir" && pwd)"
  _rest="${PATH}"
  while [ -n "$_rest" ]; do
    case "$_rest" in
      *:*) _elem="${_rest%%:*}"; _rest="${_rest#*:}" ;;
      *)   _elem="${_rest}"; _rest="" ;;
    esac
    if [ -n "$_elem" ]; then
      _elem_abs="$(cd "$_elem" 2>/dev/null && pwd)" || true
      if [ "$_elem_abs" = "$_abs" ]; then return 0; fi
    fi
  done
  return 1
}

print_path_instructions() {
  _dir="$1"
  _abs="$(cd "$_dir" && pwd)"
  echo ""
  echo "The install directory is not on your PATH. Add it so you can run 'ralph' from any terminal:"
  echo ""
  echo "  export PATH=\"${_abs}:\$PATH\""
  echo ""
  echo "To make this permanent, add the line above to your shell profile:"
  echo "  - Bash: ~/.bashrc or ~/.bash_profile"
  echo "  - Zsh:  ~/.zshrc"
  echo "  - Fish: run 'set -U fish_user_paths ${_abs} \$fish_user_paths'"
  echo ""
  echo "Then run: ralph version"
}

# --- 2. Parse args ---
VERSION=""
INSTALL_DIR="${RALPH_INSTALL_DIR:-}"
while [ $# -gt 0 ]; do
  case "$1" in
    -h|--help) usage ;;
    --dir)
      [ -n "$2" ] || { echo "Missing value for --dir" >&2; exit 1; }
      INSTALL_DIR="$2"
      shift 2
      ;;
    *)
      if [ -z "$VERSION" ]; then
        VERSION="$1"
        shift
      else
        echo "Unexpected argument: $1" >&2
        exit 1
      fi
      ;;
  esac
done
if [ -z "$INSTALL_DIR" ]; then
  INSTALL_DIR="$(default_install_dir)"
fi

# --- 3. Resolve release (no v) ---
RELEASE="$(resolve_release)"

# --- 4. Build download URL ---
ARTIFACT="$(artifact_name "$RELEASE")"
URL="$(download_url "$RELEASE")"

# --- 5. Install ---
mkdir -p "$INSTALL_DIR"
if ! [ -d "$INSTALL_DIR" ]; then
  echo "Cannot create or use install directory: $INSTALL_DIR" >&2
  exit 1
fi

echo "Installing Ralph ${RELEASE} to ${INSTALL_DIR}..."
if ! curl -sSfL -o "${INSTALL_DIR}/ralph${SUF}" "$URL"; then
  echo "Download failed. If the release does not exist yet, build locally: make build && cp bin/ralph ${INSTALL_DIR}/" >&2
  exit 1
fi
chmod +x "${INSTALL_DIR}/ralph${SUF}"

# --- 6. Path check and message ---
echo "Install complete. Binary: ${INSTALL_DIR}/ralph${SUF}"
if is_dir_in_path "$INSTALL_DIR"; then
  echo "Install directory is on your PATH. Run: ralph version"
else
  print_path_instructions "$INSTALL_DIR"
fi
