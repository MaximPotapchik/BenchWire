# shellcheck shell=bash

safe_rm() {
    local TARGET="$1"
    
    if [[ -z "${TARGET// }" ]]; then
        echo "safe_rm: BLOCKED -- empty path" >&2
        return 1
    fi
    
    local ABS_TARGET
    ABS_TARGET="$(realpath "$TARGET" 2>/dev/null)"
    if [[ -z "$ABS_TARGET" ]]; then
        echo "safe_rm: cannot resolve path: $TARGET" >&2
        return 1
    fi
    
    if [[ "$ABS_TARGET" != "$ALLOWED_DELETE_DIR/"* && "$ABS_TARGET" != "$ALLOWED_DELETE_DIR" ]]; then
        echo "safe_rm: BLOCKED - $TARGET is not inside $ALLOWED_DELETE_DIR" >&2
        return 1
    fi
    
    rm -rf "$ABS_TARGET"
}

