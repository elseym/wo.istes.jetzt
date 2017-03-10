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

	if len(rawInput) == 6 {
		h, _ = strconv.Atoi(rawInput[:2])
		m, _ = strconv.Atoi(rawInput[2:4])
		s, _ = strconv.Atoi(rawInput[4:])
	}
	if len(rawInput) == 5 {
		h, _ = strconv.Atoi(rawInput[:1])
		m, _ = strconv.Atoi(rawInput[1:3])
		s, _ = strconv.Atoi(rawInput[3:])
	}
	if len(rawInput) == 4 {
		h, _ = strconv.Atoi(rawInput[:2])
		m, _ = strconv.Atoi(rawInput[2:])
	}
	if len(rawInput) == 3 {
		h, _ = strconv.Atoi(rawInput[:1])
		m, _ = strconv.Atoi(rawInput[1:])
	}
	if len(rawInput) < 3 {
		h, _ = strconv.Atoi(rawInput)
	}

	return tzlib.Time(h, m, s)
}
