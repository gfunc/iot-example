package entities

import (
	"fmt"
	"regexp"
	"time"
)

const timeFormat = `2006-01-02 15:04:05`

var eventRegex, tempRegex, qualityRegex regexp.Regexp

func init() {
	eventRegex = *regexp.MustCompile(`^(.+),(\d\d\d\d-\d\d-\d\d \d\d:\d\d:\d\d),`)
	tempRegex = *regexp.MustCompile(`\d\d,([0-9]*[.]?[0-9]+);$`)
	qualityRegex = *regexp.MustCompile(`[a-zA-Z]+:[0-9]*[.]?[0-9]+`)
}

type Event struct {
	WaterMark time.Time `json:"water_mark"`
	DeviceID  string    `json:"device_id"`
	EventTime time.Time `json:"event_time"`
}

func (e *Event) GetDeviceID() string {
	return e.DeviceID
}

func (e *Event) ParseBinary(b []byte) (err error) {
	str := string(b)
	sub := eventRegex.FindStringSubmatch(str)
	if sub == nil || len(sub) < 2 {
		return fmt.Errorf("wrong format")
	}
	e.DeviceID = sub[1]
	e.EventTime, err = time.Parse(timeFormat, sub[2])
	e.WaterMark = time.Now()
	return
}
