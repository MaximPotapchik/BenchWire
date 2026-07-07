# Roadmap
 
These are future plans for expansion of this project.
 
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
 
Config lives in `.env` or a similar mechanism, gated behind an explicit
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
more DBs would be useful.
 
## Additional command utilities
 
These may either be their own utility or a feature flag for `bench.sh`.
 
- `--mode scan` | Take a file or comma-separated list of opcodes and
run single or compare mode across all of them in one invocation, one
combined report at the end. High chance of needing a separate utility.
Ties in well with a potential scheduler feature.

- `--validate` | Preflight check. Validates whether the run should
begin or not, in case the user has a large N of runs. Would prevent
wasted compute.

- `--resume` | Would be paired with a system to track N of runs, in
case the runs were interrupted for whatever reason.

## Run progress and ETA
 
`bench.sh` currently prints `Run i/N complete` and nothing else. For a
1000+ run batch this becomes an unhelpful wall of identical lines. Can
track wall-clock duration per run, keep a rolling average and print
something like:
 
```
[===========                    ] 34% (340/1000) | avg 42ms/run | ETA 27s
```
 
## Also under consideration
 
- `Modularization` | Split/Refactor parts as more module-first
features. Keeps `bench.sh` from growing into a monolith as features
land.

- `Scheduler` | Automated periodic or event-triggered runs (nightly, or
on every LLVM commit), feeding directly into the InfluxDB + build_sha
pipeline above.

- `JSON export` | Output alongside the markdown summary, for CI or
dashboard consumption without requiring a live InfluxDB instance.

- `Parsing & other optimizations` | It may be viable to refactor
certain sections of the Python analysis with C++, after profiling.

- `LLVM exegesis analysis sweep` | This would likely orchestrate
llvm-exegesis's own `-mode=analysis` and
`-analysis-inconsistencies-output-file` across many opcodes, rather
than reimplementing scheduling-model comparison from scratch.

