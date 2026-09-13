package common

type ProgressMessage struct {
	TaskID     string  `json:"task_id"`
	Percentage float64 `json:"percentage"`
}

var Broadcast = make(chan ProgressMessage)
