#!/usr/bin/env sh
# Uninstall Ralph by removing the binary from the first standard location that contains it.
# Checks (in order): /usr/local/bin, ~/.local/bin, ~/bin (same as install.sh defaults).
# If you installed with --dir to a custom path, remove that binary manually.
# User config (e.g. ~/.config/ralph/ralph-config.yml) is not removed.

set -e

# Candidate directories in same order as install.sh defaults (FHS/XDG + fallbacks).
# Uninstall removes from the first candidate that contains ralph or ralph.exe.
get_candidate_dirs() {
  if [ -d /usr/local/bin ]; then
    echo "/usr/local/bin"
  fi
  if [ -n "$HOME" ] && [ -d "${HOME}/.local/bin" ]; then
    echo "${HOME}/.local/bin"
  fi
  if [ -n "$HOME" ] && [ -d "${HOME}/bin" ]; then
    echo "${HOME}/bin"
  fi
}

FOUND=""
for DIR in $(get_candidate_dirs); do
  if [ -f "${DIR}/ralph" ] || [ -f "${DIR}/ralph.exe" ]; then
    FOUND="$DIR"
    break
  fi
done

if [ -z "$FOUND" ]; then
  echo "ralph binary not found in standard locations (/usr/local/bin, ~/.local/bin, ~/bin)." >&2
  echo "If you installed with --dir to a custom path, remove that binary manually." >&2
  exit 1
fi

if [ -f "${FOUND}/ralph" ]; then
  rm -f "${FOUND}/ralph"
  echo "Removed ${FOUND}/ralph"
fi
if [ -f "${FOUND}/ralph.exe" ]; then
  rm -f "${FOUND}/ralph.exe"
  echo "Removed ${FOUND}/ralph.exe"
fi

echo "Uninstall complete."
echo "User config in ${RALPH_CONFIG_HOME:-$HOME/.config/ralph}/ (e.g. ralph-config.yml) was not removed."
