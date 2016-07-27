package main

import (
	"encoding/json"
	"net/http"
	"time"
)

type LocSpec struct {
	UtcOffset int      `json:"utc_offset"`
	LocalTime string   `json:"local_time"`
	Locations []string `json:"locations"`
}

type Response struct {
	Status  int       `json:"status"`
	Message string    `json:"message"`
	Payload []LocSpec `json:"payload"`
}

func (r *Response) appendPayloadFor(h, m int) {
	offset := om.CalculatePreviousUtcOffset(h, m)

	r.Payload = append(r.Payload, LocSpec{
		offset,
		om.GetLocaltime(offset).Format(time.RFC1123Z),
		om.GetCities(offset),
	})
}

func (r *Response) RespondJSON(w http.ResponseWriter) {
	w.WriteHeader(r.Status)
	w.Header().Add("Content-Type", "application/json; charset=utf-8")
	data, _ := json.Marshal(r)
	w.Write(data)
}

func (r *Response) SetNotFound(err error) {
	r.Status = http.StatusNotFound
	r.Message = err.Error()
}
