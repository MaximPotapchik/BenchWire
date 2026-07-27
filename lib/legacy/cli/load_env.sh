# shellcheck shell=bash

declare -a RAWENV=()

# TODO: - Doesn't have to be strictly .env. Can go for yaml/toml.
# - More validation checks.
# - Potential try refactoring for this format:
# EXEGESIS_FLAGS="--mcpu=native --mode=latency --etc."
load_env() {
    # Grabs the entire .env and tosses it into an array.
    get_env_array() {
        while IFS= read -r RAWLINE || [[ -n "$RAWLINE" ]]; do
            if [[ -z "$RAWLINE" ]]; then
                continue
            fi

            if [[ "$RAWLINE" =~ ^[[:space:]]*# ]]; then
                continue
            fi
            
            if [[ "$RAWLINE" =~ ^-- ]]; then
                RAWENV+=( "$RAWLINE" )
            else
                eval "export $RAWLINE"
                RAWENV+=( "$RAWLINE" )
            fi
        done < "$SCRIPT_DIR/.env"
    }

    if [[ ! -f "$SCRIPT_DIR/.env" ]]; then
        echo "No .env found. Copy .env.example to .env and fill it in first." >&2
        exit 1
    fi

    get_env_array
}
