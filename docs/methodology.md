# Comparison methodology

Compare mode isolates implementation differences from noise introduced
by run-to-run conditions.

Sequential comparison, all of A then all of B, confounds the two.
Anything that drifts over the run window (thermal state, frequency
scaling, background load) concentrates entirely in whichever side ran
second.

## The three orderings

`single` | Runs a single binary, no comparison involved.

`sequential` | All A runs, then all B runs. Simplest for A/B, most exposed
to drift. Fine for a sanity check, not sufficient to support a defended
claim.

`cycling` | Strict A, B, A, B alternation. Spreads drift evenly across
both sides instead of concentrating it in one.

`random interleaving` | Run order shuffled, counts held exactly even (N
per side, by construction, not by approximation of the shuffle). 
Matches the approach Google Benchmark uses via 
'--benchmark_enable_random_interleaving=true' to reduce systematic
bias across a suite. The trade-off: this lowers drift-correlated bias but
introduces state sensitivity between consecutive runs.

## Reading output

Percentile spread (P75 to P99.9) carries as much information as the
mean. Near-identical means with divergent tails is common.
