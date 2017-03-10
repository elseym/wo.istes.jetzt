package responders

import (
	"testing"
	"github.com/elseym/go-tzlib"
	"time"
)

func expectTimeMatches(actual, expected time.Time, t *testing.T) {
	if (actual.String() != expected.String()) {
		t.Errorf("Time %s does not match expected %s", actual.String(), expected.String());
	}
}

func TestParseTime(t *testing.T) {
	// 6 digits
	expectTimeMatches(ParseTime("062020"), tzlib.Time(6, 20, 20), t);
	expectTimeMatches(ParseTime("161010"), tzlib.Time(16, 10, 10), t);

	// 5 digits
	expectTimeMatches(ParseTime("42020"), tzlib.Time(4, 20, 20), t);
	expectTimeMatches(ParseTime("81010"), tzlib.Time(8, 10, 10), t);

	// 4 digits
	expectTimeMatches(ParseTime("1620"), tzlib.Time(16, 20, 0), t);
	expectTimeMatches(ParseTime("0620"), tzlib.Time(6, 20, 0), t);

	// 3 digts
	expectTimeMatches(ParseTime("620"), tzlib.Time(6, 20, 0), t);
	expectTimeMatches(ParseTime("920"), tzlib.Time(9, 20, 0), t);

	// 2 digts
	expectTimeMatches(ParseTime("06"), tzlib.Time(6, 0, 0), t);
	expectTimeMatches(ParseTime("16"), tzlib.Time(16, 0, 0), t);

	// 1 digit
	expectTimeMatches(ParseTime("6"), tzlib.Time(6, 0, 0), t);
	expectTimeMatches(ParseTime("9"), tzlib.Time(9, 0, 0), t);
}
