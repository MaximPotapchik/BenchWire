package cleanup

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func resolveExisting(path string) (string, error) {
	if resolved, err := filepath.EvalSymlinks(path); err == nil {
		return resolved, nil
	}

	parent := filepath.Dir(path)
	resolvedParent, err := filepath.EvalSymlinks(parent)
	if err != nil {
		return path, nil
	}

	return filepath.Join(resolvedParent, filepath.Base(path)), nil
}

func SafeRemove(target, allowedDir string) error {
	if strings.TrimSpace(target) == "" {
		return fmt.Errorf("[BenchWire] SafeRemove: blocked, empty path")
	}
	
	abs, err := filepath.Abs(target)
	if err != nil {
		return fmt.Errorf("[BenchWire] SafeRemove: cannot resolve path: %s", target)
	}

	abs, err = resolveExisting(abs)
	if err != nil {
		return fmt.Errorf("[BenchWire] SafeRemove: cannot resolve path: %s", target)
	}

	absAllowed, err := filepath.Abs(allowedDir)
	if err != nil {
		return fmt.Errorf("[BenchWire] SafeRemove: cannot resolve allowed dir: %s", allowedDir)
	}

	absAllowed, err = resolveExisting(absAllowed)
	if err != nil {
		return fmt.Errorf("[BenchWire] SafeRemove: cannot resolve allowed dir: %s", allowedDir)
	}

	if abs != absAllowed && !strings.HasPrefix(abs, absAllowed+string(filepath.Separator)) {
		return fmt.Errorf("[BenchWire] SafeRemove: blocked, %s is not inside %s", target, allowedDir)
	}

	return os.RemoveAll(abs)
}
