package broadcast

import (
	"../utils"
	"./messages"
	"strings"
)

type MajorityAckUniformReliableBroadcast struct {
	indChannel        chan messages.IndMessage
	reqChannel        chan messages.ReqMessage
	ipAddress         string
	delivered         map[string] bool
	pending           map[string] bool
	ack               map[string] int
	targetIpAddresses []string
	numberOfProcesses int
	beb               BestEffortBroadcast
}

func (urb *MajorityAckUniformReliableBroadcast) Init(address string, targetIpAddresses []string, debug bool) *MajorityAckUniformReliableBroadcast {
	urb.indChannel = make(chan messages.IndMessage)
	urb.reqChannel = make(chan messages.ReqMessage)
	urb.ipAddress = address
	urb.delivered = make(map[string] bool)
	urb.pending = make(map[string] bool)
	urb.ack = make(map[string] int)
	urb.targetIpAddresses = targetIpAddresses
	urb.numberOfProcesses = len(targetIpAddresses)
	urb.beb = *NewBestEffortBroadcast(debug)

	return urb
}

func NewMajorityAckUniformReliableBroadcast(address string, targetIpAddresses []string, debug bool) *MajorityAckUniformReliableBroadcast {
	return new(MajorityAckUniformReliableBroadcast).Init(address, targetIpAddresses, debug)
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
	key := formatKey(urb.IpAddress(), message.Message())
	reqMsg := *messages.NewReqMessage(message.To(), key)
	urb.pending[key] = true
	urb.beb.PushReqMessageToChannel(reqMsg)
}

func (urb *MajorityAckUniformReliableBroadcast) KeepDelivering() {
	keys := utils.Filter(utils.Keys(urb.pending), func (v string) bool {
		return urb.pending[v]
	})

	for i := 0; i < len(keys); i++ {
		k := keys[i]
		_, ok := urb.delivered[k]
		if !ok && urb.canDeliver(k) {
			urb.pending[k] = false
			urb.delivered[k] = true
			urb.PushIndMessageToChannel(unformatKey(keys[i]))
		}
	}
}

func (urb *MajorityAckUniformReliableBroadcast) bebDeliver(message messages.IndMessage) {
	key := message.Message()
	count, okAck := urb.ack[key]

	if okAck {
		urb.ack[key] = count + 1
	} else {
		urb.ack[key] = 1
	}
	_, ok := urb.pending[key]

	if !ok {
		reqMsg := *messages.NewReqMessage(urb.targetIpAddresses, message.Message())
		urb.beb.PushReqMessageToChannel(reqMsg)
		urb.pending[key] = true
	}

	urb.KeepDelivering()
}

func (urb *MajorityAckUniformReliableBroadcast) canDeliver(key string) bool {
	return urb.ack[key] > urb.numberOfProcesses/2
}

func formatKey(ip string, message string) string {
	return ip + ";" + message
}

func unformatKey(key string) messages.IndMessage {
	s := strings.Split(key, ";")
	indMsg := *messages.NewIndMessage(s[0], s[1])
	return indMsg
}
