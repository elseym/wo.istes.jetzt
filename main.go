package main

import (
	"errors"
	"flag"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"regexp"
	"strconv"
)

type config struct {
	Map  *string
	Host *string
	Port *int
}

var rx *regexp.Regexp
var om OffsetMap
var omFile string

func main() {
	rx = regexp.MustCompile(`/(\d\d?)(\d\d)$`)

	cfg := config{
		flag.String("map", "./offset_map.json", "Location of 'offset_map.json'"),
		flag.String("host", "localhost", "Hostname or IP-Address to bind to"),
		flag.Int("port", 1620, "Port number to listen on"),
	}
	flag.Parse()

	var err error
	fmt.Print("[wo.istes.jetzt] loading offset map... ")
	om, err = LoadOffsetMap(*cfg.Map)
	if err != nil {
		os.Exit(1)
	}
	fmt.Println("ok.")

	addr := fmt.Sprintf("%s:%d", *cfg.Host, *cfg.Port)

	http.HandleFunc("/", dannHandler)
	fmt.Printf("[wo.istes.jetzt] listening on 'http://%s:%d'... ", *cfg.Host, *cfg.Port)
	http.ListenAndServe(addr, nil)
	fmt.Println("done.")
}

func dannHandler(w http.ResponseWriter, r *http.Request) {
	response := Response{http.StatusOK, "OK", make([]LocSpec, 0)}

	if h, m, err := parseURLForTime(r.URL); err == nil {
		response.appendPayloadFor(h, m)
		if r.URL.Query().Get("mode") == "12h" {
			response.appendPayloadFor((h+12)%24, m)
		}
	} else {
		response.SetNotFound(err)
	}

	response.RespondJSON(w)
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
