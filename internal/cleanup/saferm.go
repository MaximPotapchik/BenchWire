package cleanup

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func SafeRemove(target, allowedDir string) error {
	if strings.TrimSpace(target) == "" {
		return fmt.Errorf("safe_rm: blocked, empty path")
	}
	abs, err := filepath.Abs(target)
	if err != nil {
		return fmt.Errorf("safe_rm: cannot resolve path: %s", target)
	}
	if resolved, err := filepath.EvalSymlinks(abs); err == nil {
		abs = resolved
	}
	absAllowed, _ := filepath.Abs(allowedDir)
	if abs != absAllowed && !strings.HasPrefix(abs, absAllowed+string(filepath.Separator)) {
		return fmt.Errorf("safe_rm: blocked, %s is not inside %s", target, allowedDir)
	}
	return os.RemoveAll(abs)
}
