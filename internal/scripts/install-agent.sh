#!/usr/bin/env bash
set -euo pipefail

readonly API_URL="https://api.podlet.dev"
readonly AGENT_VERSION="v0.1.0"
readonly DOWNLOAD_ROOT_URL="https://get.podlet.dev/agent"
readonly DOWNLOAD_BASE_URL="${DOWNLOAD_ROOT_URL}/${AGENT_VERSION}"

readonly BINARY_DIR="/usr/local/bin"
readonly BINARY_PATH="${BINARY_DIR}/podlet-agent"

readonly DATA_DIR="/var/lib/podlet"
readonly UNINSTALLER_PATH="${DATA_DIR}/uninstall-agent.sh"

readonly CONFIG_DIR="/etc/podlet"
readonly CONFIG_FILE="${CONFIG_DIR}/server.json"

readonly SERVICE_DIR="/etc/systemd/system"
readonly SERVICE_FILE="${SERVICE_DIR}/podlet-agent.service"

readonly TOTAL_STEPS=7
CURRENT_STEP=0

# --------------------------------------------------
# Output helpers
# --------------------------------------------------

if [[ -t 1 && -z "${NO_COLOR:-}" ]]; then
    readonly CYAN=$'\033[1;36m'
    readonly GREEN=$'\033[1;32m'
    readonly RED=$'\033[1;31m'
    readonly DIM=$'\033[2m'
    readonly RESET=$'\033[0m'
else
    readonly CYAN=""
    readonly GREEN=""
    readonly RED=""
    readonly DIM=""
    readonly RESET=""
fi

step() {
    CURRENT_STEP=$((CURRENT_STEP + 1))
    printf '\n%s[%d/%d]%s %s\n' \
        "$CYAN" "$CURRENT_STEP" "$TOTAL_STEPS" "$RESET" "$1"
}

success() {
    printf '  %s✓%s %s\n' "$GREEN" "$RESET" "$1"
}

info() {
    printf '  %s%s%s\n' "$DIM" "$1" "$RESET"
}

error() {
    printf '  %s✗%s %s\n' "$RED" "$RESET" "$1" >&2
}

die() {
    error "$1"
    exit 1
}

download() {
    local url="$1"
    local destination="$2"

    if [[ -t 2 ]]; then
        curl -fL --retry 3 --progress-bar \
            "$url" -o "$destination"
    else
        curl -fsSL --retry 3 \
            "$url" -o "$destination"
    fi
}

printf '\n%sPodlet Server Agent Installer%s\n' "$CYAN" "$RESET"
info "Version ${AGENT_VERSION}"

# --------------------------------------------------
# 1. System compatibility
# --------------------------------------------------

step "Checking system compatibility"

if [[ "${EUID}" -ne 0 ]]; then
    die "This installer must run as root. Use: curl -fsSL ${DOWNLOAD_ROOT_URL}/install.sh | sudo bash"
fi

export HOME="/root"
export XDG_CONFIG_HOME="/etc"
export XDG_DATA_HOME="/var/lib"

if [[ "$(uname -s)" != "Linux" ]]; then
    die "Podlet currently supports Linux only."
fi

if [[ ! -r /etc/os-release ]]; then
    die "Cannot determine Linux distribution."
fi

# shellcheck source=/dev/null
source /etc/os-release

case "${ID:-}" in
    ubuntu | debian)
        success "Detected ${PRETTY_NAME:-$ID}"
        ;;
    *)
        die "Podlet currently supports Ubuntu and Debian only."
        ;;
esac

if [[ ! -d /run/systemd/system ]]; then
    die "Podlet requires systemd as the init system."
fi

success "systemd is available"

case "$(uname -m)" in
    x86_64 | amd64)
        ARCH="amd64"
        ;;
    aarch64 | arm64)
        ARCH="arm64"
        ;;
    *)
        die "Unsupported architecture: $(uname -m)"
        ;;
esac

readonly BINARY_NAME="podlet-agent-linux-${ARCH}"

success "Architecture: ${ARCH}"

# --------------------------------------------------
# 2. Install prerequisites
# --------------------------------------------------

step "Installing prerequisites"

