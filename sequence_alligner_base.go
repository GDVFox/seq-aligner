package main

const gapByte = byte('-')

type action byte

const (
	letterAction action = iota
	firstGapAction
	secondGapAction
	zeroAction
)

// SequenceAlignerConfig набор параметров для конфигурации SequenceAligner
type SequenceAlignerConfig struct {
	AllowLocal      bool
	GapStartPenalty bool
	GapEndPenalty   bool
	GapPenalty      int
}

type sequenceAlignerBase struct {
	allowLocal      bool
	gapStartPenalty bool
	gapEndPenalty   bool
	gapPenalty      int
	scorer          Scorer
}

func (a *sequenceAlignerBase) getGapPenalty(i, max int) int {
	// если не штрафуем за gap в начале и находимся в начале какуй-либо последовательности,
	// то штраф за gap в этой последовательности нужно убрать.
	if !a.gapStartPenalty && i == 0 {
		return 0
	}

	// если не штрафуем за gap в конце и прошли какую-либо последовательность до конца,
	// то штраф за gap в этой последовательности нужно убрать.
	if !a.gapEndPenalty && i == max {
		return 0
	}

	return a.gapPenalty
}
