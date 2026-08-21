# BenchWire 

An automated design of experiments system, currently used for [`LLVM-exegesis`](https://llvm.org/docs/CommandGuide/llvm-exegesis.html).
It runs single or A/B comparison benchmarks automatically, computes statistics,
and plots the result. Rich config customization. The goal is to have a modular,
extensible framework capable of targetting any benchmarker, with continuous
integration in mind. Along with making an LLVM-based alternative for uops.info
tables.

![example image of plot](results/examples/exampleplot.png)

<details>
<summary>Example statistical output</summary>

**Methodology:** random interleaving

**Runs:** 20

**CPU:** znver2

**Triple:** x86_64-unknown-linux-gnu

**Min instructions:** 10000

### Summary Statistics

| Metric | Build A | Build B | Diff |
|---|---|---|---|
| Mean | 1.0107 | 1.0093 | 0.13% (Build A) |
| Median | 1.0091 | 1.0093 | 0.02% (Build B) |
| Standard Deviation | 0.0108 | 0.0037 | 98.09% (Build A) |
| Variance | 1.1742e-04 | 1.3725e-05 | 158.14% (Build A) |
| Coefficient of Variation | 1.07% | 0.37% | 97.99% (Build A) |
| Range | 0.0779 | 0.0277 | 95.08% (Build A) |
| IQR | 0.0008 | 0.0005 | 43.14% (Build A) |
| Min | 1.0055 | 1.0053 | 0.02% (Build A) |
| Max | 1.0834 | 1.0330 | 4.76% (Build A) |

### Percentile Statistics

| Percentile | Build A | Build B | Diff |
|---|---|---|---|
| P50 | 1.0091 | 1.0093 | 0.02% (Build B) |
| P75 | 1.0094 | 1.0094 | 0.00% (Build A) |
| P90 | 1.0098 | 1.0098 | 0.00% (Build A) |
| P99 | 1.0541 | 1.0236 | 2.93% (Build A) |
| P99.9 | 1.0805 | 1.0321 | 4.58% (Build A) |
| P99.99 | 1.0831 | 1.0329 | 4.74% (Build A) |

</details>

## Overview

1. Runs a benchmarker per `specMatrix`, solo or head-to-head, based
on that matrix's `sequence`. A single `config.yaml` can define multiple
independent comparisons in one batch.

2. Supports four different run methodologies (single, sequential, cycling,
random interleaving) specifically to control for time-based bias like thermal
drift and frequency scaling skewing numbers, see [`docs/methodology.md`](docs/methodology.md)
for why this matters.

3. Parses the benchmarker's outputs and automatically performs statistical 
analysis on it.

4. Produces a Catppuccin-themed plot and a markdown stats summary (mean, 
median, stddev, CoV, percentiles up to P99.99 + more) for every run.

## Features

1. **YAML-driven config system.** Benchmark runs are declarative with default
settings, flags, presets, and more. There are intuitive inheritance structures
built-in.

2. **Customization.** Every config can be edited with your own features.
Various different combinations are possible, all modifiable by the user. For 
more information, see [`docs/config.md`](docs/config.md).

3. **Config generation.** Simply use `./BenchWire --buildConfig=` followed by
the location of your desired preset. It will be filled with your specified 
paths, and prompt you whether to run it immediately after, or not. Comes with a
list of presets. See [`docs/presets.md`](docs/presets.md).

4. **Opcode sweep generation.** Creating a config with a list of every opcode 
`LLVM-exegesis` can run on your architecture is done using:
`./BenchWire --buildConfig=sweep`. A text list of every opcode per bucket, is
generated as well. It can run 3 different `LLVM-exegesis` modes back-to-back.
More info in [`docs/cli.md`](docs/cli.md).

5. **Full statistical aggregation.** Mean, median, standard deviation,
variance, coefficient of variation, range, interquartile range, and
percentiles. Computed across every run, for any measurement a preset
tracks. Further statistical coverage will be implemented.

6. **Markdown reporting.** Single-file, per-matrix, or per-run outputs. These 
are automatically diffed when using comparison methodologies. Generated plots
are embedded as well.

7. **Plot generation.** Plots statistics using an automated plot generation API
powered by Matplotlib. Plots are in a Catppuccin-Mocha theme.

8. **Multi-benchmarker support.** Results can be read from yaml, json, or
built entirely in memory. Wide reaching benchmarker support is the goal.

9. **Noise Control.** Artificial jitter can be added via different ways.
Both via randomization of cooldowns, and methodology.

## Requirements

- `LLVM-exegesis` and its dependencies. 
- Bash
- Python 3

## Quick Start

```bash
git clone https://github.com/MaximPotapchik/BenchWire
cd BenchWire
chmod +x setup.sh
./setup.sh
```

For the full list of `config.yaml` settings, see [`docs/commands.md`](docs/commands.md).

This will run the setup script, creating the `config.yaml` and populating it
with `config.example.yaml`.

Edit `config.yaml` with your binary path(s), along with your desired 
LLVM-exegesis flags, then run:

```bash
./benchwire
```

Benchmarker results appear in `results/yaml/` (raw exegesis output per run) and 
the analysis lands in the chosen directory (a plot + a markdown stats summary,
timestamped). 

You may also use `./benchwire --buildConfig=` with a preset to automatically
create a config and use it.

## Comparison methodology

Set via `methodology:` in `config.yaml`:

- `single` | Runs a single binary.
- `sequential` (default) | All A runs, then all B runs.
- `cycling` | A, B, A, B, etc.
- `random interleaving` | Shuffled order, same N runs each.

Full reasoning in [`docs/methodology.md`](docs/methodology.md). Short
version: naive back-to-back comparison lets anything that drifts over
time (thermal state, frequency scaling, whatever else is on the box) get
absorbed entirely into whichever side ran second.

## Docs

- [`docs/cli.md`](docs/cli.md) | Every CLI option BenchWire uses.
- [`docs/config.md`](docs/config.md) | Every `config.yaml` key, what it does.
- [`docs/known-issues.md`](docs/known-issues.md) | Current gaps and rough edges.
- [`docs/methodology.md`](docs/methodology.md) | Why run ordering matters.
- [`docs/presets.md`](docs/presets.md) | Built-in presets, ready to run.
- [`docs/roadmap.md`](docs/roadmap.md) | What's planned but not built yet.

## In progress

Version 0.0.6 will include wider support in the opcode sweep.

## Contributing

If you would like to contribute, look at `docs/known-issues.md`. If you have
experience with Golang and/or Python, your feedback/help would be appreciated. 
Open an issue before PRing.

Compiler/LLVM background is helpful but not necessary.

## License

MIT, see [`LICENSE`](LICENSE).
