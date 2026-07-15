package exegesis

import (
	"fmt"
	"os"
	"regexp"
	"strings"
)

func Format(methodology string, rawArgs []string) ([]string, error) {
	var formatted []string

	if bin, ok := os.LookupEnv("EXEGESISBIN"); ok {
		formatted = append(formatted, bin)
	} else if binA, ok := os.LookupEnv("EXEGESISBINA"); ok {
		binB := os.Getenv("EXEGESISBINB")
		formatted = append(formatted, binA+"|"+binB)
	} else {
		return nil, fmt.Errorf("no EXEGESISBIN specified, add your exegesis bin locations into .env")
	}


	benchmarksFile := regexp.MustCompile(`^-{0,2}benchmarks-file[AB]?=?$`)

	for _, opt := range rawArgs {
		lhs, _, found := strings.Cut(opt, "=")
		if !found {
			continue
		}
		lhs += "="

		if benchmarksFile.MatchString(lhs) {
			continue
		}
		if strings.HasPrefix(lhs, "--") {
			formatted = append(formatted, opt)
		}
	}
	return formatted, nil
}
