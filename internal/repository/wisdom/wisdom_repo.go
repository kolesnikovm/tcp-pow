package wisdomrepo

import (
	"encoding/csv"
	"fmt"
	"os"
)

type WisdomRepo struct {
	quotes [][]string
}

func NewWisdomRepo(filePath string) (*WisdomRepo, error) {
	const op = "wisdomrepo.NewWisdomRepo"

	f, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("%s: failed to read input file '%s': %w", op, filePath, err)
	}
	defer f.Close()

	csvReader := csv.NewReader(f)

	records, err := csvReader.ReadAll()
	if err != nil {
		return nil, fmt.Errorf("%s: failed to parse file '%s' as CSV: %w", op, filePath, err)
	}

	return &WisdomRepo{
		quotes: records,
	}, nil
}
