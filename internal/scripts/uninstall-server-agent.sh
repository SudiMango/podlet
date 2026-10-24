#!/usr/bin/env bash
set -euo pipefail

readonly BINARY_PATH="${HOME}/.local/bin/podlet-agent"
readonly DATA_DIR="${XDG_DATA_HOME:-${HOME}/.local/share}/podlet"
readonly UNINSTALLER_PATH="${DATA_DIR}/uninstall-server-agent.sh"
readonly USER_CONFIG_DIR="${XDG_CONFIG_HOME:-${HOME}/.config}"
readonly CONFIG_DIR="${USER_CONFIG_DIR}/podlet"
readonly CONFIG_FILE="${CONFIG_DIR}/server.json"
readonly SERVICE_FILE="${USER_CONFIG_DIR}/systemd/user/podlet-agent.service"

if [[ "${EUID}" -eq 0 ]]; then
    echo "Error: run this uninstaller without sudo or root." >&2
    exit 1
fi

FROM_AGENT="false"
if [[ "${1:-}" == "--from-agent" ]]; then
    FROM_AGENT="true"
elif [[ "$#" -gt 0 ]]; then
    echo "Error: unknown argument: $1" >&2
    exit 1
fi

if ! command -v systemctl >/dev/null 2>&1; then
    echo "Error: required command not found: systemctl" >&2
    exit 1
fi

if [[ "$FROM_AGENT" == "true" ]]; then
    systemctl --user disable podlet-agent.service >/dev/null 2>&1 || true
else
    systemctl --user disable --now podlet-agent.service >/dev/null 2>&1 || true
fi
rm -f "$SERVICE_FILE"
systemctl --user daemon-reload
systemctl --user reset-failed podlet-agent.service >/dev/null 2>&1 || true

rm -f "$BINARY_PATH"
rm -f "$CONFIG_FILE"
rmdir "$CONFIG_DIR" 2>/dev/null || true
rm -f "$UNINSTALLER_PATH"
rmdir "$DATA_DIR" 2>/dev/null || true

echo "Podlet server agent uninstalled successfully."
