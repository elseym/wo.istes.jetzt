package main

import (
	"errors"
	"net/http"
	"net/url"
	"regexp"
	"strconv"
)

var rx *regexp.Regexp
var om OffsetMap

func main() {
	om, _ = LoadOffsetMap("./offset_map.json")
	rx = regexp.MustCompile(`/(\d\d?)(\d\d)$`)

	http.HandleFunc("/", dannHandler)
	http.ListenAndServe("127.0.0.1:1620", nil)
}

func dannHandler(w http.ResponseWriter, r *http.Request) {
	response := Response{http.StatusOK, "OK", make([]LocSpec, 0)}

	if h, m, err := parseURLForTime(r.URL); err == nil {
		response.appendPayloadFor(h, m)
	} else {
		response.SetNotFound(err)
	}

	response.Respond(w)
}

func parseURLForTime(url *url.URL) (h int, m int, err error) {
	matches := rx.FindStringSubmatch(url.EscapedPath())

	if len(matches) < 2 {
		return 0, 0, errors.New("could not parse time from url '" + url.Path + "'")
	}

	h, err = strconv.Atoi(matches[1])
	m, err = strconv.Atoi(matches[2])

	return h % 24, m % 60, nil
}
