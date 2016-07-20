package main

import (
	"encoding/json"
	"io/ioutil"
	"math"
	"sort"
	"strconv"
	"time"
)

type offsetMapData map[int][]string
type offsetMapKeys []int

func (omk offsetMapKeys) hasWithinBounds(x int) bool {
	return omk[0] <= x && x <= omk[len(omk)-1]
}

type OffsetMap struct {
	data offsetMapData
	keys offsetMapKeys
}

func (om *OffsetMap) UnmarshalJSON(data []byte) (err error) {
	input := make(map[string][]string, 0)
	err = json.Unmarshal(data, &input)

	if err != nil {
		return err
	}

	om.data = make(offsetMapData)
	om.keys = make(offsetMapKeys, 0)

	for offsetString, cities := range input {
		offset, err := strconv.Atoi(offsetString)

		if err != nil {
			return err
		}

		om.data[offset] = cities
		om.keys = append(om.keys, offset)
	}

	sort.Ints(om.keys)

	return err
}

func LoadOffsetMap(filename string) (om OffsetMap, err error) {
	input, err := ioutil.ReadFile(filename)
	if err == nil {
		err = json.Unmarshal(input, &om)
	}

	return
}

func (om OffsetMap) GetOffsets() offsetMapKeys {
	return om.keys
}

func (om OffsetMap) GetCities(offset int) (cities []string) {
	cities, _ = om.data[offset]
	return
}

func (om OffsetMap) GetLocaltimeIn(offset int) time.Time {
	return time.Now().In(time.FixedZone(strconv.Itoa(offset), offset))
}

func (om OffsetMap) GetPreviousOffset(offset int) int {
	pos := sort.SearchInts(om.keys, offset) - 1
	return om.keys[pos]
}

func (om OffsetMap) GetOffsetFor(h, m int) int {
	now := time.Now().UTC()
	then := time.Date(now.Year(), now.Month(), now.Day(), h, m, 0, 0, now.Location())
	delta := then.Sub(now).Seconds()

	if !om.keys.hasWithinBounds(int(delta)) {
		sgn := math.Copysign(1, delta)
		then = then.Add(-24 * time.Hour * time.Duration(sgn))
		delta = then.Sub(now).Seconds()
	}

	return int(delta)
}

func (om OffsetMap) GetOffsetForUpcoming(h, m int) int {
	return om.GetPreviousOffset(om.GetOffsetFor(h, m))
}
