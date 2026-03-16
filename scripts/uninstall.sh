#!/usr/bin/env sh
# Uninstall Ralph by removing the binary from the first standard location that contains it.
# Checks (in order): /usr/local/bin, ~/.local/bin, ~/bin — same as install.sh defaults.
# If you installed with --dir to a custom path, remove that binary manually.
# User config (e.g. ~/.config/ralph/ralph-config.yml) is not removed.

set -e

# --- Setup ---
CONFIG_HOME="${RALPH_CONFIG_HOME:-$HOME/.config/ralph}"

# --- 1. Detect platform (SUF = .exe on Windows, else ""; matches install.sh binary name) ---
_os="$(uname -s | tr '[:upper:]' '[:lower:]')"
case "$_os" in
  mingw*|msys*|cygwin*) SUF=".exe" ;;
  *)                   SUF="" ;;
esac

# --- Helpers ---
# Candidate directories in same order as install.sh defaults.
get_candidate_dirs() {
  echo "/usr/local/bin"
  [ -n "$HOME" ] && echo "${HOME}/.local/bin"
  [ -n "$HOME" ] && echo "${HOME}/bin"
}

# First candidate dir that contains ralph or ralph.exe; empty if none.
find_install_dir() {
  for _dir in $(get_candidate_dirs); do
    if [ -f "${_dir}/ralph" ] || [ -f "${_dir}/ralph.exe" ]; then
      echo "$_dir"
      return
    fi
  done
}

# --- 2. Find install location ---
INSTALL_DIR="$(find_install_dir)"

if [ -z "$INSTALL_DIR" ]; then
  echo "ralph binary not found in standard locations (/usr/local/bin, ~/.local/bin, ~/bin)." >&2
  echo "If you installed with --dir to a custom path, remove that binary manually." >&2
  exit 1
fi

# --- 3. Remove binary ---
if [ -f "${INSTALL_DIR}/ralph" ]; then
  rm -f "${INSTALL_DIR}/ralph"
  echo "Removed ${INSTALL_DIR}/ralph"
fi
if [ -f "${INSTALL_DIR}/ralph.exe" ]; then
  rm -f "${INSTALL_DIR}/ralph.exe"
  echo "Removed ${INSTALL_DIR}/ralph.exe"
fi

# --- 4. Done ---
echo "Uninstall complete."
echo "User config in ${CONFIG_HOME}/ (e.g. ralph-config.yml) was not removed."
