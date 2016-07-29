package main

import (
	"encoding/json"
	"net/http"
	"time"
)

// LocSpec represents one piece of payload
type LocSpec struct {
	UtcOffset int      `json:"utc_offset"`
	LocalTime string   `json:"local_time"`
	Delta     int      `json:"delta"`
	Locations []string `json:"locations"`
}

// Response represents the main response
type Response struct {
	Status  int       `json:"status"`
	Message string    `json:"message"`
	Payload []LocSpec `json:"payload"`
}

// AppendPayloadFor takes two ints, hours and minutes, and appends one
// LocSpec to the response representing the calculated data for the
// requested parameters
func (r *Response) AppendPayloadFor(h, m int) {
	utcOffset, delta := om.CalculateOffsets(h, m)
	localtime := om.GetLocaltime(utcOffset).Format(time.RFC1123Z)
	cities := om.GetCities(utcOffset)
	r.Payload = append(r.Payload, LocSpec{utcOffset, localtime, delta, cities})
}

// RespondJSON renders the response as json and writes it to w
func (r *Response) RespondJSON(w http.ResponseWriter) {
	w.WriteHeader(r.Status)
	w.Header().Add("Content-Type", "application/json; charset=utf-8")
	w.Header().Add("Access-Control-Allow-Origin", "*")
	data, _ := json.Marshal(r)
	w.Write(data)
}

// SetNotFound prepares the response to be rendered with status 404
// and message set to the provided error's message
func (r *Response) SetNotFound(err error) {
	r.Status = http.StatusNotFound
	r.Message = err.Error()
}
