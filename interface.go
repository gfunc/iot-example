package iot_practise

import "encoding/json"

type EventEntity interface {
	GetDeviceID() string
	ParseBinary([]byte) error
}

type EventAlert struct {
	Event EventEntity
	Msg   string
}

func (ea *EventAlert) Error() string {
	if ea.Event == nil {
		return ea.Msg
	}
	b, err := json.Marshal(ea)
	if err != nil {
		return "error marshal json for event alert" + err.Error()
	}
	return string(b) + ea.Msg
}

type EventInspector interface {
	Inspect(entity EventEntity, errChan chan<- error)
}
