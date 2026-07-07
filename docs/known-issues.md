# Known issues

Current gaps in the implementation.

## Custom command mode (option 2) is sequential only

There's no methodology selection in option 2, it runs the A batch then
the B batch back to back with no gap between them, and no cycling or
interleaving choice. Those require `.env` mode.

## Custom command mode has no label support

`LABEL` / `LABELA` / `LABELB` only come from `.env`, option 2 has no
prompt for them, so plots from option 2 fall back to generic `Run` /
`A` / `B` labels instead of anything descriptive. Cosmetic, but editing
`.env` is the only way, even if you're using flags at the prompt.

## `analyze.py` only understands latency mode output

`ExegesisMode` is hardcoded to `"Latency"` and `LoadRun` grabs the first
entry in `measurements`. `uops`, or `inverse_throughput` mode output isn't
parsed currently.

## No validation on `.env` values

A bad `.env` fails with a raw bash or exegesis error instead of a message
pointing at the actual line. Could make it ask for what it doesn't have
as well.

## No resume support

If a run batch is interrupted partway through, the partial yaml files
stay in `results/yaml/`. The next `bench.sh` invocation deletes them and
starts over from zero, there's no resume.

