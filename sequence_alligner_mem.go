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
	gapPenalty int
	scorer     Scorer

	upBuffer   []int
	downBuffer []int
}

// NewSequenceAlignerMem возвращает новый объект SequenceAlignerMem
func NewSequenceAlignerMem(gPen int, scorer Scorer) *SequenceAlignerMem {
	return &SequenceAlignerMem{
		gapPenalty: gPen,
		scorer:     scorer,
	}
}

// Align производит оптимальное глобальное выравнивание двух последовательностей
func (a *SequenceAlignerMem) Align(str1, str2 string) (string, string, int) {
	a.upBuffer = make([]int, len(str1)+1)
	a.downBuffer = make([]int, len(str1)+1)

	actions := a.findActions(str1, str2, &coord{0, 0}, &coord{len(str1), len(str2)})
	alignedStr1, alignedStr2 := &strings.Builder{}, &strings.Builder{}
	score := 0

	i, j, actionIndex := 0, 0, 0
	for i < len(str1) || j < len(str2) {
		switch actions[actionIndex] {
		case letterAction:
			alignedStr1.WriteByte(str1[i])
			alignedStr2.WriteByte(str2[j])
			score += a.scorer.Score(str1[i], str2[j])
			i++
			j++
		case firstGapAction:
			alignedStr1.WriteByte(gapByte)
			alignedStr2.WriteByte(str2[j])
			score += a.gapPenalty
			j++
		case secondGapAction:
			alignedStr1.WriteByte(str1[i])
			alignedStr2.WriteByte(gapByte)
			score += a.gapPenalty
			i++
		}
		actionIndex++
	}

	return alignedStr1.String(), alignedStr2.String(), score
}

func (a *SequenceAlignerMem) findActions(str1, str2 string, f, t *coord) []action {
	if f.j == t.j {
		res := make([]action, t.i-f.i)
		for i := 0; i < t.i-f.i; i++ {
			res[i] = secondGapAction
		}
		return res
	}

	size := (t.j - f.j)
	upSize := size / 2
	downSize := (size - (size+1)%2) / 2

	a.findUp(str1, str2, f, &coord{t.i, f.j + upSize})
	a.findDown(str1, str2, &coord{f.i, t.j - downSize}, t)

	i, j := f.i, f.j+upSize
	act, v := firstGapAction, a.upBuffer[i]+a.downBuffer[i]+a.gapPenalty

	for k := f.i; k <= t.i; k++ {
		current := a.upBuffer[k] + a.downBuffer[k] + a.gapPenalty
		if current > v {
			i, v = k, current
		}
	}

	for k := f.i; k < t.i; k++ {
		current := a.upBuffer[k] + a.downBuffer[k+1] + a.scorer.Score(str1[k], str2[j])
		if current > v {
			i, v = k, current
			act = letterAction
		}
	}

	fNext := &coord{i, j + 1}
	if act == letterAction {
		fNext.i++
	}

	tNext := &coord{i, j}

	res := make([]action, 0)
	res = append(res, a.findActions(str1, str2, f, tNext)...)
	res = append(res, act)
	res = append(res, a.findActions(str1, str2, fNext, t)...)
	return res
}

func (a *SequenceAlignerMem) findUp(str1, str2 string, f, t *coord) {
	a.upBuffer[f.i] = 0
	for i := f.i + 1; i <= t.i; i++ {
		a.upBuffer[i] = a.upBuffer[i-1] + a.gapPenalty
	}

	var tmp int
	for j := f.j; j < t.j; j++ {
		tmp, a.upBuffer[f.i] = a.upBuffer[f.i], a.upBuffer[f.i]+a.gapPenalty
		for i := f.i + 1; i <= t.i; i++ {
			val, _ := MaxOfThreeInt(
				tmp+a.scorer.Score(str1[i-1], str2[j]),
				a.upBuffer[i]+a.gapPenalty,
				a.upBuffer[i-1]+a.gapPenalty,
			)

			tmp, a.upBuffer[i] = a.upBuffer[i], val
		}
	}
}

func (a *SequenceAlignerMem) findDown(str1, str2 string, f, t *coord) {
	a.downBuffer[t.i] = 0
	for i := t.i - 1; i >= f.i; i-- {
		a.downBuffer[i] = a.downBuffer[i+1] + a.gapPenalty
	}

	var tmp int
	for j := t.j; j > f.j; j-- {
		tmp, a.downBuffer[t.i] = a.downBuffer[t.i], a.downBuffer[t.i]+a.gapPenalty
		for i := t.i - 1; i >= f.i; i-- {
			val, _ := MaxOfThreeInt(
				tmp+a.scorer.Score(str1[i], str2[j-1]),
				a.downBuffer[i]+a.gapPenalty,
				a.downBuffer[i+1]+a.gapPenalty,
			)

			tmp, a.downBuffer[i] = a.downBuffer[i], val
		}
	}
}
