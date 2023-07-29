package file_parser

import "sync"

type Parser interface {
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
