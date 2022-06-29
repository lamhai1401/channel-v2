package client

// Message message
type Message struct {
	Topic   string      `mapstructure:"topic"`
	Event   string      `mapstructure:"event"`
	Ref     string      `mapstructure:"ref"`
	Payload interface{} `mapstructure:"payload"`
}
