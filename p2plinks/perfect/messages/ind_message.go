package messages

type IndMessage struct {
	from string
	message string
}

func (msg *IndMessage) Init(from string, message string) *IndMessage {
	msg.from = from
	msg.message = message

	return msg
}

func NewIndMessage(from string, message string) *IndMessage {
	return new(IndMessage).Init(from, message)
}

func (msg *IndMessage) From() string {
	return msg.from
}

func (msg *IndMessage) Message() string {
	return msg.message
}