install_prerequisites() {
    local cmd
    local -a packages=()
    local -a missing=()

    # Commands used by the installer.
    local -a commands=(
        curl
        git
        jq
        sha256sum
        systemctl
        install
        awk
    )

    for cmd in "${commands[@]}"; do
        if ! command -v "$cmd" >/dev/null 2>&1; then
            missing+=("$cmd")

            case "$cmd" in
                curl)      packages+=("curl") ;;
                git)       packages+=("git") ;;
                jq)        packages+=("jq") ;;
                sha256sum) packages+=("coreutils") ;;
                systemctl) packages+=("systemd") ;;
                install)   packages+=("coreutils") ;;
                awk)       packages+=("gawk") ;;
            esac
        fi
    done

    if [[ ${#missing[@]} -eq 0 ]]; then
        success "All prerequisites are already installed"
        return 0
    fi

    info "Missing: ${missing[*]}"

    if ! command -v apt-get >/dev/null 2>&1; then
        die "apt-get is required to install prerequisites."
    fi

    # Remove duplicate package names.
    local -A seen=()
    local -a unique_packages=()
    local package

    for package in "${packages[@]}"; do
        if [[ -z "${seen[$package]:-}" ]]; then
            unique_packages+=("$package")
            seen[$package]=1
        fi
    done

    info "Installing: ${unique_packages[*]}"

    export DEBIAN_FRONTEND=noninteractive

    apt-get update
    apt-get install -y --no-install-recommends "${unique_packages[@]}"

    for cmd in "${commands[@]}"; do
        if ! command -v "$cmd" >/dev/null 2>&1; then
            die "Failed to install prerequisite: ${cmd}"
        fi
    done

    success "All prerequisites installed successfully"
}

install_prerequisites

# --------------------------------------------------
# 3. Server configuration
# --------------------------------------------------

step "Configuring server"

if [[ ! -r /dev/tty || ! -w /dev/tty ]]; then
    die "An interactive terminal is required for registration."
fi

read -r -p "  Server name: " SERVER_NAME </dev/tty

if [[ -z "${SERVER_NAME//[[:space:]]/}" ]]; then
    die "Server name is required."
fi

read -r -p "  Registration token: " REGISTRATION_TOKEN </dev/tty

if [[ -z "$REGISTRATION_TOKEN" ]]; then
    die "Registration token is required."
fi

success "Server configuration complete"

# --------------------------------------------------
# 4. Download agent
# --------------------------------------------------

step "Downloading Podlet agent"

TEMP_DIR="$(mktemp -d)"
chmod 0700 "$TEMP_DIR"

trap 'rm -rf "$TEMP_DIR"' EXIT

info "Downloading agent ${AGENT_VERSION} (${ARCH})..."

download \
    "${DOWNLOAD_BASE_URL}/${BINARY_NAME}" \
    "${TEMP_DIR}/${BINARY_NAME}"

success "Agent binary downloaded"

info "Downloading checksum manifest..."

download \
    "${DOWNLOAD_BASE_URL}/checksums.txt" \
    "${TEMP_DIR}/checksums.txt"

success "Checksum manifest downloaded"

info "Downloading uninstaller..."

download \
    "${DOWNLOAD_ROOT_URL}/uninstall-agent.sh" \
    "${TEMP_DIR}/uninstall-agent.sh"

success "Uninstaller downloaded"

# --------------------------------------------------
# 5. Verify download
# --------------------------------------------------

step "Verifying download integrity"

EXPECTED_CHECKSUM="$(
    awk -v filename="$BINARY_NAME" \
        '$2 == filename { print $1 }' \
        "${TEMP_DIR}/checksums.txt"
)"

if [[ ! "$EXPECTED_CHECKSUM" =~ ^[a-fA-F0-9]{64}$ ]]; then
    die "Valid checksum for ${BINARY_NAME} was not found."
fi

if ! printf '%s  %s\n' \
    "$EXPECTED_CHECKSUM" \
    "${TEMP_DIR}/${BINARY_NAME}" \
    | sha256sum --check --status; then
    die "Agent checksum verification failed."
fi

success "SHA-256 checksum verified"

# --------------------------------------------------
# 6. Register and install server
# --------------------------------------------------

step "Registering and installing server"

info "Registering server with Podlet..."

jq -n \
    --arg name "$SERVER_NAME" \
    --arg token "$REGISTRATION_TOKEN" \
    '{
        name: $name,
        registration_token: $token
    }' >"${TEMP_DIR}/register-request.json"

