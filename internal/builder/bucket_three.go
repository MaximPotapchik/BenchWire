package builder

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"golang.org/x/sys/cpu"
)

// TODO: For now this is seperate from two. Further precision in opcode choice
// is needed.  We can unify ExtractBucket functions and just use a struct
// holding skipFlags, and other variables for precision. 
func ExtractBucketThree(incPath string) ([]OpcodeInfo, error) {
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
		"UnmodeledSideEffects", "MCID::Pseudo",
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

func FilterBucketThree(info []OpcodeInfo) []string {
	var opcodes []string

	for _, r := range info {
		if r.Skipped {
			continue
		}

		if r.Encoding == EncodingEVEX && !cpu.X86.HasAVX512F {
			continue
		}

		// Filters pseudo instructions that make it past the other skipped flags.
		if r.TSFlag == 0xc00000 || r.TSFlag == 0x1800000 || r.TSFlag == 0x1400000 ||
		   r.TSFlag == 0x400000 || r.TSFlag == 0x800000 {
			continue
		}

		opcodes = append(opcodes, r.Name)
	}

	sort.Strings(opcodes)
	return opcodes
}

func MakeBucketThree(incPath string, sweepDir string) ([]string, string, error) {
	infos, err := ExtractBucketThree(incPath)
	if err != nil {
		return nil, "", err
	}
	opcodes := FilterBucketThree(infos)

	bucketDir := filepath.Join(sweepDir, "buckets")
	if err := os.MkdirAll(bucketDir, 0755); err != nil {
		fmt.Fprintln(os.Stderr, "[BenchWire] Couldn't create bucket three directory for sweep. ", err)
	}

	bucketFile := filepath.Join(bucketDir, "bucketThree.txt")
	content := []byte(strings.Join(opcodes, "\n"))
	if err := os.WriteFile(bucketFile, content, 0644); err != nil {
		fmt.Fprintln(os.Stderr, "[BenchWire] Couldn't create bucket three. ", err)
		return opcodes, "", err
	}

	return opcodes, bucketFile, nil
}

