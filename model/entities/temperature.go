package entities

import "strconv"

type TemperatureEvent struct {
	Event
	Temperature float64 `json:"tmp"`
}

func (e *TemperatureEvent) ParseBinary(b []byte) (err error) {
	if err = e.Event.ParseBinary(b); err != nil {
		return err
	}
	sub := tempRegex.FindStringSubmatch(string(b))
	tempStr := sub[1]
	e.Temperature, err = strconv.ParseFloat(tempStr, 64)
	return
}
