#!/usr/bin/env bash
set -euo pipefail

readonly CLI_VERSION="v0.1.0"
readonly DOWNLOAD_ROOT_URL="https://get.podlet.dev/cli"
readonly DOWNLOAD_BASE_URL="${DOWNLOAD_ROOT_URL}/${CLI_VERSION}"
readonly API_URL="https://api.podlet.dev"

readonly BINARY_DIR="${HOME}/.local/bin"
readonly BINARY_PATH="${BINARY_DIR}/podlet"

readonly DATA_DIR="${XDG_DATA_HOME:-${HOME}/.local/share}/podlet"
readonly UNINSTALLER_PATH="${DATA_DIR}/uninstall-cli.sh"

readonly USER_CONFIG_DIR="${XDG_CONFIG_HOME:-${HOME}/.config}"
readonly CONFIG_DIR="${USER_CONFIG_DIR}/podlet"
readonly CONFIG_FILE="${CONFIG_DIR}/config.json"

readonly TOTAL_STEPS=6
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

printf '\n%sPodlet CLI Installer%s\n' "$CYAN" "$RESET"
info "Version ${CLI_VERSION}"

# --------------------------------------------------
# 1. System compatibility
# --------------------------------------------------

step "Checking system compatibility"

if [[ "${EUID}" -eq 0 ]]; then
    die "Run without sudo or root. Use: curl -fsSL ${DOWNLOAD_ROOT_URL}/install.sh | bash"
fi

if [[ "$(uname -s)" != "Linux" ]]; then
    die "Podlet CLI currently supports Linux only."
fi

success "Linux detected"

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

readonly BINARY_NAME="podlet-linux-${ARCH}"

success "Architecture: ${ARCH}"

# --------------------------------------------------
# 2. Verify prerequisites
# --------------------------------------------------

step "Checking prerequisites"

for cmd in curl sha256sum install awk grep mktemp; do
    if ! command -v "$cmd" >/dev/null 2>&1; then
        die "Required command not found: ${cmd}"
    fi
done

success "All prerequisites are available"

# --------------------------------------------------
# 3. Download CLI
# --------------------------------------------------

step "Downloading Podlet CLI"

TEMP_DIR="$(mktemp -d)"
chmod 0700 "$TEMP_DIR"

trap 'rm -rf "$TEMP_DIR"' EXIT

info "Downloading Podlet CLI ${CLI_VERSION} (${ARCH})..."

download \
    "${DOWNLOAD_BASE_URL}/${BINARY_NAME}" \
    "${TEMP_DIR}/${BINARY_NAME}"

success "CLI binary downloaded"

info "Downloading checksum manifest..."

download \
    "${DOWNLOAD_BASE_URL}/checksums.txt" \
    "${TEMP_DIR}/checksums.txt"

success "Checksum manifest downloaded"

info "Downloading uninstaller..."

download \
    "${DOWNLOAD_ROOT_URL}/uninstall-cli.sh" \
    "${TEMP_DIR}/uninstall-cli.sh"

success "Uninstaller downloaded"

# --------------------------------------------------
# 4. Verify download
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
    die "CLI checksum verification failed."
fi

success "SHA-256 checksum verified"

# --------------------------------------------------
# 5. Install CLI
# --------------------------------------------------

step "Installing Podlet CLI"

info "Installing binary..."

install -d -m 0755 "$BINARY_DIR"

install -m 0755 \
    "${TEMP_DIR}/${BINARY_NAME}" \
    "$BINARY_PATH"

success "CLI binary installed"

info "Installing uninstaller..."

install -d -m 0755 "$DATA_DIR"

install -m 0755 \
    "${TEMP_DIR}/uninstall-cli.sh" \
    "$UNINSTALLER_PATH"

success "Uninstaller installed"

info "Configuring Podlet CLI..."

install -d -m 0700 "$CONFIG_DIR"

if [[ ! -e "$CONFIG_FILE" ]]; then
    printf '{\n  "base_url": "%s"\n}\n' \
        "$API_URL" >"${TEMP_DIR}/config.json"

    install -m 0600 \
        "${TEMP_DIR}/config.json" \
        "$CONFIG_FILE"

    success "Default configuration created"
else
    success "Existing configuration preserved"
fi

# --------------------------------------------------
# 6. Configure PATH
# --------------------------------------------------

step "Configuring PATH"

configure_path() {
    local shell_rc
    local path_entry
    local shell_name

    if [[ ":${PATH}:" == *":${BINARY_DIR}:"* ]]; then
        success "Podlet binary directory is already in PATH"
        return 0
    fi

    shell_name="${SHELL:-/bin/bash}"

    case "$shell_name" in
        */zsh)
            shell_rc="${HOME}/.zshrc"
            ;;
        */bash)
            shell_rc="${HOME}/.bashrc"
            ;;
        *)
            shell_rc="${HOME}/.profile"
            ;;
    esac

    path_entry='export PATH="$HOME/.local/bin:$PATH"'

    # Create the shell configuration file if needed.
    touch "$shell_rc"

    # Avoid duplicate PATH entries.
    if ! grep -Fxq "$path_entry" "$shell_rc"; then
        printf '\n%s\n' "$path_entry" >>"$shell_rc"
        success "Added ${BINARY_DIR} to PATH in ${shell_rc}"
    else
        success "PATH entry already exists in ${shell_rc}"
    fi

    info "Open a new terminal to use the updated PATH."
}

configure_path

# --------------------------------------------------
# Installation summary
# --------------------------------------------------

echo ""
printf '%s========================================%s\n' "$GREEN" "$RESET"
printf '%s  Podlet CLI installed successfully!%s\n' "$GREEN" "$RESET"
printf '%s========================================%s\n' "$GREEN" "$RESET"
echo ""

printf '  Version:     %s\n' "$CLI_VERSION"
printf '  Binary:      %s\n' "$BINARY_PATH"
printf '  Config:      %s\n' "$CONFIG_FILE"
printf '  Uninstaller: %s\n' "$UNINSTALLER_PATH"

echo ""

if [[ ":${PATH}:" == *":${BINARY_DIR}:"* ]]; then
    printf '  Run: podlet --help\n'
else
    printf '  Open a new terminal, then run: podlet --help\n'
fi

echo ""
