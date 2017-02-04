package action

type Action struct {
	Action    string
	Timestamp time.Time
	Msg       map[string]string
}
