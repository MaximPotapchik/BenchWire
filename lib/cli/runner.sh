#!/bin/bash
set -e

argsarr=("location" "--outputlocationcli=|outputlocation" "--test1=|test1" "--test2=|test2")
comparearr=("locationA|locationB" "--outputlocationcliA=|outputlocationA" "--outputlocationcliB=|outputlocationB" "--testA=|testA" "--testB=|testB")

# TODO: This can be further modularized. Should probably get its own folder.
runner() {
    # sleep/cooldown function
   . "$(dirname "${BASH_SOURCE[0]}")/msleep.sh"

    local -n METHODOLOGY=$1
    local RUNS=$2
    local -n ARGS=$3

    local COMMAND=""
    local COMMANDA=""
    local COMMANDB=""

    if [[ "${ARGS[0]}" == *\|* ]]; then
        COMMANDA+="${ARGS[0]%%|*}"
        COMMANDB+="${ARGS[0]##*|}"
    else
        COMMAND="${ARGS[0]}"
    fi

    if [[ "${ARGS[1]#|}" == *A ]]; then
        local curI=3
        local total=${#ARGS[@]}

        while [[ $curI -lt $total ]]; do 
            local opt="${ARGS[$curI]}"
            local cli="${opt%%|*}"
            local val="${opt##*|}"
            # Check 2nd last letter for A. cli ends with = usually.
            local check="${cli: -2: 1}"

            if [[ "$check" == "A" ]]; then
                COMMANDA+=" $cli$val"
            else
                COMMANDB+=" $cli$val"
            fi 
            ((curI++))
        done
    else
        for opt in "${ARGS[@]:2}"; do
            local cli="${opt%%|*}"
            local val="${opt##*|}"
            COMMAND+=" $cli$val"
        done
    fi

    run_cmd() {
        local RUN=$1
        if [[ "${ARGS[1]#|}" == *A ]]; then
            local outputlocationcliA="${ARGS[1]%%|*}"
            local outputlocationA="${ARGS[1]##*|}"
            RUNCMDA="$COMMANDA $outputlocationcliA$outputlocationA/Arun_$RUN.yaml"

            local outputlocationcliB="${ARGS[2]%%|*}"
            local outputlocationB="${ARGS[2]##*|}"
            RUNCMDB="$COMMANDB $outputlocationcliB$outputlocationB/Brun_$RUN.yaml"
        else
            local outputlocationcli="${ARGS[1]%%|*}"
            local outputlocation="${ARGS[1]##*|}"
            RUNCMD="$COMMAND $outputlocationcli$outputlocation/run_$RUN.yaml"
        fi
    }

    case "$METHODOLOGY" in
        "single")
            for i in $(seq 1 $RUNS); do 
                run_cmd "$i"
                echo "$RUNCMD"
                msleep
            done
            ;;
        "sequential")
            for i in $(seq 1 $RUNS); do 
                run_cmd "$i"
                echo "$RUNCMDA"
                msleep
            done

            for i in $(seq 1 $RUNS); do 
                run_cmd "$i"
                echo "$RUNCMDB"
                msleep
            done
            ;;
        "cycling")
            for i in $(seq 1 $RUNS); do 
                run_cmd "$i"
                echo "$RUNCMDA"
                msleep
                echo "$RUNCMDB"
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
                    echo "$RUNCMDA"
                    ((COUNTA++))
                else
                    run_cmd $COUNTB
                    echo "$RUNCMDB"
                    ((COUNTB++))
                fi 
            done < <(shuf -i 1-"$TOTALRUNS" | while read -r num; do echo $((num % 2)); done)
            ;;
    esac
}

SINGLE="single"
CYCLING="cycling"
SEQUENTIAL="sequential"
INTERLEAVING="random interleaving"

runner SINGLE 5 argsarr
runner SEQUENTIAL 3 comparearr
runner CYCLING 5 comparearr
runner INTERLEAVING 3 comparearr
