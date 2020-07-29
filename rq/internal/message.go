package internal

type TimedMessage struct {
	Time    int64  `json:"time"`
	Payload string `json:"payload"`
}
