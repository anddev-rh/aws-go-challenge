package utils

type Status string

const (
	StatusCompleted   Status = "completed"
	StatusIncompleted Status = "incompleted"
	StatusReady       Status = "ready for shipping"
)
