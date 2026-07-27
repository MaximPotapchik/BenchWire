# shellcheck shell=bash
#
declare -a FORMATTEDARRAY=()

# TODO: - The gate validation here should be strong because everything
# breaks otherwise.
# - Allow cli mode here.
exegesis_formatter() {
    local METHODOLOGY=$1
    local -n RAWARRAY=$2

    if [[ -v EXEGESISBIN ]]; then
        FORMATTEDARRAY=("$EXEGESISBIN")
    elif [[ -v EXEGESISBINA ]]; then
        FORMATTEDARRAY=("$EXEGESISBINA|$EXEGESISBINB")
    else 
        echo "No EXEGESISBIN specified. Add your Exegesis bin locations into .env"
        exit 1
    fi

    for i in "${RAWARRAY[@]}"; do
        local LHS="${i%%=*}="
        local RHS="${i##*=}"
        local CLICHECK="${LHS:0:2}"

        # Regex for if --benchmarks-file exists and if it ends with A/B.
        if [[ "$LHS" =~ ^-{0,2}benchmarks-file[AB]?=?$ ]]; then
            continue
        fi 

        if [[ "$CLICHECK" == "--" ]]; then
            FORMATTEDARRAY+=("$i")
        fi
    done
}
