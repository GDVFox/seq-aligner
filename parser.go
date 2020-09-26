package main

import (
	"bufio"
	"io"
	"strings"
	"unicode"

	"github.com/pkg/errors"
)

// Possible parse errors
var (
	ErrBadHeader = errors.New("fasta parser: bad header")
)

// FastaParser parses a sequence of objects from reader
type FastaParser struct {
	reader *bufio.Reader
}

// NewFastaParser returns new FastaParser
func NewFastaParser(r io.Reader) *FastaParser {
	return &FastaParser{
		reader: bufio.NewReader(r),
	}
}

// Next gets next object from reader.
// Returns io.EOF if all objects were read.
func (p *FastaParser) Next() (*Sequence, error) {
	header, err := p.reader.ReadString('\n')
	if err != nil {
		return nil, err
	}

	descr, err := p.parseHeader(header)
	if err != nil {
		return nil, err
	}

	valueBuilder := &strings.Builder{}
	for {
		b, err := p.reader.ReadByte()
		if err != nil {
			// end of file
			if err == io.EOF {
				break
			}
			return nil, err
		}
		// end of current object
		if b == '>' {
			if err := p.reader.UnreadByte(); err != nil {
				return nil, err
			}
			break
		}
		// ignore newline
		if unicode.IsSpace(rune(b)) {
			continue
		}

		valueBuilder.WriteByte(b)
	}

	return &Sequence{
		Description: descr,
		Value:       valueBuilder.String(),
	}, nil
}

func (p *FastaParser) parseHeader(h string) (string, error) {
	if len(h) == 0 || h[0] != '>' {
		return "", ErrBadHeader
	}

	return strings.TrimSpace(h[1:]), nil
}
