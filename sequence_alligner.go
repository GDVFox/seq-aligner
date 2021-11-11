package main

import (
	"strings"
)

// SequenceAligner вспомогательный объект для глобального выравнивания
type SequenceAligner struct {
	sequenceAlignerBase
}

// NewSequenceAligner возвращает новый объект SequenceAligner
func NewSequenceAligner(cfg *SequenceAlignerConfig, scorer Scorer) *SequenceAligner {
	return &SequenceAligner{
		sequenceAlignerBase: sequenceAlignerBase{
			allowLocal:      cfg.AllowLocal,
			gapStartPenalty: cfg.GapStartPenalty,
			gapEndPenalty:   cfg.GapEndPenalty,
			gapPenalty:      cfg.GapPenalty,
			scorer:          scorer,
		},
	}
}

// Align производит оптимальное глобальное выравнивание двух последовательностей
func (a *SequenceAligner) Align(str1, str2 string) (string, string, int) {
	actions, score := a.findActions(str1, str2)
	alignedStr1, alignedStr2 := &strings.Builder{}, &strings.Builder{}

	i, j := len(actions)-1, len(actions[0])-1

WriteCycle:
	for {
		if i == 0 && j == 0 {
			break WriteCycle
		}

		switch actions[i][j] {
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
		case zeroAction:
			break WriteCycle
		}
	}

	return Reverse(alignedStr1.String()), Reverse(alignedStr2.String()), score
}

func (a *SequenceAligner) findActions(str1, str2 string) ([][]action, int) {
	dp, actions := a.buildBaseMatrices(len(str1)+1, len(str2)+1)
	for i := 1; i <= len(str1); i++ {
		for j := 1; j <= len(str2); j++ {
			val, indx := MaxOfThreeInt(
				dp[i-1][j-1]+a.scorer.Score(str1[i-1], str2[j-1]), // i-1 и j-1 потому что с 1
				dp[i][j-1]+a.getGapPenalty(i, len(str1)),
				dp[i-1][j]+a.getGapPenalty(j, len(str2)),
			)
			if a.allowLocal && val < 0 {
				val = 0
				indx = int(zeroAction)
			}

			dp[i][j] = val
			actions[i][j] = action(indx)
		}
	}

	maxI, maxJ := len(str1), len(str2)
	if a.allowLocal {
		maxVal := -int(^uint(0)>>1) - 1 // minimal int value
		for i := 0; i < len(dp); i++ {
			for j := 0; j < len(dp[i]); j++ {
				if dp[i][j] >= maxVal {
					maxI, maxJ = i, j
					maxVal = dp[i][j]
				}
			}
		}

		actions = actions[:maxI+1]
		for i := 0; i < len(actions); i++ {
			actions[i] = actions[i][:maxJ+1]
		}
	}

	return actions, dp[maxI][maxJ]
}

func (a *SequenceAligner) buildBaseMatrices(rowCount, colCount int) ([][]int, [][]action) {
	dp := make([][]int, rowCount)
	actions := make([][]action, rowCount)
	for i := 0; i < rowCount; i++ {
		dp[i] = make([]int, colCount)
		actions[i] = make([]action, colCount)
	}

	dp[0][0] = 0
	// если за gap в начале не штрафуем, то не нужно пердвычислять границу из gap
	for i := 1; i < colCount; i++ {
		dp[0][i] = dp[0][i-1] + a.getGapPenalty(0, rowCount)
		actions[0][i] = firstGapAction
	}
	for i := 1; i < rowCount; i++ {
		dp[i][0] = dp[i-1][0] + a.getGapPenalty(0, colCount)
		actions[i][0] = secondGapAction
	}

	return dp, actions
}
