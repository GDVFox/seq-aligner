package main

import "github.com/pkg/errors"

var (
	// ErrInvalidSymbol символ не из альфавита в последовательности
	ErrInvalidSymbol = errors.New("invalid symbol")
)

// Scorer оценивает разницу между символами
type Scorer interface {
	Score(a, b byte) int
}

// Adapter вспомогательная структура для работы с последовательностями некоторого алфавита
type Adapter interface {
	Scorer
	Validate(seq string) error
}

// MatrixAdapter адаптер на основе матрицы
type MatrixAdapter struct {
	inner   [][]int
	symbols map[byte]int
}

// Score оценить разницу по матрице.
// В случае, если a или b не принадлежат алфавиту паникует.
func (s *MatrixAdapter) Score(a, b byte) int {
	return s.inner[s.symbols[a]][s.symbols[b]]
}

// Validate проверяет строку на соответсвие алфавиту
func (s *MatrixAdapter) Validate(seq string) error {
	for _, a := range seq {
		if _, ok := s.symbols[byte(a)]; !ok {
			return ErrInvalidSymbol
		}
	}
	return nil
}

// NewDNAAdapter возвращает новый объект для работы с последовательностями нуклеотидов
func NewDNAAdapter() *MatrixAdapter {
	return &MatrixAdapter{
		inner: [][]int{
			{5, -4, -4, -4},
			{-4, 5, -4, -4},
			{-4, -4, 5, -4},
			{-4, -4, -4, 5},
		},
		symbols: map[byte]int{
			'A': 0,
			'T': 1,
			'G': 2,
			'C': 3,
		},
	}
}

// NewProteinAdapter возвращает новый объект для работы с последовательностями аминокислот
func NewProteinAdapter() *MatrixAdapter {
	return &MatrixAdapter{
		inner: [][]int{
			{4, -1, -2, -2, 0, -1, -1, 0, -2, -1, -1, -1, -1, -2, -1, 1, 0, -3, -2, 0},
			{-1, 5, 0, -2, -3, 1, 0, -2, 0, -3, -2, 2, -1, -3, -2, -1, -1, -3, -2, -3},
			{-2, 0, 6, 1, -3, 0, 0, 0, 1, -3, -3, 0, -2, -3, -2, 1, 0, -4, -2, -3},
			{-2, -2, 1, 6, -3, 0, 2, -1, -1, -3, -4, -1, -3, -3, -1, 0, -1, -4, -3, -3},
			{0, -3, -3, -3, 9, -3, -4, -3, -3, -1, -1, -3, -1, -2, -3, -1, -1, -2, -2, -1},
			{-1, 1, 0, 0, -3, 5, 2, -2, 0, -3, -2, 1, 0, -3, -1, 0, -1, -2, -1, -2},
			{-1, 0, 0, 2, -4, 2, 5, -2, 0, -3, -3, 1, -2, -3, -1, 0, -1, -3, -2, -2},
			{0, -2, 0, -1, -3, -2, -2, 6, -2, -4, -4, -2, -3, -3, -2, 0, -2, -2, -3, -3},
			{-2, 0, 1, -1, -3, 0, 0, -2, 8, -3, -3, -1, -2, -1, -2, -1, -2, -2, 2, -3},
			{-1, -3, -3, -3, -1, -3, -3, -4, -3, 4, 2, -3, 1, 0, -3, -2, -1, -3, -1, 3},
			{-1, -2, -3, -4, -1, -2, -3, -4, -3, 2, 4, -2, 2, 0, -3, -2, -1, -2, -1, 1},
			{-1, 2, 0, -1, -3, 1, 1, -2, -1, -3, -2, 5, -1, -3, -1, 0, -1, -3, -2, -2},
			{-1, -1, -2, -3, -1, 0, -2, -3, -2, 1, 2, -1, 5, 0, -2, -1, -1, -1, -1, 1},
			{-2, -3, -3, -3, -2, -3, -3, -3, -1, 0, 0, -3, 0, 6, -4, -2, -2, 1, 3, -1},
			{-1, -2, -2, -1, -3, -1, -1, -2, -2, -3, -3, -1, -2, -4, 7, -1, -1, -4, -3, -2},
			{1, -1, 1, 0, -1, 0, 0, 0, -1, -2, -2, 0, -1, -2, -1, 4, 1, -3, -2, -2},
			{0, -1, 0, -1, -1, -1, -1, -2, -2, -1, -1, -1, -1, -2, -1, 1, 5, -2, -2, 0},
			{-3, -3, -4, -4, -2, -2, -3, -2, -2, -3, -2, -3, -1, 1, -4, -3, -2, 11, 2, -3},
			{-2, -2, -2, -3, -2, -1, -2, -3, 2, -1, -1, -2, -1, 3, -3, -2, -2, 2, 7, -1},
			{0, -3, -3, -3, -1, -2, -2, -3, -3, 3, 1, -2, 1, -1, -2, -2, 0, -3, -1, 4},
		},
		symbols: map[byte]int{
			'A': 0,
			'R': 1,
			'N': 2,
			'D': 3,
			'C': 4,
			'Q': 5,
			'E': 6,
			'G': 7,
			'H': 8,
			'I': 9,
			'L': 10,
			'K': 11,
			'M': 12,
			'F': 13,
			'P': 14,
			'S': 15,
			'T': 16,
			'W': 17,
			'Y': 18,
			'V': 19,
		},
	}
}

// DefaultAdapter объект по умолчанию подходит для работы с произвольными последовательностями
type DefaultAdapter struct {
}

// NewDefaultAdapter возвращает новый объект DNAScorer
func NewDefaultAdapter() *DefaultAdapter {
	return &DefaultAdapter{}
}

// Score если символы сопадают — 1, иначе — 0.
func (s *DefaultAdapter) Score(a, b byte) int {
	if a == b {
		return 1
	}
	return -1
}

// Validate любая цепочка символов без '-' (зарезервированный символ), считается валидной
func (s *DefaultAdapter) Validate(str string) error {
	for _, a := range str {
		if a == '-' {
			return ErrInvalidSymbol
		}
	}

	return nil
}

func buildAdapter(mode string) Adapter {
	switch mode {
	case dnaMode:
		return NewDNAAdapter()
	case proteinMode:
		return NewProteinAdapter()
	}

	return NewDefaultAdapter()
}

func validate(a Adapter, seqs []*Sequence) error {
	for i, seq := range seqs {
		if err := a.Validate(seq.Value); err != nil {
			return errors.Wrapf(err, "sequence %d", i)
		}
	}
	return nil
}
