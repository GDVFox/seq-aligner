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
	s.aligner = NewSequenceAlignerMem(-10, NewDNAAdapter())
}

func (s *SequenceAlignerMemTestSuite) TestAlign() {
	for _, c := range []struct {
		a             string
		b             string
		expA          string
		expB          string
		expectedScore int
	}{
		// полное совпадение. ничего не делаем.
		{
			a:             "AAAA",
			b:             "AAAA",
			expA:          "AAAA",
			expB:          "AAAA",
			expectedScore: 20,
		},
		// вхождение как подстроки. нужно только выровнять длину гэпами.
		{
			a:             "AAAA",
			b:             "AA",
			expA:          "AAAA",
			expB:          "AA--",
			expectedScore: -10,
		},
		// отличие в одном символе. в данной конфигурации 2 гэпа дороже, поэтому так и осталяем
		{
			a:             "AAAT",
			b:             "AAAA",
			expA:          "AAAT",
			expB:          "AAAA",
			expectedScore: 11,
		},
		// можно совместить только один символ, куда вставим гэп неважно
		{
			a:             "AAT",
			b:             "AC",
			expA:          "AAT",
			expB:          "AC-",
			expectedScore: -9,
		},
		// дефолтный тест из задания
		{
			a:             "AATCG",
			b:             "AACG",
			expA:          "AATCG",
			expB:          "AA-CG",
			expectedScore: 10,
		},
		// штраф за gap с двух сторон.
		// максимум, который тут можно получить — совместить 3 символа
		{
			a:             "AATTTTTTAATCGGGGGGGG",
			b:             "AACC",
			expA:          "AATTTTTTAATCGGGGGGGG",
			expB:          "AAC--------C--------",
			expectedScore: -149,
		},
		// большие последовательности проверены с помощью
		// https://www.ebi.ac.uk/Tools/psa/emboss_needle/
		{
			a: "GCGCGTGCGCGGAAGGAGCCAAGGTGAAGTTGTAGCAGTGTGTCAGAAGAGGTGCGTGGC" +
				"ACCATGCTGTCCCCCGAGGCGGAGCGGGTGCTGCGGTACCTGGTCGAAGTAGAGGAGTTG",
			b: "GACTTGTGGAACCTACTTCCTGAAAATAACCTTCTGTCCTCCGAGCTCTCCGCACCCGTG" +
				"GATGACCTGCTCCCGTACACAGATGTTGCCACCTGGCTGGATGAATGTCCGAATGAAGCG",
			expA: "GCGCGTGCGCGGAAGGAGCCAAGGTGAAGTTGTAGCAGTGTGTCAGAAGAGGTGCGTGG" +
				"CACCA-TG-CTGTCCCCCGAGGCGGA-GCGGGTGCTGC--GGTACCTGGTCGAA-GTA-GAG-GAGTTG",
			expB: "GA-CTTGTG-GAACCTA-CTTCC-TGAAAATA-ACCT-TCTGTCCTCCGAGCT-CTCCGCACCCGTGGA" +
				"TGACCTGCTCC-CGTACACAGATGTTGCCACCTGGCTGGATGAATGTCCGAATGAAGCG",
			expectedScore: -41,
		},
	} {
		a, b, score := s.aligner.Align(c.a, c.b)
		s.Equal(c.expA, a)
		s.Equal(c.expB, b)
		s.Equal(c.expectedScore, score)
	}
}

func TestSequenceAlignerMemSuite(t *testing.T) {
	suite.Run(t, new(SequenceAlignerMemTestSuite))
}
