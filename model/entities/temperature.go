package entities

import (
	"fmt"
	"regexp"
	"strconv"
)

var tempRegex regexp.Regexp

func init() {
	tempRegex = *regexp.MustCompile(`\d\d,([0-9]*[.]?[0-9]+);$`)
}

type TemperatureEvent struct {
	Event
	Temperature float64 `json:"tmp"`
}

func (e *TemperatureEvent) ParseBinary(b []byte) (err error) {
	if err = e.Event.ParseBinary(b); err != nil {
		return err
	}
	sub := tempRegex.FindStringSubmatch(string(b))
	if len(sub) < 2 {
		return fmt.Errorf("wrong format")
	}
	tempStr := sub[1]
	e.Temperature, err = strconv.ParseFloat(tempStr, 64)
	return
}
