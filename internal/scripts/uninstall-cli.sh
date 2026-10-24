#!/usr/bin/env bash
set -euo pipefail

readonly BINARY_PATH="${HOME}/.local/bin/podlet"
readonly DATA_DIR="${XDG_DATA_HOME:-${HOME}/.local/share}/podlet"
readonly UNINSTALLER_PATH="${DATA_DIR}/uninstall-cli.sh"
readonly USER_CONFIG_DIR="${XDG_CONFIG_HOME:-${HOME}/.config}"
readonly CONFIG_DIR="${USER_CONFIG_DIR}/podlet"
readonly CONFIG_FILE="${CONFIG_DIR}/config.json"

if [[ "${EUID}" -eq 0 ]]; then
    echo "Error: run this uninstaller without sudo or root." >&2
    exit 1
fi

if [[ ! -x "$BINARY_PATH" ]]; then
    echo "Error: Podlet CLI was not found at ${BINARY_PATH}. Nothing was removed." >&2
    exit 1
fi

echo "Logging out of Podlet..."
if ! "$BINARY_PATH" auth logout; then
    echo "Error: logout failed. Nothing was removed." >&2
    exit 1
fi

rm -f "$BINARY_PATH"
rm -f "$CONFIG_FILE"
rmdir "$CONFIG_DIR" 2>/dev/null || true
rm -f "$UNINSTALLER_PATH"
rmdir "$DATA_DIR" 2>/dev/null || true

echo "Podlet CLI uninstalled successfully."
echo "Run 'hash -r' if your shell still resolves the podlet command."
