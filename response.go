package main

import (
	"encoding/json"
	"net/http"
	"time"
)

type LocSpec struct {
	UtcOffset int      `json:"utc_offset"`
	LocalTime string   `json:"local_time"`
	Delta     int      `json:"delta"`
	Locations []string `json:"locations"`
}

type Response struct {
	Status  int       `json:"status"`
	Message string    `json:"message"`
	Payload []LocSpec `json:"payload"`
}

func (r *Response) AppendPayloadFor(h, m int) {
	utcOffset, delta := om.CalculateOffsets(h, m)
	localtime := om.GetLocaltime(utcOffset).Format(time.RFC1123Z)
	cities := om.GetCities(utcOffset)

	r.Payload = append(r.Payload, LocSpec{utcOffset, localtime, delta, cities})
}

func (r *Response) RespondJSON(w http.ResponseWriter) {
	w.WriteHeader(r.Status)
	w.Header().Add("Content-Type", "application/json; charset=utf-8")
	w.Header().Add("Access-Control-Allow-Origin", "*")
	data, _ := json.Marshal(r)
	w.Write(data)
}

func (r *Response) SetNotFound(err error) {
	r.Status = http.StatusNotFound
	r.Message = err.Error()
}
