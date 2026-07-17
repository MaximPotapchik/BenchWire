#!/bin/bash
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"

HASPYTHON3=0
HASPIP=0
HASGO=0

if command -v python3 >/dev/null 2>&1; then
    HASPYTHON3=1
else
    HASPYTHON3=0
    echo "[BenchWire] Please install python3."
fi

if command -v pip &> /dev/null; then
    HASPIP=1
else
    HASPIP=0
    echo "[BenchWire] Please install pip."
fi

# Python dependancies. 
if [[ $HASPYTHON3 -eq 1 && $HASPIP -eq 1 ]]; then
    if [[ -f $SCRIPT_DIR/requirements.txt ]]; then
        if python3 -c "import numpy, matplotlib, yaml" 2>/dev/null; then
            echo "[BenchWire] Python deps already satisfied, skipping install."
        else
            pip install -r requirements.txt
        fi
    else
        echo "[BenchWire] requirements.txt not found. Please git clone the repository again."
    fi
fi

# .env check.
if [[ ! -f "$SCRIPT_DIR/.env" || ! -s "$SCRIPT_DIR/.env" ]]; then
    if [[ -f "$SCRIPT_DIR/.env.example" ]]; then
        cp "$SCRIPT_DIR/.env.example" "$SCRIPT_DIR/.env"
        echo "[BenchWire] .env file empty, .env.example used. Please fill it with your own variables."
    else
        echo "[BenchWire] .env.example not found."
        echo "[BenchWire] Please git clone the repository again or create .env manually."
    fi
fi

# Go installations.
if command -v go >/dev/null 2>&1; then
    go build -o "$SCRIPT_DIR/benchwire" "$SCRIPT_DIR/cmd/benchwire"
    HASGO=1
    
else
    HASGO=0
    os=$(uname -s | tr '[:upper:]' '[:lower:]')
    arch=$(uname -m)
    [[ "$arch" == "x86_64" ]] && arch="amd64"

    curl -L -o "$SCRIPT_DIR/benchwire.tar.gz" \
        "https://github.com/MaximPotapchik/BenchWire/releases/latest/download/BenchWire_${os}_${arch}.tar.gz"
    tar -xzf "$SCRIPT_DIR/benchwire.tar.gz" -C "$SCRIPT_DIR"
    rm -f "$SCRIPT_DIR/benchwire.tar.gz"
fi

echo "[BenchWire] After populating your .env, use ./benchwire to run the harness."
echo "[BenchWire] To see available commands, check docs/commands.md."
echo "[BenchWire] If you would like to use the legacy bash command, use ./bench.sh"

chmod +x bench.sh
