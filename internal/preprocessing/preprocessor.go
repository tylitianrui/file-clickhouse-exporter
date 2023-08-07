package preprocessing

import (
	"sync"
)

type PreprocessingHandler struct {
	source map[string]string
	mu     sync.RWMutex
}

func (ph *PreprocessingHandler) WithData(source map[string]string) {
	ph.mu.Lock()
	defer ph.mu.Unlock()
	ph.source = source
}

func (ph *PreprocessingHandler) DO(format string) map[string]string {
	ph.mu.RLock()
	defer ph.mu.RUnlock()

	return nil
}
