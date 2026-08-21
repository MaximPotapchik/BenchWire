# Config

Reference for every `config.yaml` key BenchWire reads. Each may have a sub-list of keys.
This is the current list: 

1. [Context](#context) | Metadata about this config batch, currently name and mode.
2. [Default](#default) | BenchWire settings that get applied by default, if specMatrix does not invoke them.
3. [Global Flags](#global-flags) | Flags applied to every target, in every `specMatrix`.
4. [Presets](#presets) | Reusable, named flag bundles, with one layer of inheritance.
5. [Spec Matrix](#spec-matrix) | The actual benchmark definitions, one entry per comparison batch.
6. [Matrix Sequence](#matrix-sequence) | Declared ordering of `specMatrix` entries.
7. [Analysis](#analysis)| Where and how BenchWire's output gets written.

## Context

| Key | Meaning | Example |
|---|---|---|
| `name` | Identifies this config batch. | `name: "Sweep"` |
| `mode` | Read by the resolver. Defaults to `"default"` if unset. | `mode: "default"` |

## Default

| Key | Meaning | Example |
|---|---|---|
| `runs` | Number of runs per target (compare mode runs this many for *each* target, not total). | `runs: 20` |
| `methodology` | `single`, `sequential`, `cycling`, or `random interleaving`. | `methodology: "random interleaving"` |
| `cooldownTimer` | The amount of time in between runs. Can be randomized with an interval. See [CooldownTimer](#cooldowntimer).  | `cooldownTimer: "5ms"` |


## Global Flags

For more information on how flags are handled, see (Flag assembly order)[#flag-assembly-order].

| Key | Meaning | Example |
|---|---|---|
| `globalFlags` | A list of flags applied to every target in every `specMatrix`. | `globalFlags: [--mcpu=native]` |

## Presets

Named, reusable flag bundles. 

| Key | Meaning | Example |
|---|---|---|
| `name` | How targets and other presets reference this one. | `name: Vex` |
| `inherit` | Names of top-level presets to pull flags from first. Only use inside child presets. | `inherit: [Vex]` |
| `flags` | The flags this preset contributes. | `flags: [--opcode-name=VEXTRACTF128rri]` |

**Inheritance is only for child presets of specMatrix only**
- Top-level presets cannot use `inherit` at all.
- A local preset can only inherit from a top-level preset, never from another local preset.

## Spec Matrix

Each entry is one comparison batch, its own methodology, its own targets,
its own run sequence.

| Key | Meaning | Example |
|---|---|---|
| `name` | Identifies this batch, used in filenames and reports. | `name: "VEXTRACTF128rri FPU pipe quick check"` |
| `benchmarker` | Which benchmarker this batch runs. Currently only `"llvm-exegesis"`. | `benchmarker: "llvm-exegesis"` |
| `methodology` | Overrides `default.methodology` for this batch. | `methodology: "sequential"` |
| `runs` | Overrides `default.runs`. | `runs: 20` |
| `cooldownTimer` | Overrides `default` `cooldownTimer`. See [CooldownTimer](#cooldowntimer). | `cooldownTimer: "10ms"` |
| `localFlags` | Extra flags for this batch, **added on top of** `globalFlags`. | `localFlags: [--mode=uops]` |
| `presets` | This batch's own presets. See [Presets](#presets). | see above |
| `targets` | The binaries being compared. See [Targets](#targets). | see below |
| `sequence` | Execution order. See [Sequence](#sequence). | see below |

### Targets

| Key | Meaning | Example |
|---|---|---|
| `label` | Legend label used in plots and reports. | `label: "Raw"` |
| `binPath` | Path to the exegesis binary. Supports `$VAR`-style env expansion. | `binPath: "$HOME/llvm-project/build-raw/bin/llvm-exegesis"` |
| `preset` | Preset name(s) this target pulls flags from. **Always a list, brackets required even for one.** | `preset: [VexLowMin]` |
| `flags` | A `target`'s own specified flags. Applied after all other flags. | `flags: [--mcpu=native]` |

### Sequence

Controls the order and groupings of each target.

| Form | Meaning | Example |
|---|---|---|
| Single target label | Solo run of one target. Forces `methodology: single` for this step, regardless of the batch's own methodology. | `- Build1` |
| Two target list | A/B comparison, using the batch's methodology. Do not use `single` here. | `- [build1, build2]` |

### CooldownTimer

| Key | Meaning | Example |
|---|---|---|
| `randomizeWithin` | Upper bound of the random range. | `randomizeWithin: 5` |
| `precision` | Unit for `randomizeWithin`. `us`, `ms` or `s`. | `precision: "ms"` |

This allows simulation of instability by randomzing sleep within the minimum allowed
by the operating system, and your number.

Example: `cooldownTimer: {randomizeWithin: 5, precision: "ms"}`

## Matrix Sequence

**Not currently supported.** Will be configurable in the future.

| Key | Meaning | Example |
|---|---|---|
| `matrixSequence` | A list of `specMatrix` names, declaring intended order. | see below |

## Analysis

Controls where reports get written.

| Key | Meaning | Example |
|---|---|---|
| `dirName` | Optional. If set, reports go in `results/{dirName}/` instead of the default location. | `dirName: "reports"` |
| `storeIn` | `"directory"` or `"single file"`. | `storeIn: "directory"` |
| `per` | How you'd like to target storeIn. `"suite"`, `"specMatrix"`, or `"run"`. **Ignored entirely if `storeIn` is `"single file"`.** | `per: "specMatrix"` |
| `plots` | Options for plots. See [Plots](#plots).  Expansion of customization options will be built. | see below |

`per` decides grouping when `storeIn: "directory"`:
- `"suite"` | Every `specMatrix` in the whole run combines into one report, in one `batch_{timestamp}/` folder.
- `"specMatrix"` | Each `specMatrix` gets its own combined report, in its own subfolder.
- `"run"` | Every individual sequence step gets its own separate file, no combining, no subfolders.

If `storeIn: "single file"`, everything above is skipped, always exactly
one report, containing every `specMatrix` and `every step, regardless of
what `per` says.

Neither `storeIn` nor `per` currently has a default, both are required
keys. You will get an error otherwise.

### Plots

| Key | Meaning | Example |
|---|---|---|
| `disable` | disables plotting when set to true | `disable: false` |

## Flag assembly order

For a given target, flags are assembled in this exact order:

1. `globalFlags`.
2. Each `specMatrix`'s own `localFlags`.
3. Each `preset` the target lists. In the specified order.
4. The `target`'s own `flags`.

## Exegesis flag passthrough

Every flag specified for a target gets forwarded to the exegesis binary as-is.
Only `--mode=analysis` is not supported. Further mode support is being pursued.
To find the commands available with llvm-exegesis use `--help`.
BenchWire doesn't hardcode or validate flag names except for: `--benchmarks-file`. 
BenchWire generates this one itself per run, so don't set it manually. This will
be customizable in the future.

## Full example `config.yaml`

```yaml
default:
  methodology: "single"
  runs: 50
  cooldownTimer: "5ms"

globalFlags:
  - --mcpu=native

presets:
  - name: Vex
    flags: [--opcode-name=VEXTRACTF128rri]

specMatrix:
  - name: "VEXTRACTF128rri FPU pipe quick check"
    benchmarker: "llvm-exegesis"
    methodology: "sequential"
    runs: 20
    localFlags:
      - --mode=uops

    presets:
      - name: VexLowMin
        inherit: [Vex]
        flags: [--min-instructions=5000]

    targets:
      - label: "build1"
        binPath: "Add your binary path here"
        preset: [VexLowMin]

      - label: "build2"
        binPath: "Add your binary path here"
        preset: [VexLowMin]
        flags:
          - --repetition-mode=duplicate

    sequence:
        - build1
        - build2
        - [build1, build2]

matrixSequence:
  - "VEXTRACTF128rri FPU pipe quick check"

analysis:
  output:
    storeIn: "directory"
    per: "specMatrix"
```
