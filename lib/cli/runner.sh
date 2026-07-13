# shellcheck shell=bash
RUNCMDA=""
RUNCMDB=""
RUNCMD=""

# TODO: This can be further modularized. Should probably get its own folder.
runner() {
    local METHODOLOGY=$1
    local RUNS=$2
    local OUTPUTDIR=$3
    local -n ARGS=$4

    local COMMAND=""
    local COMMANDA=""
    local COMMANDB=""

    if [[ "${ARGS[0]}" == *\|* ]]; then
        COMMANDA+="${ARGS[0]%%|*}"
        COMMANDB+="${ARGS[0]##*|}"
    else
        COMMAND="${ARGS[0]}"
    fi

    if [[ "$METHODOLOGY" != "single" ]]; then
        local curI=1
        local total=${#ARGS[@]}

        while [[ $curI -lt $total ]]; do 
            local opt="${ARGS[$curI]}"
            local cli="${opt%%=*}="
            local val="${opt##*=}"
            # Check 2nd last letter for A. cli ends with = usually.
            local check="${cli: -2: 1}"

            if [[ "$check" == "A" ]]; then
                COMMANDA+=" ${cli%??}=${val}"
            else
                COMMANDB+=" ${cli%??}=${val}"
            fi 
            ((curI++))
        done
    else
        for opt in "${ARGS[@]:2}"; do
            COMMAND+=" $opt"
        done
    fi

    run_cmd() {
        local RUN=$1
        if [[ "$METHODOLOGY" != "single" ]]; then
            RUNCMDA="$COMMANDA --benchmarks-file=$OUTPUTDIR/Arun_$RUN.yaml"
            RUNCMDB="$COMMANDB --benchmarks-file=$OUTPUTDIR/Brun_$RUN.yaml"
        else
            RUNCMD="$COMMAND --benchmarks-file=$OUTPUTDIR/run_$RUN.yaml"
        fi
    }

    case "$METHODOLOGY" in
        "single")
            for i in $(seq 1 "$RUNS"); do 
                run_cmd "$i"
                $RUNCMD 
                msleep
            done
            ;;
        "sequential")
            for i in $(seq 1 $RUNS); do 
                run_cmd "$i"
                echo "$RUNCMDA"
                $RUNCMDA
                msleep
            done

            for i in $(seq 1 $RUNS); do 
                run_cmd "$i"
                $RUNCMDB
                msleep
            done
            ;;
        "cycling")
            for i in $(seq 1 $RUNS); do 
                run_cmd "$i"
                $RUNCMDA
                msleep
                $RUNCMDB
                msleep
            done
            ;;
        "random interleaving")
            COUNTA=1
            COUNTB=1
            TOTALRUNS=$((RUNS * 2))
        
            while read -r choice; do 
                if [[ "$choice" == 0 ]]; then
                    run_cmd $COUNTA
                    $RUNCMDA
                    ((COUNTA++))
                    msleep
                else
                    run_cmd $COUNTB
                    $RUNCMDB
                    msleep
                    ((COUNTB++))
                fi 
            done < <(shuf -i 1-"$TOTALRUNS" | while read -r num; do echo $((num % 2)); done)
            ;;
    esac
}
