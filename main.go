package main

import (
	"flag"
	"fmt"
	"net/http"
	"net/url"
	"os"
)

var (
	cfg = struct {
		Map  *string
		Host *string
		Port *int
	}{
		flag.String("map", "./offset_map.json", "Location of 'offset_map.json'"),
		flag.String("host", "localhost", "Hostname or IP-Address to bind to"),
		flag.Int("port", 1620, "Port number to listen on"),
	}
	om OffsetMap
)

func main() {
	var err error
	fmt.Print("[wo.istes.jetzt] loading offset map... ")
	flag.Parse()

	om, err = LoadOffsetMap(*cfg.Map)
	if err != nil {
		os.Exit(1)
	}
	fmt.Println("ok.")

	http.HandleFunc("/", dannHandler)
	fmt.Printf("[wo.istes.jetzt] listening on 'http://%s:%d'... ", *cfg.Host, *cfg.Port)
	http.ListenAndServe(fmt.Sprintf("%s:%d", *cfg.Host, *cfg.Port), nil)
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

func parseURLForTime(url *url.URL) (h, m int, err error) {
	var d int
	_, err = fmt.Sscanf(url.Path, "/%4d", &d)
	return d / 100 % 24, d % 100 % 60, err
}
