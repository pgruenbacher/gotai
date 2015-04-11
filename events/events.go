package events

import (
	"github.com/pgruenbacher/gotai/utils"
)

type Event struct {
	Id string
}

func NewEvent() Event {
	return Event{
		Id: utils.RandSeq(6),
	}
}
