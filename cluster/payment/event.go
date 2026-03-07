package payment

type EventType int

const (
	EventTypeNewConribution EventType = iota
)

type Event struct {
	Type EventType `json:"type"`
	// TODO: add contrib here
}
