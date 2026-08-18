# CLI

This is a reference to all CLI that BenchWire uses. Broken down by sections.

1. [Config](#config) | Options for running or building a config file.
2. [Cleanup](#cleanup) | Removing generated outputs.
3. [Other](#other) | Misc utilities or not implemented yet.

## Config

Every Config related CLI option.

### --config

Runs a config file from any location provided. If you generate a config with 
`--buildConfig`, you can target it with this. This will run the config
directly, instead of using the base config.yaml.

**Example:** --config=/configs/custom/customMatrix.yaml

### --buildConfig

Builds a config from a preset. Takes a string. The following are the different
ways to invoke this.

1. Preset names existing inside the `configs` directory, can be passed directly.
Will ask for a new config name and the required binary path(s). 

**Example:** `--buildConfig=50RunAdd64.yaml`

2. Preset path if you place it outside of `configs`. It must be inside the
BenchWire directory.

**Example:** `--buildConfig=/customDir/customConfig.yaml` 

3. Sweep. Does an opcode sweep for `LLVM-exegesis`. Currently only supports 
X86.

**Example:** `--buildConfig=sweep`

It will ask for 3 prerequisites.

First: The full filepath location of your `LLVM-exegesis` binary. 
For example `/root/llvm-project/build/bin/llvm-exegesis`.

Second: The full filepath location of your X86 tablegenerated instruction info.
For example `/root/llvm-project/build/lib/Target/X86/X86GenInstrInfo.inc`.

Third: The `LLVM-exegesis` mode you wish to use. For example: latency, uops,
inverse_throughput, or full for all three.

## Cleanup

### --clean

Removes everything under `results/`. Safe, only ever touches paths inside
the results directory. Run as a bare argument, no value.

**Use:** ./benchwire clean

## Other

### --ci

Not implemented yet.
