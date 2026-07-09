# shellcheck shell=bash

#TODO: Add timing calibration. Ramp ups, and randomization.
msleep() {
    local Ms=${1:-${COOLDOWNTIMER:-500}}
    local Sec=$((Ms / 1000))
    local Rem=$((Ms % 1000))
    sleep "$(printf "%d.%03d" "$Sec" "$Rem")"
}
