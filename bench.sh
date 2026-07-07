#!/bin/bash
set -e 

# default vars
BENCHMODE="single"
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
OUTPUT_DIR="$SCRIPT_DIR/results/yaml"
PLOTS_DIR="$SCRIPT_DIR/results/plots"

# Safely deleting runs
ALLOWED_DELETE_DIR="$SCRIPT_DIR/results/yaml"

safe_rm() {
    local TARGET="$1"
    
    if [[ -z "${TARGET// }" ]]; then
        echo "safe_rm: BLOCKED -- empty path" >&2
        return 1
    fi
    
    local ABS_TARGET="$(realpath "$TARGET" 2>/dev/null)"
    if [[ -z "$ABS_TARGET" ]]; then
        echo "safe_rm: cannot resolve path: $TARGET" >&2
        return 1
    fi
    
    if [[ "$ABS_TARGET" != "$ALLOWED_DELETE_DIR/"* && "$ABS_TARGET" != "$ALLOWED_DELETE_DIR" ]]; then
        echo "safe_rm: BLOCKED -- $TARGET is not inside $ALLOWED_DELETE_DIR" >&2
        return 1
    fi
    
    rm -rf "$ABS_TARGET"
}

# --mode parse
while [[ $# -gt 0 ]]; do 
    case $1 in
        --mode)
            BENCHMODE="$2"
            shift 2
            ;;
        *)
            echo "Invalid argument: $1" >&2
            exit 1
            ;;
    esac
done

if [[ "$BENCHMODE" != "single" && "$BENCHMODE" != "compare" ]]; then
    echo "Error: --mode must be single or compare" >&2
    exit 1
fi

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

# option 1 | Use .env
if [[ "$UsingEnv" == true ]]; then
    if [[ ! -f "$SCRIPT_DIR/.env" ]]; then
        echo "No .env found. Copy .env.example to .env and fill it in first." >&2
        exit 1
    fi
    set -a
    source "$SCRIPT_DIR/.env"
    set +a
fi

for FILE in "$ALLOWED_DELETE_DIR"/*; do
    [[ -e "$FILE" ]] && safe_rm "$FILE"
done

# sleep/cooldown function
msleep() {
    local Ms=${COOLDOWNTIMER:-500}
    local Sec=$((Ms / 1000))
    local Rem=$((Ms % 1000))
    sleep "$(printf "%d.%03d" "$Sec" "$Rem")"
}

# TODO: gate the exegesis mode at the start so that it can do latency + uop.
if [[ "$UsingEnv" == true ]]; then

    if [[ "$BENCHMODE" == "single" ]]; then
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
        
        if [[ "$METHODOLOGY" == "cycling" ]]; then
            for i in $(seq 1 $RUNS); do
                "$EXEGESIS_BINA" \
                    --mode="$EXEGESIS_MODEA" \
                    --opcode-name="$OPCODEA" \
                    --mcpu="$MCPUA" \
                    --benchmarks-file="$OUTPUT_DIR/Arun_$i.yaml"
                echo "A run $i/$RUNS complete"
                msleep

                "$EXEGESIS_BINB" \
                    --mode="$EXEGESIS_MODEB" \
                    --opcode-name="$OPCODEB" \
                    --mcpu="$MCPUB" \
                    --benchmarks-file="$OUTPUT_DIR/Brun_$i.yaml"
                echo "B run $i/$RUNS complete"
                msleep
            done

        elif [[ "$METHODOLOGY" == "random interleaving" ]]; then

            COUNTA=1
            COUNTB=1
            TOTALRUNS=$((RUNS * 2))

            while read -r choice; do 
                if [[ "$choice" == 0 ]]; then
                    "$EXEGESIS_BINA" \
                        --mode="$EXEGESIS_MODEA" \
                        --opcode-name="$OPCODEA" \
                        --mcpu="$MCPUA" \
                        --benchmarks-file="$OUTPUT_DIR/Arun_$COUNTA.yaml"
                    echo "A run $COUNTA/$RUNS complete"
                    msleep
                    ((COUNTA++))
                else
                    "$EXEGESIS_BINB" \
                        --mode="$EXEGESIS_MODEB" \
                        --opcode-name="$OPCODEB" \
                        --mcpu="$MCPUB" \
                        --benchmarks-file="$OUTPUT_DIR/Brun_$COUNTB.yaml"
                    echo "B run $COUNTB/$RUNS complete"
                    msleep
                    ((COUNTB++))
                fi 
            done < <(shuf -i 1-"$TOTALRUNS" | while read -r num; do echo $((num % 2)); done)

        else
            for i in $(seq 1 $RUNS); do
                "$EXEGESIS_BINA" \
                    --mode="$EXEGESIS_MODEA" \
                    --opcode-name="$OPCODEA" \
                    --mcpu="$MCPUA" \
                    --benchmarks-file="$OUTPUT_DIR/Arun_$i.yaml"
                echo "A run $i/$RUNS complete"
                msleep
            done
            
            echo "2 second cooldown between comparison" 
            sleep 2

            for i in $(seq 1 $RUNS); do
                "$EXEGESIS_BINB" \
                    --mode="$EXEGESIS_MODEB" \
                    --opcode-name="$OPCODEB" \
                    --mcpu="$MCPUB" \
                    --benchmarks-file="$OUTPUT_DIR/Brun_$i.yaml"
                echo "B run $i/$RUNS complete"
                msleep
            done
        fi
    fi
fi

# TODO: Complete section. Currently only supports sequential METHODOLOGY.
# option 2 goes here.
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
if [[ "$BENCHMODE" == "compare" ]]; then
    if [ -n "$LABELA" ] && [ -n "$LABELB" ]; then
        LABELS="${LABELA}|${LABELB}"
    else
        LABELS="A|B"
    fi
else
    LABELS="${LABEL:-Run}"
fi

# Python cmd - should be given accurate args.
python3 analyze.py "${BENCHMODE}" "${RUNS}" "${LABELS}" "${METHODOLOGY:-sequential}" "${COOLDOWNTIMER:-500}"
