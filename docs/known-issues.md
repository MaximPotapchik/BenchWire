# Known issues

Current gaps in the implementation.

## No resume support

If a batch is interrupted partway through, the partial yaml files stay in
`results/yaml/`. The next invocation deletes or overwrites them and starts over
from zero, there's no resume.

## Errors being uncaught

There are gaps in reporting of where certain config variables or benchmarker
outputs error. 

## Unoptimized aggregator loop.

Currently the analysis duration is slow, bottlenecked by an inefficient 
algorithm. This is currently being addressed.
