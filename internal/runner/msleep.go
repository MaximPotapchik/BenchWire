package runner

import (
	"benchwire/internal/config"
	"time"
	"math/rand/v2"
	"strings"
	"fmt"
	"strconv"
	"unicode"
)

func msleep(cooldown config.CooldownTimer) (int64, error) {
	split := strings.IndexFunc(cooldown.Value, unicode.IsLetter)
	if split == -1 {
		return 0, fmt.Errorf("cooldown %q: missing unit, expected ms or s", cooldown.Value)
	}
	number, precision := cooldown.Value[:split], cooldown.Value[split:]

	duration, err := strconv.Atoi(number)
	if err != nil {
		return 0, fmt.Errorf("cooldown %q: invalid number: %w", cooldown.Value, err)
	}

	switch precision {
		case "ms":
			duration *= 1_000_000
		case "s":
			duration *= 1_000_000_000
		default:
			return 0, fmt.Errorf("cooldown %q: unknown precision %q expected ms or s", cooldown.Value, precision)
	}

	const minNs = 50_000
	if duration < minNs {
		return 0, fmt.Errorf("cooldown %dns below minimum 50us", duration)
	}

	actual := duration
	if cooldown.Randomized {
		actual = minNs + rand.IntN(duration - minNs)
	}

	time.Sleep(time.Duration(actual) * time.Nanosecond)
	return int64(actual), nil
}
