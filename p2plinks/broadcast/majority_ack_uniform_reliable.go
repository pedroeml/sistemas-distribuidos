package broadcast

import (
	"./messages"
)

type MajorityAckUniformReliableBroadcast struct {
	indChannel chan messages.IndMessage
	reqChannel chan messages.ReqMessage
	ipAddress string
	// delivered Set
	// pending Set
	// ack Set
	beb BestEffortBroadcast
}

func (urb *MajorityAckUniformReliableBroadcast) Init(address string) *MajorityAckUniformReliableBroadcast {
	urb.indChannel = make(chan messages.IndMessage)
	urb.reqChannel = make(chan messages.ReqMessage)
	urb.ipAddress = address
	// delivered = make(Set)
	// pending = make(Set)
	// ack = make(Set)
	urb.beb = *NewBestEffortBroadcast()

	return urb
}

func NewMajorityAckUniformReliableBroadcast(address string) *MajorityAckUniformReliableBroadcast {
	return new(MajorityAckUniformReliableBroadcast).Init(address)
}

func (urb *MajorityAckUniformReliableBroadcast) IpAddress() string {
	return urb.ipAddress
}

func (urb *MajorityAckUniformReliableBroadcast) PushIndMessageToChannel(message messages.IndMessage) {
	urb.indChannel <- message
}

func (urb *MajorityAckUniformReliableBroadcast) PopIndMessageFromChannel() messages.IndMessage {
	msg := <- urb.indChannel
	return msg
}

func (urb *MajorityAckUniformReliableBroadcast) PushReqMessageToChannel(message messages.ReqMessage) {
	urb.reqChannel <- message
}

func (urb *MajorityAckUniformReliableBroadcast) Start() {
	urb.beb.Start(urb.ipAddress)
	go urb.KeepSending()
	go urb.KeepDelivering()
}

func (urb *MajorityAckUniformReliableBroadcast) KeepSending() {
	for {
		select {
		case reqMsg := <- urb.reqChannel:
			urb.Broadcast(reqMsg)
		case indMsg := <- urb.beb.GetIndChannel():
			urb.bebDeliver(indMsg)
		}
	}
}

func (urb *MajorityAckUniformReliableBroadcast) Broadcast(message messages.ReqMessage) {
// pending = pending U { (urb.IpAddress(), message.Message()) }
	urb.beb.PushReqMessageToChannel(message)
}

func (urb *MajorityAckUniformReliableBroadcast) KeepDelivering() {
//	for {
//		for each (source, message) in pending {
//			if urb.canDeliver(message) ^ not(message E delivered) {
// 				delivered = delivered U { (message.From(), message.Message()) }
//				urb.indChannel <- messages.NewIndMessage(source, message)
// 			}
//		}
//	}
}

func (urb *MajorityAckUniformReliableBroadcast) bebDeliver(message messages.IndMessage) {
//	ack[message] = ack[message] U message
// 	if not ((message.From(), message.Message()) E pending)
//		pending = pending U { (message.From(), message.Message()) }
//		urb.beb.Broadcast(?????)
}

func (urb *MajorityAckUniformReliableBroadcast) canDeliver(message messages.IndMessage) bool {
//	return #(ack[message]) > N/2
	return true
}
