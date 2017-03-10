package responders

import (
	"time"
	"github.com/elseym/go-tzlib"
	"strconv"
)

// ParseTime parses strings to time in a biased manner:
// 6 digits hhmmss
// 5 digits hmmss
// 4 digits hhmm
// 3 digits hmm
// 2 digits hh
// 1 digit h
func ParseTime(rawInput string) time.Time {
	h, m, s := 0, 0, 0

	// from is inclusive, to exclusive
	extractInt := func(from, to int) int {
		i, err := strconv.Atoi(rawInput[from:to])
		if err != nil {
			return 0
		}
		return i
	}

	switch len(rawInput) {
	case 6:
		h = extractInt(0, 2)
		m = extractInt(2, 4)
		s = extractInt(4, 6)
	case 5:
		h = extractInt(0, 1)
		m = extractInt(1, 3)
		s = extractInt(3, 5)
	case 4:
		h = extractInt(0, 2)
		m = extractInt(2, 4)
	case 3:
		h = extractInt(0, 1)
		m = extractInt(1, 3)
	case 2:
		h = extractInt(0, 2)
	case 1:
		h = extractInt(0, 1)
	}

	return tzlib.Time(h, m, s)
}
