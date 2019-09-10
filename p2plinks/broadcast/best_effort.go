package broadcast

import (
	"../perfect"
	perfectMessages "../perfect/messages"
	"./messages"
)

type BestEffortBroadcast struct {
	indChannel chan messages.IndMessage
	reqChannel chan messages.ReqMessage
	link perfect.Link
}

func (beb *BestEffortBroadcast) Init() *BestEffortBroadcast {
	beb.indChannel = make(chan messages.IndMessage)
	beb.reqChannel = make(chan messages.ReqMessage)
	beb.link = *perfect.NewPerfectLink()

	return beb
}

func NewBestEffortBroadcast() *BestEffortBroadcast {
	return new(BestEffortBroadcast).Init()
}

func (beb *BestEffortBroadcast) GetIndChannel() chan messages.IndMessage {
	return beb.indChannel
}

func (beb *BestEffortBroadcast) GetReqChannel() chan messages.ReqMessage {
	return beb.reqChannel
}

func (beb *BestEffortBroadcast) PushIndMessageToChannel(message messages.IndMessage) {
	beb.indChannel <- message
}

func (beb *BestEffortBroadcast) PopIndMessageFromChannel() messages.IndMessage {
	msg := <- beb.indChannel
	return msg
}

func (beb *BestEffortBroadcast) PushReqMessageToChannel(message messages.ReqMessage) {
	beb.reqChannel <- message
}

func (beb *BestEffortBroadcast) Start(address string) {
	beb.link.Start(address)
	go beb.KeepSending()
}

func (beb *BestEffortBroadcast) KeepSending() {
	for {
		select {
		case reqMsg := <- beb.reqChannel:
			beb.Broadcast(reqMsg)
		case indMsg := <- beb.link.GetIndChannel():
			beb.Deliver(*mapToBroadcastIndMessage(indMsg))
		}
	}
}

func (beb *BestEffortBroadcast) Broadcast(message messages.ReqMessage) {
	for i := 0; i < len(message.To()); i++ {
		msg := *mapToPerfectReqMessage(message, i)
		beb.link.PushReqMessageToChannel(msg)
	}
}

func (beb *BestEffortBroadcast) Deliver(message messages.IndMessage) {
	beb.indChannel <- message
}

func mapToPerfectReqMessage(message messages.ReqMessage, index int) *perfectMessages.ReqMessage {
	return perfectMessages.NewReqMessage(message.To()[index], message.Message())
}

func mapToBroadcastIndMessage(message perfectMessages.IndMessage) *messages.IndMessage {
	return messages.NewIndMessage(message.From(), message.Message())
}
