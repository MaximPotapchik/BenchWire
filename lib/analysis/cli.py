import argparse
from enum import Enum

class METHODOLOGY(str, Enum):
    SINGLE = "single"
    SEQUENTIAL = "sequential"
    CYCLING = "cycling"
    RANDOMINTERLEAVING = "random interleaving"

# TODO: output formats line, once support is added.
# cli variable grabbing.
cli_args = argparse.ArgumentParser(description="Python CLI args")
cli_args.add_argument("methodology", type=str, choices=list(METHODOLOGY))
cli_args.add_argument("runs", type=int)
cli_args.add_argument("cooldowntimer", type=int)
cli_args.add_argument("labels", type=str)
cli_args.add_argument("collectedstats", type=str)
# cli_args.add_argument("outputformats", type=str) # This is for when there are different
# output file formats.

args = cli_args.parse_args()

# TODO: Change the collectedStats split, to catch more safely.
cliOpts = {
    "methodology" : args.methodology,
    "runs" : args.runs,
    "cooldownTimer" : args.cooldowntimer,
    "labels" : args.labels,
    "collectedStats" : args.collectedstats.split(",")
}
