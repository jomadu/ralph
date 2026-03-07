#!/usr/bin/env sh
# uninstall.sh - Remove Ralph binary and install state file written by install.sh.
# O6 T3. Spec: R4. Plan: PLAN_INSTALL_UNINSTALL.md.
# Usage: ./scripts/uninstall.sh
# State file: ~/.config/ralph/install-state (single line: absolute install directory).

set -e

STATE_DIR="${HOME}/.config/ralph"
STATE_FILE="$STATE_DIR/install-state"

# --- Discover install location from state file
if [ ! -f "$STATE_FILE" ]; then
  echo "Ralph does not appear to be installed (no install state found)." >&2
  exit 0
fi

INSTALL_DIR_ABS="$(cat "$STATE_FILE" | head -n1)"
# Basic sanity: must be an absolute path
case "$INSTALL_DIR_ABS" in
  /*) ;;
  *)
    echo "Error: Invalid install state (expected absolute path)." >&2
    exit 2
    ;;
esac

# --- Remove binary/binary.exe (rm -f so missing binary does not fail)
rm -f "$INSTALL_DIR_ABS/ralph" "$INSTALL_DIR_ABS/ralph.exe"

# --- Remove state file so future install gets clean state
rm -f "$STATE_FILE"

echo "Uninstalled ralph from $INSTALL_DIR_ABS."
