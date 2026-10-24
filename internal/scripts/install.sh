#!/usr/bin/env bash
set -euo pipefail

readonly API_URL="https://api.podlet.dev"
readonly AGENT_VERSION="v0.1.0"
readonly DOWNLOAD_BASE_URL="https://downloads.podlet.dev/server-agent/${AGENT_VERSION}"
readonly BINARY_DIR="${HOME}/.local/bin"
readonly BINARY_PATH="${BINARY_DIR}/podlet-agent"
readonly USER_CONFIG_DIR="${XDG_CONFIG_HOME:-${HOME}/.config}"
readonly CONFIG_DIR="${USER_CONFIG_DIR}/podlet"
readonly CONFIG_FILE="${CONFIG_DIR}/server.json"
readonly SERVICE_DIR="${USER_CONFIG_DIR}/systemd/user"
readonly SERVICE_FILE="${SERVICE_DIR}/podlet-agent.service"

if [[ "${EUID}" -eq 0 ]]; then
    echo "Error: run this installer without sudo or root." >&2
    echo "Use: curl -fsSL https://downloads.podlet.dev/server-agent/install.sh | bash" >&2
    exit 1
fi

if [[ "$(uname -s)" != "Linux" ]]; then
    echo "Error: the Podlet agent currently supports Linux only." >&2
    exit 1
fi

for command in curl jq sha256sum systemctl install awk; do
    if ! command -v "$command" >/dev/null 2>&1; then
        echo "Error: required command not found: $command" >&2
        exit 1
    fi
done

if [[ ! -r /dev/tty ]]; then
    echo "Error: an interactive terminal is required for registration." >&2
    exit 1
fi

read -r -p "Server name: " SERVER_NAME </dev/tty
if [[ -z "${SERVER_NAME//[[:space:]]/}" ]]; then
    echo "Error: server name is required." >&2
    exit 1
fi

read -r -s -p "Registration token: " REGISTRATION_TOKEN </dev/tty
printf '\n' >/dev/tty
if [[ -z "$REGISTRATION_TOKEN" ]]; then
    echo "Error: registration token is required." >&2
    exit 1
fi

case "$(uname -m)" in
    x86_64 | amd64)
        ARCH="amd64"
        ;;
    aarch64 | arm64)
        ARCH="arm64"
        ;;
    *)
        echo "Error: unsupported architecture: $(uname -m)" >&2
        exit 1
        ;;
esac

readonly BINARY_NAME="podlet-agent-linux-${ARCH}"
TEMP_DIR="$(mktemp -d)"
trap 'rm -rf "$TEMP_DIR"' EXIT

echo "Downloading Podlet agent ${AGENT_VERSION}..."
curl -fsSL "${DOWNLOAD_BASE_URL}/${BINARY_NAME}" -o "${TEMP_DIR}/${BINARY_NAME}"
curl -fsSL "${DOWNLOAD_BASE_URL}/checksums.txt" -o "${TEMP_DIR}/checksums.txt"

EXPECTED_CHECKSUM="$(awk -v filename="$BINARY_NAME" '$2 == filename { print $1 }' "${TEMP_DIR}/checksums.txt")"
if [[ -z "$EXPECTED_CHECKSUM" ]]; then
    echo "Error: checksum for ${BINARY_NAME} was not found." >&2
    exit 1
fi

printf '%s  %s\n' "$EXPECTED_CHECKSUM" "${TEMP_DIR}/${BINARY_NAME}" | sha256sum --check --status

install -d -m 0755 "$BINARY_DIR"
install -m 0755 "${TEMP_DIR}/${BINARY_NAME}" "$BINARY_PATH"
install -d -m 0700 "$CONFIG_DIR"
install -d -m 0755 "$SERVICE_DIR"

cat >"$SERVICE_FILE" <<EOF
[Unit]
Description=Podlet Server Agent

[Service]
Type=simple
ExecStart="${BINARY_PATH}"
Restart=on-failure
RestartSec=5
NoNewPrivileges=true
PrivateTmp=true

[Install]
WantedBy=default.target
EOF
chmod 0644 "$SERVICE_FILE"
systemctl --user daemon-reload

echo "Registering server with Podlet..."
REGISTER_BODY="$(jq -n \
    --arg name "$SERVER_NAME" \
    --arg token "$REGISTRATION_TOKEN" \
    '{name: $name, registration_token: $token}')"

HTTP_STATUS="$(curl -sS \
    -o "${TEMP_DIR}/registration.json" \
    -w '%{http_code}' \
    -X POST "${API_URL}/servers/register" \
    -H 'Content-Type: application/json' \
    --data "$REGISTER_BODY" || true)"

if [[ ! "$HTTP_STATUS" =~ ^2[0-9][0-9]$ ]]; then
    echo "Error: server registration failed with HTTP ${HTTP_STATUS:-unknown}." >&2
    if [[ -s "${TEMP_DIR}/registration.json" ]]; then
        jq -r '.message // "Unknown API error"' "${TEMP_DIR}/registration.json" >&2 2>/dev/null || true
    fi
    exit 1
fi

if ! jq -e '
    (.id | type == "string" and length > 0) and
    (.token | type == "string" and length > 0)
' "${TEMP_DIR}/registration.json" >/dev/null; then
    echo "Error: Podlet returned an invalid registration response." >&2
    exit 1
fi

jq --arg api_url "$API_URL" \
    '{id, token, api_url: $api_url}' \
    "${TEMP_DIR}/registration.json" >"${TEMP_DIR}/server.json"
install -m 0600 "${TEMP_DIR}/server.json" "$CONFIG_FILE"

unset REGISTRATION_TOKEN REGISTER_BODY

echo "Starting Podlet agent..."
systemctl --user enable --now podlet-agent.service

if ! systemctl --user is-active --quiet podlet-agent.service; then
    echo "Error: the Podlet agent failed to start." >&2
    journalctl --user -u podlet-agent.service --no-pager -n 20 >&2 || true
    exit 1
fi

LINGER_ENABLED="false"
if command -v loginctl >/dev/null 2>&1; then
    loginctl enable-linger "$USER" >/dev/null 2>&1 || true
    LINGER_ENABLED="$(loginctl show-user "$USER" --property=Linger --value 2>/dev/null || true)"
fi

SERVER_ID="$(jq -r '.id' "${TEMP_DIR}/registration.json")"
echo "Podlet agent installed successfully."
echo "Server ID: ${SERVER_ID}"
echo "Status: systemctl --user status podlet-agent"
echo "Logs: journalctl --user -u podlet-agent -f"

if [[ "$LINGER_ENABLED" != "yes" ]]; then
    echo "Note: this system did not allow login lingering. The agent will run while your user session is active."
fi
