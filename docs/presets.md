# Presets

A list of presets that come with BenchWire. Currently only supports 
`LLVM-exegesis`. To use one of these, simply type: 
./BenchWire --buildConfig=preset_name.yaml

Every preset is modifiable. You can create your own presets as well, inside
`configs/presets/`. 

If you find that your own preset is not found, but you think it should be,
please don't hesitate to submit a pull request to add it.

## LLVM-exegesis
| Name | Description |
| --- | --- |
| 50RunAdd64.yaml | Runs ADD64rr 50 times on a single binary. |
| BinComparisonAdd64.yaml | Compares 2 exegesis Binaries in the same way 50RunAdd64 does. Useful for regression detection. |
