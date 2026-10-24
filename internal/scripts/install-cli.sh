#!/usr/bin/env bash
set -euo pipefail

readonly CLI_VERSION="v0.1.0"
readonly DOWNLOAD_ROOT_URL="https://downloads.podlet.dev/cli"
readonly DOWNLOAD_BASE_URL="${DOWNLOAD_ROOT_URL}/${CLI_VERSION}"
readonly API_URL="https://api.podlet.dev"
readonly BINARY_DIR="${HOME}/.local/bin"
readonly BINARY_PATH="${BINARY_DIR}/podlet"
readonly DATA_DIR="${XDG_DATA_HOME:-${HOME}/.local/share}/podlet"
readonly UNINSTALLER_PATH="${DATA_DIR}/uninstall-cli.sh"
readonly USER_CONFIG_DIR="${XDG_CONFIG_HOME:-${HOME}/.config}"
readonly CONFIG_DIR="${USER_CONFIG_DIR}/podlet"
readonly CONFIG_FILE="${CONFIG_DIR}/config.json"

if [[ "${EUID}" -eq 0 ]]; then
    echo "Error: run this installer without sudo or root." >&2
    echo "Use: curl -fsSL https://downloads.podlet.dev/cli/install.sh | bash" >&2
    exit 1
fi

if [[ "$(uname -s)" != "Linux" ]]; then
    echo "Error: the Podlet CLI currently supports Linux only." >&2
    exit 1
fi

for command in curl sha256sum install awk; do
    if ! command -v "$command" >/dev/null 2>&1; then
        echo "Error: required command not found: $command" >&2
        exit 1
    fi
done

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

readonly BINARY_NAME="podlet-linux-${ARCH}"
TEMP_DIR="$(mktemp -d)"
trap 'rm -rf "$TEMP_DIR"' EXIT

echo "Downloading Podlet CLI ${CLI_VERSION}..."
curl -fsSL "${DOWNLOAD_BASE_URL}/${BINARY_NAME}" -o "${TEMP_DIR}/${BINARY_NAME}"
curl -fsSL "${DOWNLOAD_BASE_URL}/checksums.txt" -o "${TEMP_DIR}/checksums.txt"
curl -fsSL "${DOWNLOAD_ROOT_URL}/uninstall-cli.sh" -o "${TEMP_DIR}/uninstall-cli.sh"

EXPECTED_CHECKSUM="$(awk -v filename="$BINARY_NAME" '$2 == filename { print $1 }' "${TEMP_DIR}/checksums.txt")"
if [[ -z "$EXPECTED_CHECKSUM" ]]; then
    echo "Error: checksum for ${BINARY_NAME} was not found." >&2
    exit 1
fi

printf '%s  %s\n' "$EXPECTED_CHECKSUM" "${TEMP_DIR}/${BINARY_NAME}" | sha256sum --check --status

install -d -m 0755 "$BINARY_DIR"
install -m 0755 "${TEMP_DIR}/${BINARY_NAME}" "$BINARY_PATH"
install -d -m 0755 "$DATA_DIR"
install -m 0755 "${TEMP_DIR}/uninstall-cli.sh" "$UNINSTALLER_PATH"
install -d -m 0700 "$CONFIG_DIR"

if [[ ! -e "$CONFIG_FILE" ]]; then
    printf '{\n  "base_url": "%s"\n}\n' "$API_URL" >"${TEMP_DIR}/config.json"
    install -m 0600 "${TEMP_DIR}/config.json" "$CONFIG_FILE"
fi

echo "Podlet CLI installed successfully."
echo "Binary: ${BINARY_PATH}"
echo "Config: ${CONFIG_FILE}"
echo "Uninstaller: ${UNINSTALLER_PATH}"

case ":${PATH}:" in
    *":${BINARY_DIR}:"*)
        echo "Run: podlet --help"
        ;;
    *)
        echo "Add ${BINARY_DIR} to your PATH, then run: podlet --help"
        echo "For bash: echo 'export PATH=\"\$HOME/.local/bin:\$PATH\"' >> ~/.bashrc && source ~/.bashrc"
        ;;
esac
