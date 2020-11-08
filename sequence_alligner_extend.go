package main

import (
	"strings"
)

// SequenceAlignerExtendConfig набор параметров для конфигурации SequenceAlignerExtend.
type SequenceAlignerExtendConfig struct {
	SequenceAlignerConfig
	ExtendGapPenalty int
}

// SequenceAlignerExtend вспомогательный объект для глобального выравнивания,
// с возможностью разного штрафа за расшиерние gap.
type SequenceAlignerExtend struct {
	sequenceAlignerBase
	extendGapPenalty int
}

// NewSequenceAlignerExtend возвращает новый объект SequenceAlignerExtend.
func NewSequenceAlignerExtend(cfg *SequenceAlignerExtendConfig, scorer Scorer) *SequenceAlignerExtend {
	return &SequenceAlignerExtend{
		sequenceAlignerBase: sequenceAlignerBase{
			gapStartPenalty: cfg.GapStartPenalty,
			gapEndPenalty:   cfg.GapEndPenalty,
			gapPenalty:      cfg.GapPenalty,
			scorer:          scorer,
		},
		extendGapPenalty: cfg.ExtendGapPenalty,
	}
}

// Align производит оптимальное глобальное выравнивание двух последовательностей
func (a *SequenceAlignerExtend) Align(str1, str2 string) (string, string, int) {
	actions, currentAction, score := a.findActions(str1, str2)
	alignedStr1, alignedStr2 := &strings.Builder{}, &strings.Builder{}

	i, j := len(str1), len(str2)
	for {
		if i == 0 && j == 0 {
			break
		}

		nextAction := action((actions[i][j] >> (currentAction * 2)) & 0b11)
		switch currentAction {
		case letterAction:
			alignedStr1.WriteByte(str1[i-1])
			alignedStr2.WriteByte(str2[j-1])
			i--
			j--
		case firstGapAction:
			alignedStr1.WriteByte(gapByte)
			alignedStr2.WriteByte(str2[j-1])
			j--
		case secondGapAction:
			alignedStr1.WriteByte(str1[i-1])
			alignedStr2.WriteByte(gapByte)
			i--
		}

		currentAction = nextAction
	}

	return Reverse(alignedStr1.String()), Reverse(alignedStr2.String()), score
}

func (a *SequenceAlignerExtend) findActions(str1, str2 string) ([][]byte, action, int) {
	match, insetion, deletion, actions := a.buildExtendMatrices(len(str1)+1, len(str2)+1)
	for i := 1; i <= len(str1); i++ {
		for j := 1; j <= len(str2); j++ {
			var indexMatch, indexInsertion, indexDeletion int
			match[i][j], indexMatch = MaxOfThreeInt(
				match[i-1][j-1]+a.scorer.Score(str1[i-1], str2[j-1]),
				insetion[i-1][j-1]+a.scorer.Score(str1[i-1], str2[j-1]),
				deletion[i-1][j-1]+a.scorer.Score(str1[i-1], str2[j-1]),
			)
			insetion[i][j], indexInsertion = MaxOfThreeInt(
				match[i][j-1]+a.getGapPenalty(i, len(str1), a.gapPenalty),
				insetion[i][j-1]+a.getGapPenalty(i, len(str1), a.extendGapPenalty),
				deletion[i][j-1]+a.getGapPenalty(i, len(str1), a.gapPenalty),
			)
			deletion[i][j], indexDeletion = MaxOfThreeInt(
				match[i-1][j]+a.getGapPenalty(j, len(str2), a.gapPenalty),
				insetion[i-1][j]+a.getGapPenalty(j, len(str2), a.gapPenalty),
				deletion[i-1][j]+a.getGapPenalty(j, len(str2), a.extendGapPenalty),
			)

			actions[i][j] = byte(indexDeletion)<<4 | byte(indexInsertion)<<2 | byte(indexMatch)
		}
	}

	score, index := MaxOfThreeInt(match[len(str1)][len(str2)], insetion[len(str1)][len(str2)], deletion[len(str1)][len(str2)])
	return actions, action(index), score
}

func (a *SequenceAlignerExtend) buildExtendMatrices(rowCount, colCount int) ([][]int, [][]int, [][]int, [][]byte) {
	match := make([][]int, rowCount)
	insetion := make([][]int, rowCount)
	deletion := make([][]int, rowCount)
	actions := make([][]byte, rowCount)
	for i := 0; i < rowCount; i++ {
		match[i] = make([]int, colCount)
		insetion[i] = make([]int, colCount)
		deletion[i] = make([]int, colCount)
		actions[i] = make([]byte, colCount)
	}

	infinity := 2*a.gapPenalty + (rowCount+colCount)*a.extendGapPenalty + 1
	match[0][0] = 0
	insetion[0][0] = infinity
	deletion[0][0] = infinity

	for i := 1; i < rowCount; i++ {
		match[i][0] = infinity
		insetion[i][0] = infinity
		deletion[i][0] = a.getGapPenalty(0, colCount, a.gapPenalty+(i-1)*a.extendGapPenalty)
		actions[i][0] = byte(secondGapAction)<<4 | byte(secondGapAction)<<2 | byte(secondGapAction)
	}

	for j := 1; j < colCount; j++ {
		match[0][j] = infinity
		insetion[0][j] = a.getGapPenalty(0, rowCount, a.gapPenalty+(j-1)*a.extendGapPenalty)
		deletion[0][j] = infinity
		actions[0][j] = byte(firstGapAction)<<4 | byte(firstGapAction)<<2 | byte(firstGapAction)
	}

	return match, insetion, deletion, actions
}

func (a *SequenceAlignerExtend) getGapPenalty(i, max, potentialPenalty int) int {
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

	return potentialPenalty
}
