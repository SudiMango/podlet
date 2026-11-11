
#!/usr/bin/env bash
set -euo pipefail

readonly BINARY_PATH="${HOME}/.local/bin/podlet"
readonly DATA_DIR="${XDG_DATA_HOME:-${HOME}/.local/share}/podlet"
readonly UNINSTALLER_PATH="${DATA_DIR}/uninstall-cli.sh"

readonly USER_CONFIG_DIR="${XDG_CONFIG_HOME:-${HOME}/.config}"
readonly CONFIG_DIR="${USER_CONFIG_DIR}/podlet"
readonly CONFIG_FILE="${CONFIG_DIR}/config.json"

readonly TOTAL_STEPS=4
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

printf '\n%sPodlet CLI Uninstaller%s\n' "$CYAN" "$RESET"

# --------------------------------------------------
# 1. Verify installation
# --------------------------------------------------

step "Checking Podlet installation"

if [[ "${EUID}" -eq 0 ]]; then
    die "Run this uninstaller without sudo or root."
fi

if [[ ! -x "$BINARY_PATH" ]]; then
    die "Podlet CLI was not found at ${BINARY_PATH}. Nothing was removed."
fi

success "Podlet CLI installation found"

# --------------------------------------------------
# 2. Log out
# --------------------------------------------------

step "Logging out of Podlet"

info "Removing saved authentication credentials..."

if ! "$BINARY_PATH" auth logout; then
    die "Logout failed. Nothing was removed."
fi

success "Logged out successfully"

# --------------------------------------------------
# 3. Remove CLI files
# --------------------------------------------------

step "Removing Podlet CLI"

rm -f "$BINARY_PATH"
success "CLI binary removed"

rm -f "$CONFIG_FILE"
rmdir "$CONFIG_DIR" 2>/dev/null || true
success "CLI configuration removed"

# --------------------------------------------------
# 4. Clean up
# --------------------------------------------------

step "Cleaning up installation"

rm -f "$UNINSTALLER_PATH"
rmdir "$DATA_DIR" 2>/dev/null || true

success "Uninstaller removed"

# --------------------------------------------------
# Uninstallation summary
# --------------------------------------------------

echo ""
printf '%s========================================%s\n' "$GREEN" "$RESET"
printf '%s  Podlet CLI uninstalled successfully!%s\n' "$GREEN" "$RESET"
printf '%s========================================%s\n' "$GREEN" "$RESET"
echo ""

info "If your shell still resolves 'podlet', run: hash -r"
echo ""
