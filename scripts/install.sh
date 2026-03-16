#!/usr/bin/env sh
# Install Ralph from a GitHub release. Usage: install.sh [VERSION] [--dir DIR]
# Default: latest release; install directory is platform- and privilege-dependent (see default_install_dir).
# Uninstall looks for the binary in the same standard locations (no state file).
# Requires: curl, jq. Supported: Linux, macOS, Windows (Git Bash); amd64, arm64.
#
# Default install locations (when RALPH_INSTALL_DIR and --dir are not set):
#   Linux:   /usr/local/bin if writable (e.g. sudo), else ~/.local/bin (XDG)
#   macOS:   /usr/local/bin if writable, else ~/bin
#   Windows: $HOME/bin

set -e

REPO="${RALPH_REPO:-jomadu/ralph}"
GITHUB_API="https://api.github.com/repos/${REPO}/releases"
BINARY_NAME="ralph"

usage() {
  echo "Usage: $0 [VERSION] [--dir DIR]"
  echo "  VERSION  Optional. Tag name (e.g. 1.0.0 or v1.0.0). Default: latest release."
  echo "  --dir    Optional. Install directory. Default: platform-dependent (see below) or RALPH_INSTALL_DIR."
  echo "  Default dir: Linux: /usr/local/bin if writable else ~/.local/bin; macOS: /usr/local/bin if writable else ~/bin; Windows: ~/bin."
  echo "Env: RALPH_REPO, RALPH_INSTALL_DIR"
  exit 0
}

# Default install directory by platform and privilege (FHS/XDG conventions).
# Use /usr/local/bin when it exists and is writable (system-wide); else user-local.
default_install_dir() {
  OS="$(uname -s | tr '[:upper:]' '[:lower:]')"
  case "$OS" in
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
    mingw*|msys*|cygwin*)
      echo "${HOME}/bin"
      ;;
    *)
      echo "${HOME}/.local/bin"
      ;;
  esac
}

# Parse args (INSTALL_DIR from env or --dir; if unset after parsing, use default_install_dir)
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

# Normalize version tag (add v if missing). No trailing newline (safe for URLs).
normalize_tag() {
  case "$1" in
    v*) printf '%s' "$1" ;;
    *)  printf '%s' "v$1" ;;
  esac
}

# Detect OS and arch for artifact name (matches Makefile build-multi)
detect_platform() {
  OS="$(uname -s | tr '[:upper:]' '[:lower:]')"
  ARCH="$(uname -m)"
  case "$OS" in
    darwin) GOOS="darwin" ;;
    linux)  GOOS="linux" ;;
    mingw*|msys*|cygwin*) GOOS="windows" ;;
    *) echo "Unsupported OS: $OS" >&2; exit 1 ;;
  esac
  case "$ARCH" in
    x86_64|amd64) GOARCH="amd64" ;;
    aarch64|arm64) GOARCH="arm64" ;;
    *) echo "Unsupported arch: $ARCH" >&2; exit 1 ;;
  esac
  if [ "$GOOS" = "windows" ]; then
    SUF=".exe"
  else
    SUF=""
  fi
  echo "${GOOS} ${GOARCH} ${SUF}"
}

# Resolve version: latest or specific tag (GitHub tag may be v1.0.0 or 1.0.0).
# Uses release tag_name (not name) for download URL.
get_version() {
  if [ -z "$VERSION" ]; then
    curl -sSfL "${GITHUB_API}/latest" | jq -r '.tag_name' | tr -d '\n'
  else
    normalize_tag "$VERSION"
  fi
}

# Artifact filename for a given tag and platform (matches Makefile build-multi output)
artifact_name() {
  TAG="$1"
  TAG_STRIP="${TAG#v}"
  echo "ralph-${TAG_STRIP}-${GOOS}-${GOARCH}${SUF}"
}

download_url() {
  TAG="$1"
  NAME="$(artifact_name "$TAG")"
  # GitHub release asset URL: /repos/OWNER/REPO/releases/assets/ASSET_ID (requires Accept header)
  # Simpler: redirect to tarball/zip or use direct asset by name from releases/tag/TAG
  # Standard pattern: https://github.com/OWNER/REPO/releases/download/TAG/ASSET_NAME
  echo "https://github.com/${REPO}/releases/download/${TAG}/${NAME}"
}

# Main
PLATFORM="$(detect_platform)"
GOOS="${PLATFORM%% *}"
REST="${PLATFORM#* }"
GOARCH="${REST%% *}"
SUF="${REST#* }"

TAG="$(get_version)"
ARTIFACT="$(artifact_name "$TAG")"
URL="$(download_url "$TAG")"

mkdir -p "$INSTALL_DIR"
if ! [ -d "$INSTALL_DIR" ]; then
  echo "Cannot create or use install directory: $INSTALL_DIR" >&2
  exit 1
fi

echo "Installing Ralph ${TAG} to ${INSTALL_DIR}..."
if ! curl -sSfL -o "${INSTALL_DIR}/ralph${SUF}" "$URL"; then
  echo "Download failed. If the release does not exist yet, build locally: make build && cp bin/ralph ${INSTALL_DIR}/" >&2
  exit 1
fi
chmod +x "${INSTALL_DIR}/ralph${SUF}"

echo "Install complete. Binary: ${INSTALL_DIR}/ralph${SUF}"
echo "Ensure ${INSTALL_DIR} is on your PATH, then run: ralph version"
