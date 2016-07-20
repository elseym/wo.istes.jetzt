package main

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/url"
	"regexp"
	"strconv"
	"time"
)

var rx *regexp.Regexp

type LocSpec struct {
	LocalTime string   `json:"local_time"`
	Locations []string `json:"locations"`
	UtcOffset int      `json:"utc_offset"`
}

type Response struct {
	Status  int       `json:"status"`
	Message string    `json:"message"`
	Payload []LocSpec `json:"payload"`
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

func main() {
	om, _ := LoadOffsetMap("./offset_map.json")
	rx = regexp.MustCompile(`^/(\d\d?)(\d\d)$`)

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		response := Response{http.StatusOK, "OK", make([]LocSpec, 0)}

		if h, m, err := parseURLForTime(r.URL); err == nil {
			offset := om.GetOffsetForUpcoming(h, m)

			response.Payload = append(response.Payload, LocSpec{
				om.GetLocaltimeIn(offset).Format(time.RFC1123Z),
				om.GetCities(offset),
				offset,
			})
		} else {
			response.Status = http.StatusNotFound
			response.Message = err.Error()
		}

		w.WriteHeader(response.Status)
		jsonresponse, _ := json.Marshal(response)
		w.Write(jsonresponse)
	})

	http.ListenAndServe("127.0.0.1:1620", nil)
}
