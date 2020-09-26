package main

import (
	"errors"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
)

// ErrWrongNumberOfFiles возвращается
var (
	ErrWrongNumberOfFiles = errors.New("expected one or two sequences files")
)

const (
	dnaMode     = "dna"
	proteinMode = "protein"
	defaultMode = "default"
)

var (
	gapValue int
	mode     string

	pretty     bool
	lineLength int
	outputFile string

	startPenalty bool
	endPenalty   bool
)

func init() {
	flag.IntVar(&gapValue, "gap", -2, "gap penalty size")
	flag.StringVar(&mode, "mode", defaultMode, "(dna|protein|default) alphabet and score table switch")

	flag.BoolVar(&pretty, "pretty", false, "enables pretty output mode")
	flag.IntVar(&lineLength, "line", 100, "line length for default output mode")
	flag.StringVar(&outputFile, "out", "", "output file name")

	flag.BoolVar(&startPenalty, "spen", false, "enables start gap penalty")
	flag.BoolVar(&endPenalty, "epen", false, "enables end gap penalty")
}

// Sequence описывает последовательность из fasta файла
type Sequence struct {
	Description string
	Value       string
}

func main() {
	flag.Parse()

	out := os.Stdout
	if outputFile != "" {
		var err error
		out, err = os.Create(outputFile)
		if err != nil {
			log.Fatalf("can not open output file: %s", err)
		}
		defer out.Close()
	}

	sequences, err := loadSequences(flag.Args())
	if err != nil {
		log.Fatalf("can not read sequences: %s", err)
	}

	adapter := buildAdapter(mode)
	if err := validate(adapter, sequences); err != nil {
		log.Fatal(err)
	}

	cfg := &SequenceAlignerConfig{
		GapPenalty:      gapValue,
		GapStartPenalty: startPenalty,
		GapEndPenalty:   endPenalty,
	}
	aligner := NewSequenceAligner(cfg, adapter)
	aligned1, aligned2, score := aligner.Align(sequences[0].Value, sequences[1].Value)

	if pretty && out != os.Stdout {
		io.WriteString(out, "WARN: can not use '--pretty' with file output!\n")
		pretty = false
	}
	if pretty {
		WritePretty(out, aligned1, aligned2)
	} else {
		WriteAlignedDefault(out, 100, aligned1, aligned2)
	}
	fmt.Fprintf(out, "Score: %d\n", score)
}
