package jl

import (
	"bufio"
	"encoding/json"
	"io"
)

type Parser struct {
	r       io.Reader
	scan    *bufio.Scanner
	printer EntryPrinter
}

func NewParser(r io.Reader, h EntryPrinter) *Parser {
	return &Parser{
		r:       r,
		scan:    bufio.NewScanner(r),
		printer: h,
	}
}

func (p *Parser) Consume() error {
	s := p.scan
	const maxBuffer int = 10 * 1024 * 1024 // 10MB line buffer
	buf := make([]byte, maxBuffer)
	s.Buffer(buf, maxBuffer)
	for s.Scan() {
		raw := s.Bytes()
		var partials map[string]json.RawMessage
		_ = json.Unmarshal(raw, &partials)
		message := &Entry{
			Partials:    partials,
			Raw:         raw,
		}
		p.printer.Print(message)
	}
	return p.scan.Err()
}

type EntryPrinter interface {
	Print(*Entry)
}

type Entry struct {
	Partials    map[string]json.RawMessage
	Raw         []byte
}
