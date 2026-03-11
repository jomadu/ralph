#!/usr/bin/env sh
# Uninstall Ralph by removing the binary from the directory recorded at install time
# and removing the install state file. Run from anywhere; uses ~/.config/ralph/install-state.
# User config (e.g. ~/.config/ralph/ralph-config.yml) is not removed.
# No PATH or symlink changes are made by the install script, so uninstall leaves no broken references.

set -e

INSTALL_STATE_DIR="${RALPH_CONFIG_HOME:-$HOME/.config/ralph}"
INSTALL_STATE_FILE="${INSTALL_STATE_DIR}/install-state"

if [ ! -f "$INSTALL_STATE_FILE" ]; then
  echo "No install state found at ${INSTALL_STATE_FILE}. Nothing to uninstall." >&2
  exit 1
fi

INSTALL_DIR="$(cat "$INSTALL_STATE_FILE" | head -1)"
if [ -z "$INSTALL_DIR" ] || [ ! -d "$INSTALL_DIR" ]; then
  echo "Invalid or missing install directory in state file. Removing state file." >&2
  rm -f "$INSTALL_STATE_FILE"
  exit 1
fi

# Remove binary (ralph or ralph.exe)
if [ -f "${INSTALL_DIR}/ralph" ]; then
  rm -f "${INSTALL_DIR}/ralph"
  echo "Removed ${INSTALL_DIR}/ralph"
fi
if [ -f "${INSTALL_DIR}/ralph.exe" ]; then
  rm -f "${INSTALL_DIR}/ralph.exe"
  echo "Removed ${INSTALL_DIR}/ralph.exe"
fi

# Remove install state so future uninstall does not fail
rm -f "$INSTALL_STATE_FILE"
echo "Removed install state. Uninstall complete."
echo "User config in ${INSTALL_STATE_DIR}/ (e.g. ralph-config.yml) was not removed."
