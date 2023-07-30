package file_parser

import "sync"

var DefaultParserController = NewParserController()

type Parser interface {
	SetFormatString(s string) error
	SetFormat([]string) error
	Parse(s string) map[string]string
	ParseColumns(idx []string, s string) []string
}

type ParserController struct {
	allParsers map[string]Parser
	mu         sync.RWMutex
}

func NewParserController() *ParserController {
	return &ParserController{
		allParsers: map[string]Parser{},
	}
}

func (pc *ParserController) RegisterParser(k string, p Parser) {
	pc.mu.Lock()
	defer pc.mu.Unlock()
	pc.allParsers[k] = p
}

func (pc *ParserController) RemoveParser(k string) {
	pc.mu.Lock()
	defer pc.mu.Unlock()
	delete(pc.allParsers, k)
}

func (pc *ParserController) GetParser(k string) (Parser, bool) {
	pc.mu.RLock()
	defer pc.mu.RUnlock()
	p, exist := pc.allParsers[k]
	return p, exist
}
