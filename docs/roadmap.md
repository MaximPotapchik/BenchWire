# Roadmap
 
These are future plans for expansion of this project. In order of the current
priorities.

## Opcode sweep

This is currently in V1 and support is being expanded. It only supports X86 
right now. All memory setup required opcodes are skipped. V2 will seek to close
this gap. 

## pyperf integration

Adding support for pyperf, as the second supported benchmarker. It will be 
integrated with this tool to develop the CI system, using it on itself. This
will also open the door to json input analysis, since that is pyperf's output
format. The json parsing in analysis can be reused for other benchmarkers.

## uops.info alternative

A goal of this project is to provide an LLVM-based, instruction table that is
completely community driven. In order for this to be accurate, a stress testing
system is needed to prove the box being used to gather information is stable
enough to run the opcode sweep without causing jitter and tainted numbers. 
Accurate readings are critical. A database holding information on every 
architecture that has been swept and approved, is required. An optimized API
that can be called by godbolt.org or similar services, is a desired capability.

## Additional config customization

The ability to create custom chains of targets in `config.yaml`. The 
`config.yaml` options will be expanded significantly. This would include:

- Looping over a config option and inserting new variables at certain parts.

- Automated long-term scheduling support. 

- Plot customization. The ability to create custom themes.

- Report customization. Allowing the use of different output formats.

## InfluxDB export
 
Right now results only exist as timestamped markdown + png in
`results/plots/`. This is fine for the purposes of statistical
aggregation, but omits fine-grained per-run analysis. Adding time
series database support would enable long-term tracking with finer
controls.
 
Schema example:
 
```
measurement: exegesis_run
tags: opcode, mcpu, exegesis_mode, label, methodology, run_batch_id
fields: value, run_index
```
 
Example line:
 
```
exegesis_run,opcode=ADD64rr,mcpu=native,exegesis_mode=latency,label=Patch123,methodology=random_interleaving,run_batch_id=8f2a value=1.0088,run_index=42 1720260000000000000
```
 
Config lives in `config.yaml` or a similar mechanism, gated behind an explicit
flag so the tool has zero InfluxDB dependency by default.
 
```
INFLUX_ENABLED=false
INFLUX_URL=http://localhost:8086
INFLUX_TOKEN=
INFLUX_ORG=
INFLUX_BUCKET=exegesis
```
 
Enables quick Grafana integration.
 
Tagging via `build_sha`, pulled from LLVM's git or any future target
repo. This can allow a future regression-tracking system to be built
on top.
 
Can be dockerized within its own separate module as well. Support for
more DBs would be useful. This does not have to be focused on InfluxDB but it 
is a solid choice.

## More under consideration
 
- `--ci` | The goal is to have this tool be used as CI for itself. Benchmarking
the benchmark automation and analysis.

- `Scheduler` | Automated periodic or event-triggered runs (nightly, or
on every LLVM commit), feeding directly into the InfluxDB + build_sha
pipeline above.

- `JSON export` | Output alongside the markdown summary, for CI or
dashboard consumption without requiring a live InfluxDB instance.

- `Enhanced analytics` | Expanding the python analysis with further
statistical categories, calculations, and more.

- `TUI` | Using https://github.com/charmbracelet/bubbletea to build a more
interactive TUI. 