unset REGISTRATION_TOKEN

if ! HTTP_STATUS="$(
    curl -sS \
        --retry 3 \
        --connect-timeout 15 \
        -o "${TEMP_DIR}/registration.json" \
        -w '%{http_code}' \
        -X POST "${API_URL}/servers/register" \
        -H 'Content-Type: application/json' \
        --data-binary @"${TEMP_DIR}/register-request.json"
)"; then
    die "Could not connect to Podlet API."
fi

if [[ ! "$HTTP_STATUS" =~ ^2[0-9][0-9]$ ]]; then
    if [[ -s "${TEMP_DIR}/registration.json" ]]; then
        jq -r \
            '.message // "Unknown API error"' \
            "${TEMP_DIR}/registration.json" >&2 2>/dev/null || true
    fi

    die "Server registration failed with HTTP ${HTTP_STATUS}."
fi

if ! jq -e '
    (.id | type == "string" and length > 0) and
    (.token | type == "string" and length > 0)
' "${TEMP_DIR}/registration.json" >/dev/null; then
    die "Podlet returned an invalid registration response."
fi

SERVER_ID="$(jq -r '.id' "${TEMP_DIR}/registration.json")"

success "Server registered: ${SERVER_ID}"

info "Installing agent binary..."

install -d -m 0755 "$BINARY_DIR"
install -d -m 0755 "$DATA_DIR"
install -d -m 0700 "$CONFIG_DIR"
install -d -m 0755 "$SERVICE_DIR"

install -m 0755 \
    "${TEMP_DIR}/${BINARY_NAME}" \
    "$BINARY_PATH"

install -m 0755 \
    "${TEMP_DIR}/uninstall-agent.sh" \
    "$UNINSTALLER_PATH"

success "Agent binary installed"

info "Saving server configuration..."

jq --arg api_url "$API_URL" \
    '{
        id,
        token,
        api_url: $api_url
    }' \
    "${TEMP_DIR}/registration.json" \
    >"${TEMP_DIR}/server.json"

install -m 0600 \
    "${TEMP_DIR}/server.json" \
    "$CONFIG_FILE"

success "Server configuration saved"

# --------------------------------------------------
# 7. Initialize and start agent
# --------------------------------------------------

step "Initializing and starting agent"

info "Running agent initialization..."
info "This includes Docker and other required components."
echo ""

# The agent handles its own progress output.
if ! "$BINARY_PATH" --initialize-only; then
    die "Agent initialization failed."
fi

echo ""
success "Agent initialization complete"

info "Configuring systemd service..."

cat >"$SERVICE_FILE" <<EOF
[Unit]
Description=Podlet Server Agent
After=network-online.target docker.service
Wants=network-online.target

[Service]
Type=simple
User=root
Group=root
Environment=HOME=/root
Environment=XDG_CONFIG_HOME=/etc
Environment=XDG_DATA_HOME=/var/lib
ExecStart=${BINARY_PATH}
Restart=on-failure
RestartSec=5
NoNewPrivileges=true
PrivateTmp=true

[Install]
WantedBy=multi-user.target
EOF

chmod 0644 "$SERVICE_FILE"

systemctl daemon-reload

success "Systemd service configured"

info "Enabling and starting Podlet agent..."

systemctl enable podlet-agent.service
systemctl restart podlet-agent.service

if ! systemctl is-active --quiet podlet-agent.service; then
    journalctl \
        -u podlet-agent.service \
        --no-pager \
        -n 20 >&2 || true

    die "Podlet agent failed to start."
fi

success "Podlet agent is running"

# --------------------------------------------------
# Installation summary
# --------------------------------------------------

echo ""
printf '%s========================================%s\n' "$GREEN" "$RESET"
printf '%s  Podlet installed successfully!%s\n' "$GREEN" "$RESET"
printf '%s========================================%s\n' "$GREEN" "$RESET"
echo ""
printf '  Server ID:     %s\n' "$SERVER_ID"
printf '  Agent version: %s\n' "$AGENT_VERSION"
printf '  Configuration: %s\n' "$CONFIG_FILE"
printf '  Uninstaller:   %s\n' "$UNINSTALLER_PATH"
echo ""
printf '  Status: systemctl status podlet-agent\n'
printf '  Logs:   journalctl -u podlet-agent -f\n'
echo ""
