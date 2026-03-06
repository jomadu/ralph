#!/usr/bin/env bash
# Simple Ralph-style loop: pipe prompt to cursor-wrapper, scan for signals, repeat or exit.
# Usage:
#   ./ralph-loop.sh run <prompt-file> [max-iterations]
#   cat prompt.md | ./ralph-loop.sh run
#
# Requires: cursor-wrapper.sh in same directory (or in scripts/), jq, cursor agent.

set -e
RALPH_ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
WRAPPER="${RALPH_ROOT}/cursor-wrapper.sh"
if [[ ! -x "$WRAPPER" ]]; then
  WRAPPER="${RALPH_ROOT}/scripts/cursor-wrapper.sh"
fi
if [[ ! -x "$WRAPPER" ]]; then
  echo "ralph-loop.sh: cursor-wrapper.sh not found in ${RALPH_ROOT} or ${RALPH_ROOT}/scripts" >&2
  exit 1
fi

SUCCESS_SIGNAL="<promise>SUCCESS</promise>"
FAILURE_SIGNAL="<promise>FAILURE</promise>"
MAX_ITERATIONS=5

usage() {
  echo "Usage: $0 run <prompt-file> [max-iterations]"
  echo "       cat prompt.md | $0 run [max-iterations]"
  echo "Default max-iterations: ${MAX_ITERATIONS}"
  exit 0
}

if [[ "${1:-}" = "" || "${1:-}" = "-h" || "${1:-}" = "--help" ]]; then
  usage
fi

if [[ "$1" != "run" ]]; then
  echo "Expected first argument: run" >&2
  usage
fi
shift

if [[ $# -gt 0 && -f "$1" && "$1" != *[0-9]* ]]; then
  PROMPT_FILE="$1"
  shift
  read_prompt() { cat "$PROMPT_FILE"; }
else
  PROMPT_FILE=""
  # Stdin only available once; buffer it for all iterations
  STDIN_PROMPT=$(cat)
  read_prompt() { echo "$STDIN_PROMPT"; }
fi

if [[ $# -gt 0 && "$1" =~ ^[0-9]+$ ]]; then
  MAX_ITERATIONS="$1"
  shift
fi

iter=0
while [[ $iter -lt $MAX_ITERATIONS ]]; do
  iter=$((iter + 1))
  echo "[ralph-loop] iteration $iter/$MAX_ITERATIONS" >&2
  output=$(read_prompt | "$WRAPPER" 2>&1) || true
  echo "$output"

  if echo "$output" | grep -q "$SUCCESS_SIGNAL"; then
    echo "[ralph-loop] success signal found, exiting 0" >&2
    exit 0
  fi
  if echo "$output" | grep -q "$FAILURE_SIGNAL"; then
    echo "[ralph-loop] failure signal found, exiting 1" >&2
    exit 1
  fi
done

echo "[ralph-loop] max iterations ($MAX_ITERATIONS) reached, exiting 2" >&2
exit 2
