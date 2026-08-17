package builder

import (
	"benchwire/internal/config"
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"slices"
	"sort"
	"strconv"
	"strings"

	"golang.org/x/sys/cpu"
	"go.yaml.in/yaml/v3"
)

// To understand how opcodes are encoded and what we look for, use: llvm/MC/MCInstrDesc.h.
// Currently only supports X86.
// TODO: - Add further architecture support.
// - Build bucket two, which is memory annotation required opcodes.
// - Move to a seperate file.

// X86BaseInfo encoding mask and shift.
const EncodingShift = 29
const EncodingMask = 0x3 << EncodingShift

// EVEX requires AVX-512 support. We look for it to make sure it is, on X86.
const (
	EncodingLegacy = 0 << EncodingShift
	EncodingVEX = 1 << EncodingShift
	EncodingXOP = 2 << EncodingShift
	EncodingEVEX = 3 << EncodingShift
)

// One row of the opcode table.
var row = regexp.MustCompile(`^\s*\{\s*\d+`)
// Name after the //.
var name = regexp.MustCompile(`//\s*(\S+)\s*$`)
// Grabs the hex number at the of (arch)Descs 
var TSFlagsHex = regexp.MustCompile(`0x([0-9A-Fa-f]+)ULL\s*\},`)
// Enum with info on the tablegen origin of the opcode.
var enumLine = regexp.MustCompile(`^\s*(\w+)\s*=\s*\d+,?\s*//\s*([^:\s]+):`)

type OpcodeInfo struct {
	Name string
	Skipped bool
	Encoding uint64
}

// Finds every opcode whose source .td is skipped.
func ReadSourced(lines []string, skipSources []string) map[string]bool {
	skipNames := map[string]bool{}
	for _, line := range lines {
		match := enumLine.FindStringSubmatch(line)
		if match == nil {
			continue
		}

		name := match[1]
		sourceFile := match[2]

		if slices.Contains(skipSources, sourceFile) {
			skipNames[name] = true
		}
	}

	return skipNames
}

func ReadLines(path string) ([]string, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	return lines, scanner.Err()
}

func findTableStart(lines []string) (int, error) {
	for i, line := range lines {
		if strings.Contains(line, "X86Descs = {") {
			return i, nil
		}
	}

	return -1, fmt.Errorf("[BenchWire] Cannot find the X86Descs table")
}

// Finds the last line of the opcode table.
func findTableEnd(lines []string, start int) (int, error) {
	var rowIdxs []int
	for i := start; i < len(lines); i++ {
		if row.MatchString(lines[i]) {
			rowIdxs = append(rowIdxs, i)
		}
	}
	if len(rowIdxs) == 0 {
		return -1, fmt.Errorf("[BenchWire] Cannot find any row lines")
	}

	end := rowIdxs[0]
	for _, idx := range rowIdxs[1:] {
		if idx != end+1 {
			break
		}
		end = idx
	}
	return end, nil
}

// Reads one table row at a time. Returns ok=false if it's not a real row. 
func parseRow(line string, sourcedNames map[string]bool, skipFlags []string) (info OpcodeInfo, ok bool) {
	if !row.MatchString(line) {
		return OpcodeInfo{}, false
	}

	match := name.FindStringSubmatch(line)
	if match == nil {
		return OpcodeInfo{}, false
	}
	name := match[1]

	skipped := false
	for _, flag := range skipFlags {
		if strings.Contains(line, flag) {
			skipped = true
			break
		}
	}
	if !skipped && sourcedNames[name] {
		skipped = true
	}

	var encoding uint64
	if tm := TSFlagsHex.FindStringSubmatch(line); tm != nil {
		if val, err := strconv.ParseUint(tm[1], 16, 64); err == nil {
			encoding = val & EncodingMask
		}
	}

	return OpcodeInfo{ Name: name, Skipped: skipped, Encoding: encoding }, true
}

// Reads the opcode flags and appends the ones we want.
func ExtractBucketOne(incPath string) ([]OpcodeInfo, error) {
	lines, err := ReadLines(incPath)
	if err != nil {
		return nil, err
	}

	start, err := findTableStart(lines)
	if err != nil {
		return nil, err
	}

	end, err := findTableEnd(lines, start)
	if err != nil {
		return nil, err
	}

	skipSources := []string{"Target.td", "X86InstrCompiler.td"}
	sourcedNames := ReadSourced(lines, skipSources)
	// These are flags we skip, found in the generated InstrInfo.inc.
	skipFlags := []string{
		"MayLoad", "MayStore", "UnmodeledSideEffects", "MCID::Pseudo",
		"MCID::Call", "MCID::Return", "MCID::Branch", "MCID::IndirectBranch",
		"MCID::Barrier", "MCID::Terminator", "MCID::UsesCustomInserter",
	}

	var results []OpcodeInfo
	for i := start; i <= end; i++ {
		if info, ok := parseRow(lines[i], sourcedNames, skipFlags); ok {
			results = append(results, info)
		}
	}
	return results, nil
}

