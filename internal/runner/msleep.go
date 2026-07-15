package runner

import "time"

func msleep(ms int) {
	time.Sleep(time.Duration(ms) * time.Millisecond)
}
