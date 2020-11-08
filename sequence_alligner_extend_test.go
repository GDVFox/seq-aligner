package main

import (
	"testing"

	"github.com/stretchr/testify/suite"
)

type SequenceAlignerExtendTestSuite struct {
	suite.Suite
	aligner *SequenceAlignerExtend
}

func (s *SequenceAlignerExtendTestSuite) SetupTest() {
	cfg := &SequenceAlignerExtendConfig{
		SequenceAlignerConfig: SequenceAlignerConfig{
			GapPenalty: -10,
		},
		ExtendGapPenalty: -1,
	}

	s.aligner = NewSequenceAlignerExtend(cfg, NewDNAAdapter())
}

func (s *SequenceAlignerExtendTestSuite) TestAlign() {
	for _, c := range []struct {
		a              string
		b              string
		enableStartPen bool
		enableEndPen   bool
		expA           string
		expB           string
		expectedScore  int
	}{
		// тест из задания
		{
			a:              "AT",
			b:              "G",
			enableStartPen: true,
			enableEndPen:   true,
			expA:           "AT",
			expB:           "-G",
			expectedScore:  -14,
		},
		// на штрафуем начало
		{
			a:              "AT",
			b:              "G",
			enableStartPen: false,
			enableEndPen:   true,
			expA:           "AT",
			expB:           "-G",
			expectedScore:  -4,
		},
		// на штрафуем конец
		{
			a:              "AT",
			b:              "G",
			enableStartPen: true,
			enableEndPen:   false,
			expA:           "AT",
			expB:           "G-",
			expectedScore:  -4,
		},
		// так как концы не штрафуются и все символы разные, самым выгодным решением является разделить 2 строки
		{
			a:              "AT",
			b:              "G",
			enableStartPen: false,
			enableEndPen:   false,
			expA:           "AT-",
			expB:           "--G",
			expectedScore:  0,
		},
		// полное совпадение. ничего не делаем.
		{
			a:              "AAAA",
			b:              "AAAA",
			enableStartPen: true,
			enableEndPen:   true,
			expA:           "AAAA",
			expB:           "AAAA",
			expectedScore:  20,
		},
		{
			a:              "ATGCCC",
			b:              "ATTTCCCC",
			enableStartPen: true,
			enableEndPen:   true,
			expA:           "A--TGCCC",
			expB:           "ATTTCCCC",
			expectedScore:  10,
		},
		// проверено с помощью https://www.ebi.ac.uk/Tools/psa/emboss_needle/
		{
			a: "GCGCGTGCGCGGAAGGAGCCAAGGTGAAGTTGTAGCAGTGTGTCAGAAGAGGTGCGTGGC" +
				"ACCATGCTGTCCCCCGAGGCGGAGCGGGTGCTGCGGTACCTGGTCGAAGTAGAGGAGTTG",
			b: "GACTTGTGGAACCTACTTCCTGAAAATAACCTTCTGTCCTCCGAGCTCTCCGCACCCGTG" +
				"GATGACCTGCTCCCGTACACAGATGTTGCCACCTGGCTGGATGAATGTCCGAATGAAGCG",
			enableStartPen: true,
			enableEndPen:   true,
			expA: "G-CGCGTGCGCGGAAGGAGCCAA---GGTGAAGTTGTAGCAGTGTGTCAGAAGAGGTGCGTGGC" +
				"ACCATGCTGTCC---CCCGAGGCGGAGCGGGTGCTGCGGTAC------------------CTGGTCGAA-GT---AGAGGAGTTG",
			expB: "GACTTGT--------GGAACCTACTTCCTGAA--AATAACCTTCTGTC---------------CTCCGAGCTCTCCGCACCCGTG" +
				"GATGACC---TGCTCCCGTACACAGATGTTGCCACCTGGCTGGATGAATGTCCGA-ATGAAGCG",
			expectedScore: 46,
		},
		// проверено с помощью https://www.ebi.ac.uk/Tools/psa/emboss_needle/
		// тут нет штрафов за крайние гэпы.
		{
			a: "GCGCGTGCGCGGAAGGAGCCAAGGTGAAGTTGTAGCAGTGTGTCAGAAGAGGTGCGTGGC" +
				"ACCATGCTGTCCCCCGAGGCGGAGCGGGTGCTGCGGTACCTGGTCGAAGTAGAGGAGTTG",
			b: "GACTTGTGGAACCTACTTCCTGAAAATAACCTTCTGTCCTCCGAGCTCTCCGCACCCGTG" +
				"GATGACCTGCTCCCGTACACAGATGTTGCCACCTGGCTGGATGAATGTCCGAATGAAGCG",
			enableStartPen: false,
			enableEndPen:   false,
			expA: "GCGCGTGCGCGGAAGGAGCCAAGGTGAAGTTGTAGCAGTGTGTCAGAAGAGGTGCGTGGC" +
				"AC-----------------CATGCTGTCCCCCGAG----GCGGAGCGGGTGCTGCGGTACCTGGT--CGAAGTAGAGGAGTTG--------------------------------",
			expB: "------------------------------------------------GACTT--GTGGAACCTACTTCCTGAAAATAACCTTCTGTCCTCCGAGCTCTCCGCACCCGTG" +
				"GATG----ACCTGCTCCCGTA-CACAGATGTTGCCACCTGGCTGGATGAATGTCCGAATGAAGCG",
			expectedScore: 70,
		},
	} {
		s.aligner.gapStartPenalty = c.enableStartPen
		s.aligner.gapEndPenalty = c.enableEndPen

		a, b, score := s.aligner.Align(c.a, c.b)
		s.Equal(c.expA, a)
		s.Equal(c.expB, b)
		s.Equal(c.expectedScore, score)
	}
}

func TestSequenceAlignerExtendSuite(t *testing.T) {
	suite.Run(t, new(SequenceAlignerExtendTestSuite))
}
