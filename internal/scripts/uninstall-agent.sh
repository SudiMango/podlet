#!/usr/bin/env bash
set -euo pipefail

readonly BINARY_PATH="/usr/local/bin/podlet-agent"
readonly DATA_DIR="/var/lib/podlet"
readonly CONFIG_DIR="/etc/podlet"
readonly SERVICE_FILE="/etc/systemd/system/podlet-agent.service"

# --------------------------------------------------
# Root check
# --------------------------------------------------

if [[ "${EUID}" -ne 0 ]]; then
    echo "Error: this uninstaller must run as root." >&2
    exit 1
fi

# --------------------------------------------------
# Parse arguments
# --------------------------------------------------

FROM_AGENT="false"

if [[ "${1:-}" == "--from-agent" ]]; then
    FROM_AGENT="true"
elif [[ "$#" -gt 0 ]]; then
    echo "Error: unknown argument: $1" >&2
    exit 1
fi

# --------------------------------------------------
# Verify prerequisites
# --------------------------------------------------

if ! command -v systemctl >/dev/null 2>&1; then
    echo "Error: required command not found: systemctl" >&2
    exit 1
fi

# --------------------------------------------------
# Disable and stop service
# --------------------------------------------------

# When invoked by the agent itself, avoid stopping
# the currently running service during uninstallation.
if [[ "$FROM_AGENT" == "true" ]]; then
    systemctl disable podlet-agent.service >/dev/null 2>&1 || true
else
    systemctl disable --now podlet-agent.service >/dev/null 2>&1 || true
fi

# --------------------------------------------------
# Remove systemd service
# --------------------------------------------------

rm -f "$SERVICE_FILE"

systemctl daemon-reload
systemctl reset-failed podlet-agent.service >/dev/null 2>&1 || true

# --------------------------------------------------
# Remove agent files
# --------------------------------------------------

rm -f "$BINARY_PATH"
rm -rf "$CONFIG_DIR"
rm -rf "$DATA_DIR"

# --------------------------------------------------
# Uninstallation complete
# --------------------------------------------------

echo "Podlet server agent uninstalled successfully."