// Checks if the X86 CPU is AVX-512 compatible.
func FilterBucketOne(info []OpcodeInfo) []string {
	var opcodes []string

	for _, r := range info {
		if r.Skipped {
			continue
		}

		if r.Encoding == EncodingEVEX && !cpu.X86.HasAVX512F {
			continue
		}

		opcodes = append(opcodes, r.Name)
	}

	sort.Strings(opcodes)
	return opcodes
}

func MakeBucketOne(incPath string, sweepDir string) ([]string, string, error) {
	infos, err := ExtractBucketOne(incPath)
	if err != nil {
		return nil, "", err
	}
	opcodes := FilterBucketOne(infos)

	bucketDir := filepath.Join(sweepDir, "buckets")
	if err := os.MkdirAll(bucketDir, 0755); err != nil {
		fmt.Fprintln(os.Stderr, "[BenchWire] Couldn't create bucket one directory for sweep. ", err)
	}

	bucketFile := filepath.Join(bucketDir, "bucketOne.txt")
	content := []byte(strings.Join(opcodes, "\n"))
	if err := os.WriteFile(bucketFile, content, 0644); err != nil {
		fmt.Fprintln(os.Stderr, "[BenchWire] Couldn't create bucket one. ", err)
		return opcodes, "", err
	}

	return opcodes, bucketFile, nil
}

// TODO: BucketTwo
func CreateSweepConfig(modeSelect []string, sweepDir string, binPath string, bucktOne []string) {
	presetFile := "sweepPreset.yaml"
	cfg, err := config.LoadYamlConfig(sweepDir, presetFile) 
	if err != nil {
        fmt.Fprintf(os.Stderr, "[BenchWire] Failed to load YAML config: %v\n", err)
    }

	fmt.Print(modeSelect[0])
	for _, mode := range modeSelect {
		for _, opcode := range bucktOne {
			opcodeFlag := "--opcode-name=" + opcode
			target := config.Target {
				Label: opcode + " " + mode,
				BinPath: binPath,
				Preset: []string{mode},
				Flags: []string{opcodeFlag},
			}
			cfg.SpecMatrix[0].Targets = append(cfg.SpecMatrix[0].Targets, target)
			cfg.SpecMatrix[0].Sequence = append(cfg.SpecMatrix[0].Sequence, []string{opcode + " " + mode})
		}
	}
	
	builtSweep, _ := yaml.Marshal(&cfg)
	sweepFile := filepath.Join(sweepDir, "sweep.yaml")
	err = os.WriteFile(sweepFile, builtSweep, 0644)
	if err != nil {
		fmt.Fprintln(os.Stderr, "[BenchWire] Couldn't create custom config for sweep.", err)
	}	
}

func BuildSweep(configDir string, mode string) (string, error) {
	var incPath string
	fmt.Print("[BenchWire] Enter your build/lib/Target/X86/X86GenInstrInfo.inc location: ")
	fmt.Scanln(&incPath)

	var binpath string
	fmt.Print("[BenchWire] Binary path: ")
	fmt.Scanln(&binpath)

	var modeSelect string
	fmt.Print("[BenchWire] Enter llvm-exegesis mode(s).\n")
	fmt.Print("[BenchWire] latency, uops, inverse_throughput, or full for all 3: ")
	fmt.Scanln(&modeSelect)

	sweepDir := filepath.Join(configDir, "configs", "sweep")
	bucketOne, _, err := MakeBucketOne(incPath, sweepDir)
	if err != nil {
		fmt.Fprintln(os.Stderr, "[BenchWire] Failed to extract bucket one:", err)
		return "", err
	}

	switch modeSelect {
		case "latency":
			CreateSweepConfig([]string{"Lat"}, sweepDir, binpath, bucketOne)
		case "uops":
			CreateSweepConfig([]string{"Uop"}, sweepDir, binpath, bucketOne)
		case "inverse_throughput":
			CreateSweepConfig([]string{"Thru"}, sweepDir, binpath, bucketOne)
		case "full":
			fmt.Print("Hit")
			CreateSweepConfig([]string{"Lat", "Uop", "Thru"}, sweepDir, binpath, bucketOne)
	}

	return sweepDir, nil
}
