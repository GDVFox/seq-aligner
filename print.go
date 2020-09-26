package main

import (
	"errors"
	"io"

	"github.com/fatih/color"
)

var (
	// ErrNotAligned возвращается, когда длина последовательностей разная
	ErrNotAligned = errors.New("aligned write: sequences are not aligned")
)

// WriteAlignedDefault запись выровненных последовательностей
// в стандартном формате с переносом каждые lineLength символов
func WriteAlignedDefault(w io.Writer, lineLength int, a, b string) error {
	if len(a) == 0 {
		return nil
	}
	if len(a) != len(b) {
		return ErrNotAligned
	}

	seqLen := len(a)
	l, r := 0, MinInt(seqLen, lineLength)

	for l < seqLen {
		io.WriteString(w, "seq1: ")
		io.WriteString(w, a[l:r])
		io.WriteString(w, "\n")
		io.WriteString(w, "seq2: ")
		io.WriteString(w, b[l:r])
		io.WriteString(w, "\n")

		l, r = r, MinInt(seqLen, r+lineLength)
	}

	return nil
}

// WritePretty выводит в разноцветном формате
// в одну строчку с дополнительными символами
func WritePretty(w io.Writer, a, b string) error {
	if len(a) == 0 {
		return nil
	}
	if len(a) != len(b) {
		return ErrNotAligned
	}
	seqLen := len(a)

	redAdapter := color.New(color.FgRed)     // gap
	greenAdapter := color.New(color.FgGreen) // match
	blueAdapter := color.New(color.FgBlue)   // missmatch
	getAdapter := func(b1, b2 byte) *color.Color {
		if b1 == gapByte || b2 == gapByte {
			return redAdapter
		} else if b1 == b2 {
			return greenAdapter
		} else {
			return blueAdapter
		}
	}
	symbolAdapter := func(b1, b2 byte) string {
		if b1 == gapByte || b2 == gapByte {
			return " "
		} else if b1 == b2 {
			return "*"
		} else {
			return "|"
		}
	}

	io.WriteString(w, "seq1: ")
	for i := 0; i < seqLen; i++ {
		getAdapter(a[i], b[i]).Fprint(w, string(a[i]))
	}
	io.WriteString(w, "\n")
	io.WriteString(w, "      ") // len("seq1: ")
	for i := 0; i < seqLen; i++ {
		getAdapter(a[i], b[i]).Fprint(w, symbolAdapter(a[i], b[i]))
	}
	io.WriteString(w, "\n")
	io.WriteString(w, "seq2: ")
	for i := 0; i < seqLen; i++ {
		getAdapter(a[i], b[i]).Fprint(w, string(b[i]))
	}
	io.WriteString(w, "\n")

	return nil
}
