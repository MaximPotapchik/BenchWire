# Commands

These are reference for every `config.yaml` variable BenchWire reads.
Two groups: 

1. **Required** | Apply regardless of mode.
2. **Targets** | A list with one entry per binary being benchmarked.

## Required

| Key | Meaning | Example |
|---|---|---|
| `runs` | Number of runs per target (compare mode runs this many for *each* target, not total). | `runs: 20` |
| `methodology` | `single`, `sequential`, `cycling`, or `random interleaving`. | `methodology: "random interleaving"` |
| `cooldownTimer` | Milliseconds slept between runs. | `cooldownTimer: 5` |
| `targets` | A list of targets. One entry for single mode, two for any compare mode. | see below |

## Targets

| Key | Meaning | Example |
|---|---|---|
| `label` | Legend label used in plots and reports. | `label: "Raw"` |
| `binPath` | Path to the exegesis binary. Supports `$VAR`-style env expansion. | `binPath: "$HOME/projects/llvm-project/build-raw/bin/llvm-exegesis"` |
| `flags` | A list of raw exegesis flags, passed straight through. | `flags: --mcpu=native` |

Single `methodology` uses one `targets` entry. Any compare mode (`sequential`, `cycling`, `random interleaving`) uses two.

## Exegesis flag passthrough

Every string in a target's `flags` list gets forwarded to the exegesis binary as-is.
Currently, only --mode=latency is supported. Further mode support is being persued.
To find the commands available with llvm-exegesis use `--help`.
BenchWire doesn't hardcode or validate flag names except for: `--benchmarks-file`. 
BenchWire generates this one itself per run, so don't set it manually. This will
be customizable in the future.

## Full example `config.yaml`

```yaml
methodology: "random interleaving"
runs: 100
cooldownTimer: 5

targets:
  - label: "Build A"
    binPath: "$HOME/llvm-project/buildA/bin/llvm-exegesis"
    flags:
      - --mcpu=native
      - --mode=latency
      - --opcode-name=ADD64rr

  - label: "Build B"
    binPath: "$HOME/llvm-project/buildB/bin/llvm-exegesis"
    flags:
      - --mcpu=native
      - --mode=latency
      - --opcode-name=ADD64rr
```
