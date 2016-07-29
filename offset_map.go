package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math"
	"os"
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

func (om *OffsetMap) LoadFromFile(filename string) (err error) {
	if _, err = os.Stat(filename); os.IsNotExist(err) {
		return fmt.Errorf("%s does not exist", filename)
	}

	var input []byte
	input, err = ioutil.ReadFile(filename)
	if err != nil {
		return fmt.Errorf("could not read from %s", filename)
	}

	err = json.Unmarshal(input, om)
	if err != nil {
		return fmt.Errorf("could not parse contents of %s", filename)
	}

	return
}

func (om OffsetMap) GetCities(offset int) (cities []string) {
	cities, _ = om.data[offset]
	return
}

func (om OffsetMap) GetLocaltime(offset int) (t time.Time) {
	name := strconv.Itoa(offset)
	t = time.Now().In(time.FixedZone(name, offset))
	return
}

func (om OffsetMap) CalculateUtcOffset(h, m int) (o int) {
	now := time.Now().UTC()
	then := time.Date(now.Year(), now.Month(), now.Day(), h, m, 0, 0, now.Location())

	for {
		delta := then.Sub(now).Seconds()
		o = int(delta)

		if om.keys.hasWithinBounds(o) {
			break
		}

		sgn := math.Copysign(1, delta)
		then = then.Add(-24 * time.Hour * time.Duration(sgn))
	}

	return
}

func (om OffsetMap) CalculatePreviousUtcOffset(h, m int) (o int) {
	offset := om.CalculateUtcOffset(h, m)
	o = om.keys[sort.SearchInts(om.keys, offset)-1]
	return
}
