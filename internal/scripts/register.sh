
#!/usr/bin/env bash
set -euo pipefail

API_URL="${PODLET_API_URL:-http://localhost:8080}"
CONFIG_DIR="${XDG_CONFIG_HOME:-$HOME/.config}/podlet"
CONFIG_FILE="$CONFIG_DIR/server.json"

NAME=""
TOKEN=""

while [[ $# -gt 0 ]]; do
    case "$1" in
        --name)
            NAME="${2:?Missing name}"
            shift 2
            ;;
        --token)
            TOKEN="${2:?Missing token}"
            shift 2
            ;;
        *)
            echo "Unknown argument: $1" >&2
            exit 1
            ;;
    esac
done

if [[ -z "$NAME" || -z "$TOKEN" ]]; then
    echo "Usage: $0 --name NAME --token TOKEN" >&2
    exit 1
fi

if ! command -v jq >/dev/null 2>&1; then
    echo "Error: jq is required. Install it with: sudo apt install jq" >&2
    exit 1
fi

# Make request and capture HTTP status separately
BODY_FILE=$(mktemp)
trap 'rm -f "$BODY_FILE"' EXIT

HTTP_STATUS=$(curl -sS \
    -o "$BODY_FILE" \
    -w "%{http_code}" \
    -X POST "$API_URL/servers/register" \
    -H "Content-Type: application/json" \
    -d "$(jq -n \
        --arg name "$NAME" \
        --arg token "$TOKEN" \
        '{name: $name, registration_token: $token}')")

if [[ "$HTTP_STATUS" -lt 200 || "$HTTP_STATUS" -ge 300 ]]; then
    echo "Registration failed (HTTP $HTTP_STATUS)" >&2
    echo "Response:" >&2
    if ! jq . "$BODY_FILE" >&2 2>/dev/null; then
        cat "$BODY_FILE" >&2
    fi
    exit 1
fi

# Validate response
if ! jq -e '
    (.id | type == "string") and
    (.name | type == "string") and
    (.token | type == "string" and length > 0)
' "$BODY_FILE" >/dev/null; then
    echo "Error: Invalid server response" >&2
    jq . "$BODY_FILE" >&2 || cat "$BODY_FILE" >&2
    exit 1
fi

# Save credentials
mkdir -p "$CONFIG_DIR"
chmod 700 "$CONFIG_DIR"
umask 077

jq '{id, token}' "$BODY_FILE" > "$CONFIG_FILE"
chmod 600 "$CONFIG_FILE"

echo "Server registered successfully!"
echo
jq '{id, name, user_id, created_at, updated_at}' "$BODY_FILE"
echo
echo "Credentials saved to: $CONFIG_FILE"
