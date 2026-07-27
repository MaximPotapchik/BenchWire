#!/bin/bash
set -e

# default vars
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
OUTPUT_DIR="$SCRIPT_DIR/results/yaml"
PLOTS_DIR="$SCRIPT_DIR/results/plots"

# add allowed to delete directories here.
ALLOWED_DELETE_DIR="$SCRIPT_DIR/results/yaml"
. "$SCRIPT_DIR/lib/cli/safe_rm.sh"

mkdir -p "$OUTPUT_DIR" "$PLOTS_DIR"

# cmd vars
UsingEnv=false
UsingCmd=false

read -p "Warning: This will delete your previous runs yaml files
1 to use ENV | 2 for custom command: " Option

# If check for initial gate.
# TODO: - Add error reasons here.
if [[ "$Option" != "1" && "$Option" != "2" ]]; then 
    echo "Invalid option, exiting." >&2
    exit 1
elif [[ "$Option" == "1" ]]; then
    UsingEnv=true
elif [[ "$Option" == "2" ]]; then
    UsingCmd=true
fi 

# Delete old runs.
for FILE in "$ALLOWED_DELETE_DIR"/*; do
    [[ -e "$FILE" ]] && safe_rm "$FILE"
done

# sleep/cooldown function
# source=lib/cli/msleep.sh
source "$SCRIPT_DIR/lib/cli/msleep.sh"
# runner function
. "$SCRIPT_DIR/lib/cli/runner.sh"


# option 1 | Use .env
# TODO: 
# - Add a/b/b/a mode.
if [[ "$UsingEnv" == true ]]; then
    . "$SCRIPT_DIR/lib/cli/load_env.sh"
    . "$SCRIPT_DIR/lib/cli/exegesis_formatter.sh"
    load_env
    exegesis_formatter "$METHODOLOGY" RAWENV
    runner "$METHODOLOGY" "$RUNS" "$OUTPUT_DIR" FORMATTEDARRAY
fi

# TODO: Complete section. Currently only supports sequential METHODOLOGY.
# option 2 goes here
# This is currently broken..
if [[ "$UsingCmd" == true ]]; then 

    if [[ "$BENCHMODE" == "single" ]]; then 
        read -p "Enter Exegesis flags: " Flags
        eval set -- "$Flags"
        RUNS="$1"
        shift
        EXEGESIS_BIN="$1"
        shift

        while [[ $# -gt 0 ]]; do 
            case $1 in 
                --mode=*)
                    EXEGESIS_MODE="${1#*=}"
                    shift
                    ;;
                --opcode-name=*)
                    OPCODE="${1#*=}"
                    shift
                    ;;
                --mcpu=*)
                    MCPU="${1#*=}"
                    shift
                    ;;
                *)
                    echo "Unknown flag: $1" >&2
                    shift
                    ;;
            esac
        done

        for i in $(seq 1 $RUNS); do
            "$EXEGESIS_BIN" \
                --mode="$EXEGESIS_MODE" \
                --opcode-name="$OPCODE" \
                --mcpu="$MCPU" \
                --benchmarks-file="$OUTPUT_DIR/run_$i.yaml"
            echo "Run $i/$RUNS complete"
            msleep
        done
    fi

    if [[ "$BENCHMODE" == "compare" ]]; then
        read -p "Enter first batch Exegesis flags: " FlagsA
        eval set -- "$FlagsA"
        RUNS="$1"
        shift
        EXEGESIS_BIN="$1"
        shift

        while [[ $# -gt 0 ]]; do 
            case $1 in 
                --mode=*)
                    EXEGESIS_MODE="${1#*=}"
                    shift
                    ;;
                --opcode-name=*)
                    OPCODE="${1#*=}"
                    shift
                    ;;
                --mcpu=*)
                    MCPU="${1#*=}"
                    shift
                    ;;
                *)
                    echo "Unknown flag: $1" >&2
                    shift
                    ;;
            esac
        done

        for i in $(seq 1 $RUNS); do 
            "$EXEGESIS_BIN" \
                --mode="$EXEGESIS_MODE" \
                --opcode-name="$OPCODE" \
                --mcpu="$MCPU" \
                --benchmarks-file="$OUTPUT_DIR/run_$i.yaml"
            echo "Run $i/$RUNS complete"
            msleep
        done

        read -p "Enter first batch Exegesis flags: " FlagsB
        eval set -- "$FlagsB"
        RUNS="$1"
        shift
        EXEGESIS_BIN="$1"
        shift

        while [[ $# -gt 0 ]]; do 
            case $1 in 
                --mode=*)
                    EXEGESIS_MODE="${1#*=}"
                    shift
                    ;;
                --opcode-name=*)
                    OPCODE="${1#*=}"
                    shift
                    ;;
                --mcpu=*)
                    MCPU="${1#*=}"
                    shift
                    ;;
                *)
                    echo "Unknown flag: $1" >&2
                    shift
                    ;;
            esac
        done

        for i in $(seq 1 $RUNS); do 
            "$EXEGESIS_BIN" \
                --mode="$EXEGESIS_MODE" \
                --opcode-name="$OPCODE" \
                --mcpu="$MCPU" \
                --benchmarks-file="$OUTPUT_DIR/run_$i.yaml"
            echo "Run $i/$RUNS complete"
            msleep
        done
    fi
fi

# Label
if [[ "$BENCHMODE" != "single" ]]; then
    if [ -n "$LABELA" ] && [ -n "$LABELB" ]; then
        LABELS="${LABELA}|${LABELB}"
    else
        LABELS="A|B"
    fi
else
    LABELS="${LABEL:-Run}"
fi

# Python cmd - should be given accurate args.
python3 analyze.py "${METHODOLOGY}" "${RUNS}"  "${COOLDOWNTIMER:-500}" "${LABELS}"
