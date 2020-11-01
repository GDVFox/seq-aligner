package main

import (
	"strings"
)

type coord struct {
	i int
	j int
}

// SequenceAlignerMem вспомогательный объект для глобального выравнивания оптимального по памяти
type SequenceAlignerMem struct {
	sequenceAlignerBase

	upBuffer   []int
	downBuffer []int
}

// NewSequenceAlignerMem возвращает новый объект SequenceAlignerMem
func NewSequenceAlignerMem(cfg *SequenceAlignerConfig, scorer Scorer) *SequenceAlignerMem {
	return &SequenceAlignerMem{
		sequenceAlignerBase: sequenceAlignerBase{
			gapStartPenalty: cfg.GapStartPenalty,
			gapEndPenalty:   cfg.GapEndPenalty,
			gapPenalty:      cfg.GapPenalty,
			scorer:          scorer,
		},
	}
}

// Align производит оптимальное глобальное выравнивание двух последовательностей
func (a *SequenceAlignerMem) Align(str1, str2 string) (string, string, int) {
	a.upBuffer = make([]int, len(str2)+1)
	a.downBuffer = make([]int, len(str2)+1)

	actions, score := a.findActions(str1, str2, &coord{0, 0}, &coord{len(str1), len(str2)})
	alignedStr1, alignedStr2 := &strings.Builder{}, &strings.Builder{}

	i, j, actionIndex := 0, 0, 0
	for i < len(str1) || j < len(str2) {
		switch actions[actionIndex] {
		case letterAction:
			alignedStr1.WriteByte(str1[i])
			alignedStr2.WriteByte(str2[j])
			i++
			j++
		case firstGapAction:
			alignedStr1.WriteByte(gapByte)
			alignedStr2.WriteByte(str2[j])
			j++
		case secondGapAction:
			alignedStr1.WriteByte(str1[i])
			alignedStr2.WriteByte(gapByte)
			i++
		}
		actionIndex++
	}

	return alignedStr1.String(), alignedStr2.String(), score
}

func (a *SequenceAlignerMem) findActions(str1, str2 string, f, t *coord) ([]action, int) {
	if f.i == t.i {
		score := 0
		res := make([]action, t.j-f.j)
		for i := 0; i < t.j-f.j; i++ {
			score += a.getGapPenalty(f.i, len(str1))
			res[i] = firstGapAction
		}
		return res, score
	}

	size := (t.i - f.i)
	upSize := size / 2
	downSize := (size - (size+1)%2) / 2

	a.findUp(str1, str2, f, &coord{f.i + upSize, t.j})
	a.findDown(str1, str2, &coord{t.i - downSize, f.j}, t)

	i, j := f.i+upSize, f.j
	act, v := secondGapAction, a.upBuffer[j]+a.downBuffer[j]+a.getGapPenalty(i, len(str1))

	for k := f.j; k <= t.j; k++ {
		current := a.upBuffer[k] + a.downBuffer[k] + a.getGapPenalty(k, len(str2))
		if current > v {
			j, v = k, current
		}
	}

	for k := f.j; k < t.j; k++ {
		current := a.upBuffer[k] + a.downBuffer[k+1] + a.scorer.Score(str1[i], str2[k])
		if current > v {
			j, v = k, current
			act = letterAction
		}
	}

	fNext := &coord{i + 1, j}
	if act == letterAction {
		fNext.j++
	}

	tNext := &coord{i, j}

	part1, _ := a.findActions(str1, str2, f, tNext)
	part2, _ := a.findActions(str1, str2, fNext, t)

	res := make([]action, 0)
	res = append(res, part1...)
	res = append(res, act)
	res = append(res, part2...)

	return res, v
}

func (a *SequenceAlignerMem) findUp(str1, str2 string, f, t *coord) {
	a.upBuffer[f.j] = 0
	for j := f.j + 1; j <= t.j; j++ {
		a.upBuffer[j] = a.upBuffer[j-1] + a.getGapPenalty(f.i, len(str1))
	}

	var tmp int
	for i := f.i; i < t.i; i++ {
		tmp, a.upBuffer[f.j] = a.upBuffer[f.j], a.upBuffer[f.j]+a.getGapPenalty(f.j, len(str2))
		for j := f.j + 1; j <= t.j; j++ {
			val, _ := MaxOfThreeInt(
				tmp+a.scorer.Score(str1[i], str2[j-1]),
				a.upBuffer[j-1]+a.getGapPenalty(i, len(str1)),
				a.upBuffer[j]+a.getGapPenalty(j, len(str2)),
			)

			tmp, a.upBuffer[j] = a.upBuffer[j], val
		}
	}
}

func (a *SequenceAlignerMem) findDown(str1, str2 string, f, t *coord) {
	a.downBuffer[t.j] = 0
	for j := t.j - 1; j >= f.j; j-- {
		a.downBuffer[j] = a.downBuffer[j+1] + a.getGapPenalty(t.i, len(str1))
	}

	var tmp int
	for i := t.i; i > f.i; i-- {
		tmp, a.downBuffer[t.j] = a.downBuffer[t.j], a.downBuffer[t.j]+a.getGapPenalty(t.j, len(str2))
		for j := t.j - 1; j >= f.j; j-- {
			val, _ := MaxOfThreeInt(
				tmp+a.scorer.Score(str1[i-1], str2[j]),
				a.downBuffer[j+1]+a.getGapPenalty(i, len(str1)),
				a.downBuffer[j]+a.getGapPenalty(j, len(str2)),
			)

			tmp, a.downBuffer[j] = a.downBuffer[j], val
		}
	}
}
