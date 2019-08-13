package messages

type ReqMessage struct {
	to string
	message string
}

func (msg *ReqMessage) Init(to string, message string) *ReqMessage {
	msg.to = to
	msg.message = message

	return msg
}

func NewReqMessage(to string, message string) *ReqMessage {
	return new(ReqMessage).Init(to, message)
}

func (msg *ReqMessage) To() string {
	return msg.to
}

func (msg *ReqMessage) Message() string {
	return msg.message
}
