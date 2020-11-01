package main

import (
	"testing"

	"github.com/stretchr/testify/suite"
)

type SequenceAlignerMemTestSuite struct {
	suite.Suite
	aligner *SequenceAlignerMem
}

func (s *SequenceAlignerMemTestSuite) SetupTest() {
	cfg := &SequenceAlignerConfig{
		GapPenalty: -10,
	}
	s.aligner = NewSequenceAlignerMem(cfg, NewDNAAdapter())
}

func (s *SequenceAlignerMemTestSuite) TestAlign() {
	for _, c := range []struct {
		a              string
		b              string
		enableStartPen bool
		enableEndPen   bool
		expA           string
		expB           string
		expectedScore  int
	}{
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
		// вхождение как подстроки. нужно только выровнять длину гэпами.
		{
			a:              "AAAA",
			b:              "AA",
			enableStartPen: true,
			enableEndPen:   true,
			expA:           "AAAA",
			expB:           "A--A",
			expectedScore:  -10,
		},
		// отличие в одном символе. в данной конфигурации 2 гэпа дороже, поэтому так и осталяем
		{
			a:              "AAAT",
			b:              "AAAA",
			enableStartPen: true,
			enableEndPen:   true,
			expA:           "AAAT",
			expB:           "AAAA",
			expectedScore:  11,
		},
		// можно совместить только один символ, куда вставим гэп неважно
		{
			a:              "AAT",
			b:              "AC",
			enableStartPen: true,
			enableEndPen:   true,
			expA:           "AAT",
			expB:           "A-C",
			expectedScore:  -9,
		},
		// дефолтный тест из задания
		{
			a:              "AATCG",
			b:              "AACG",
			enableStartPen: true,
			enableEndPen:   true,
			expA:           "AATCG",
			expB:           "AA-CG",
			expectedScore:  10,
		},
		// штраф за gap с двух сторон.
		// максимум, который тут можно получить — совместить 3 символа
		{
			a:              "AATTTTTTAATCGGGGGGGG",
			b:              "AACC",
			enableStartPen: true,
			enableEndPen:   true,
			expA:           "AATTTTTTAATCGGGGGGGG",
			expB:           "--------AA-C-------C",
			expectedScore:  -149,
		},
		// прочерки слева бесплатные.
		// и хотя на первый взгляд кажется, что ничего изменится не должно,
		// на самом деле в --------AA-C-------C гэпы после C слишком дорогие.
		{
			a:              "AATTTTTTAATCGGGGGGGG",
			b:              "AACC",
			expA:           "AATTTTTTAATCGGGGGGGG",
			expB:           "----------------AACC",
			enableStartPen: false,
			enableEndPen:   true,
			expectedScore:  -16,
		},
		// зеркальный предыдущему случай, только тут ещё и 2 совпадения
		{
			a:              "AATTTTTTAATCGGGGGGGG",
			b:              "AACC",
			expA:           "AATTTTTTAATCGGGGGGGG",
			expB:           "AACC----------------",
			enableStartPen: true,
			enableEndPen:   false,
			expectedScore:  2,
		},
		// так как gap по краям бесплатные,
		// то выгоднее всего подвинуть вторую строку в нужное место
		{
			a:              "TTTTTTAATCGGGGGGGG",
			b:              "AACC",
			expA:           "TTTTTTAATCGGGGGGGG",
			expB:           "------AACC--------",
			enableStartPen: false,
			enableEndPen:   false,
			expectedScore:  11,
		},
		// большие последовательности проверены с помощью
		// https://www.ebi.ac.uk/Tools/psa/emboss_needle/
		{
			a: "GCGCGTGCGCGGAAGGAGCCAAGGTGAAGTTGTAGCAGTGTGTCAGAAGAGGTGCGTGGC" +
				"ACCATGCTGTCCCCCGAGGCGGAGCGGGTGCTGCGGTACCTGGTCGAAGTAGAGGAGTTG",
			b: "GACTTGTGGAACCTACTTCCTGAAAATAACCTTCTGTCCTCCGAGCTCTCCGCACCCGTG" +
				"GATGACCTGCTCCCGTACACAGATGTTGCCACCTGGCTGGATGAATGTCCGAATGAAGCG",
			enableStartPen: true,
			enableEndPen:   true,
			expA: "GCGCGTGCGCGGAAGGAGCCAAGGTGAAGTTGTAGCAGTGTGTCAGAAGAGGTGCGTGGC" +
				"ACCA-TGC-TGTCCCCCGAGGCGGAG-CGGGTGCTGCGG--TACCTGGTCGAA-GTA-GAG-GAGTTG",
			expB: "G-AC-TT-GTGGAAC-CTACTTCCTGAAA--ATAACCTTCTGTCCTCCGAGCT" +
				"-CTCCGCACCCGTGGATGACCTGC-TCCCGTACACAGATGTTGCCACCTGGCTGGATGAATGTCCGAATGAAGCG",
			expectedScore: -41,
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

func TestSequenceAlignerMemSuite(t *testing.T) {
	suite.Run(t, new(SequenceAlignerMemTestSuite))
}
