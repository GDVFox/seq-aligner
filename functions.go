package main

// MaxInt возвращает максимум из двух int
func MaxInt(a, b int) int {
	if a > b {
		return a
	}

	return b
}

// MinInt возвращает минимум из двух int
func MinInt(a, b int) int {
	if a <= b {
		return a
	}

	return b
}

// MaxOfThreeInt возвращает максимум из 3х int и порядковый номер этого int c 0
func MaxOfThreeInt(a, b, c int) (int, int) {
	if (a >= b) && (a >= c) {
		return a, 0
	}

	if (b >= a) && (b >= c) {
		return b, 1
	}

	return c, 2
}

// Reverse разворот строки
func Reverse(s string) string {
	runes := []rune(s)
	for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
		runes[i], runes[j] = runes[j], runes[i]
	}
	return string(runes)
}
