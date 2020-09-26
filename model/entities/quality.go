package entities

import (
	"strconv"
	"strings"
)

type QualityEvent struct {
	Event
	Indexes map[string]float64 `json:"indexes"`
}

func (e *QualityEvent) ParseBinary(b []byte) (err error) {
	if err = e.Event.ParseBinary(b); err != nil {
		return err
	}
	sub := qualityRegex.FindAllStringSubmatch(string(b), -1)
	e.Indexes = make(map[string]float64, 0)
	for i := 0; i < len(sub); i++ {
		kv := strings.Split(sub[i][0], ":")
		v, err := strconv.ParseFloat(kv[1], 64)
		if err != nil {
			return err
		}
		e.Indexes[kv[0]] = v
	}
	return
}
