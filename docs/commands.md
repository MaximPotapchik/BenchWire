# Commands

These are reference for every `.env` variable BenchWire reads.
Two groups: 

1. **Shared** | Required. Apply regardless of mode.
2. **Per-mode** | Only used in the mentioned mode.

## Shared

| Variable | Meaning | Example |
|---|---|---|
| `RUNS` | Number of runs per side (compare mode runs this many for *each* side, not total). | `RUNS=5` |
| `METHODOLOGY` | `single`, `sequential`, `cycling`, or `random interleaving`. | `METHODOLOGY="random interleaving"` |
| `COOLDOWNTIMER` | Milliseconds slept between runs. | `COOLDOWNTIMER=1` |
| `LABEL` | Legend label, single mode only. | `LABEL="Raw"` |

## Single mode

| Variable | Meaning | Example |
|---|---|---|
| `EXEGESISBIN` | Path to the exegesis binary. | `EXEGESISBIN="$HOME/projects/llvm-project/build/bin/llvm-exegesis"` |
| `--flag=value` | Any exegesis flag, no suffix, passed straight through. | `--mcpu=native` |

## Compare modes | A side

| Variable | Meaning | Example |
|---|---|---|
| `EXEGESISBINA` | Path to the first binary being compared. | `EXEGESISBINA="$HOME/projects/llvm-project/build-raw/bin/llvm-exegesis"` |
| `LABELA` | Legend label for A. | `LABELA="Raw"` |
| `--flagA=value` | Any exegesis flag, suffixed `A`. | `--mcpuA=native` |

## Compare modes | B side

| Variable | Meaning | Example |
|---|---|---|
| `EXEGESISBINB` | Path to the second binary being compared. | `EXEGESISBINB="$HOME/projects/llvm-project/build-symbolic/bin/llvm-exegesis"` |
| `LABELB` | Legend label for B. | `LABELB="libpfm"` |
| `--flagB=value` | Any exegesis flag, suffixed `B`. | `--mcpuB=native` |

## Exegesis raw flag passthrough

Any line starting with `--` gets forwarded to the exegesis binary as-is.
This isn't a fixed list. `--mcpu`, `--mode`, `--opcode-name` are the ones with
support currently. To find the commands available with llvm-exegesis use `--help`.
BenchWire doesn't hardcode or validate flag names. 

One exception: `--benchmarks-files` is filtered out even if present.
BenchWire generates this one itself per run, so don't set it manually.
This will be customizable in the future.

## Full example `.env`

```
EXEGESISBINA="$HOME/projects/llvm-project/build-raw/bin/llvm-exegesis"
RUNS=5
--mcpuA=native
--modeA=latency
--opcode-nameA=ADD64rr
LABELA="Raw"

EXEGESISBINB="$HOME/projects/llvm-project/build-symbolic/bin/llvm-exegesis"
--mcpuB=native
--modeB=latency
--opcode-nameB=ADD64rr
LABELB="libpfm"

METHODOLOGY="random interleaving"
COOLDOWNTIMER=5
```
