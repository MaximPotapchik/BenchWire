# shellcheck shell=bash

# TODO: Perhaps a wider search, or .env source customization.
load_env() {
    if [[ ! -f "$SCRIPT_DIR/.env" ]]; then
        echo "No .env found. Copy .env.example to .env and fill it in first." >&2
        exit 1
    fi
    set -a
    source "$SCRIPT_DIR/.env"
    set +a
}
