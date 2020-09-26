package main

import "os"

func readNFromFile(filename string, n int) ([]*Sequence, error) {
	f, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	parser := NewFastaParser(f)
	seqs := make([]*Sequence, n)
	for i := 0; i < n; i++ {
		s, err := parser.Next()
		if err != nil {
			return nil, err
		}
		seqs[i] = s
	}

	return seqs, nil
}

func loadFromFile(filename string) ([]*Sequence, error) {
	return readNFromFile(filename, 2)
}

func loadFromFiles(filename1, filename2 string) ([]*Sequence, error) {
	s1, err := readNFromFile(filename1, 1)
	if err != nil {
		return nil, err
	}

	s2, err := readNFromFile(filename2, 1)
	if err != nil {
		return nil, err
	}

	return append(s1, s2...), nil
}

func loadSequences(fileNames []string) ([]*Sequence, error) {
	if len(fileNames) == 1 {
		return loadFromFile(fileNames[0])
	}
	if len(fileNames) == 2 {
		return loadFromFiles(fileNames[0], fileNames[1])
	}
	return nil, ErrWrongNumberOfFiles
}
