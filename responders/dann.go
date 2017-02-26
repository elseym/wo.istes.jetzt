package responders

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	tzlib "github.com/elseym/go-tzlib"
)

type dann struct {
	tzlib   *tzlib.Tzlib
	Status  int        `json:"status"`
	Message string     `json:"message"`
	Payload []dannzone `json:"payload"`
}

type dannzone struct {
	Requested string         `json:"requested"`
	Message   string         `json:"message"`
	Delta     int            `json:"delta"`
	Localtime string         `json:"localtime"`
	Timezone  tzlib.Timezone `json:"timezone"`
}

// DannResponder returns a new and initialised DannResponder
func DannResponder(l *tzlib.Tzlib) *dann {
	d := &dann{tzlib: l}

	return d.reset()
}

// ServeHTTP handles it
func (d dann) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	frg := strings.Split(r.URL.Path, "/")

	if len(frg) != 1 {
		d.respond404(w)
		return
	}

	ts := parseURL(frg[len(frg)-1])

	if r.URL.Query().Get("mode") == "12h" {
		ts = twelveify(ts)
	}

	for _, t := range ts {
		d.append(t)
	}

	d.respond(w)
}

// respond writes headers and body to response
func (d dann) respond(w http.ResponseWriter) {
	w.WriteHeader(d.Status)
	w.Header().Add("Content-Type", "application/json; charset=utf-8")
	w.Header().Add("Access-Control-Allow-Origin", "*")

	data, _ := json.Marshal(d)
	d.reset()
	w.Write(data)
}

// respond404 sets 404 before responding
func (d dann) respond404(w http.ResponseWriter) {
	d.reset()
	d.Status = 404
	d.Message = "Not Found"
	d.respond(w)
}

// append searches for the requested Time and adds
// timezone information to the response
func (d *dann) append(t time.Time) *dann {
	z, err := d.tzlib.WhereWillItBe(t)
	dz := dannzone{Requested: t.Format("15:04:05"), Message: "OK"}

	if err == nil {
		dz.Delta = int(z.Until(t) / time.Second)
		dz.Localtime = z.Localtime().Format(time.RFC1123Z)
		dz.Timezone = z
	} else {
		dz.Message = err.Error()
	}

	d.Payload = append(d.Payload, dz)

	return d
}

// reset removes the payload and resets status to 200 OK
func (d *dann) reset() *dann {
	d.Payload = make([]dannzone, 0, 2)
	d.Status = 200
	d.Message = "OK"

	return d
}

// twelveify adds after each Time its twelve hour future complement
func twelveify(ts []time.Time) (rts []time.Time) {
	for i := 0; i < len(ts); i++ {
		rts = append(rts, ts[i], ts[i].Add(12*time.Hour))
	}

	return
}

// parseURL transforms "/some/prefix/2342,1620,0" to
// Times 23:42, 16:20, 0:00 with current date
func parseURL(u string) (ts []time.Time) {
	vs := strings.Split(strings.TrimPrefix(u, "/"), ",")

	for _, v := range vs {
		ts = append(ts, parseTime(v))
	}

	return
}

// parseTime parses strings to time in a biased manner:
// > 4 digits: (...)(h)hmmss
// > 2 digits: (h)hmm
// < 3 digits: (h)h
func parseTime(u string) time.Time {
	h, m, s := 0, 0, 0

	fmt.Sscan(u, &h)

	// set h to rightmost six digits
	h %= 1e6

	// h has more than four digits:
	// assume (h)hmmss, extract seconds
	if h > 9999 {
		s = h % 100
		h /= 100
	}

	// h has more than two digits:
	// assume (h)hmm, extract minutes
	if h > 99 {
		m = h % 100
		h /= 100
	}

	// h has less than three digits:
	// assume (h)h, use as hours
	return tzlib.Time(h, m, s)
}
