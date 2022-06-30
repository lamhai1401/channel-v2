package client

// Message message
type Message struct {
	Topic   string      `mapstructure:"topic" json:"topic"`
	Event   string      `mapstructure:"event" json:"event"`
	Ref     string      `mapstructure:"ref" json:"ref"`
	Payload interface{} `mapstructure:"payload" json:"payload"`
}
