package action

import (
	"time"
)

type Action struct {
	Action    string
	Timestamp time.Time
	Msg       map[string]string
}
