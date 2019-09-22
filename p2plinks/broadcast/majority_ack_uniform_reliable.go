package broadcast

import (
	"../utils"
	"./messages"
	"fmt"
	"strings"
	"sync"
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
	mutex             *sync.Mutex
}

func (urb *MajorityAckUniformReliableBroadcast) Init(address string, targetIpAddresses []string) *MajorityAckUniformReliableBroadcast {
	urb.indChannel = make(chan messages.IndMessage)
	urb.reqChannel = make(chan messages.ReqMessage)
	urb.ipAddress = address
	urb.delivered = make(map[string] bool)
	urb.pending = make(map[string] bool)
	urb.ack = make(map[string] int)
	urb.targetIpAddresses = targetIpAddresses
	urb.numberOfProcesses = len(targetIpAddresses)
	urb.mutex = &sync.Mutex{}
	urb.beb = *NewBestEffortBroadcast()

	return urb
}

func NewMajorityAckUniformReliableBroadcast(address string, targetIpAddresses []string) *MajorityAckUniformReliableBroadcast {
	return new(MajorityAckUniformReliableBroadcast).Init(address, targetIpAddresses)
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
	key := formatKey(urb.IpAddress(), message.Message())
	reqMsg := *messages.NewReqMessage(message.To(), key)
	urb.mutex.Lock()
	urb.pending[key] = true
	urb.beb.PushReqMessageToChannel(reqMsg)
	urb.mutex.Unlock()
}

func (urb *MajorityAckUniformReliableBroadcast) KeepDelivering() {
	for {
		urb.mutex.Lock()
		keys := utils.Filter(utils.Keys(urb.pending), func (v string) bool {
			return urb.pending[v]
		})
		urb.mutex.Unlock()

		for i := 0; i < len(keys); i++ {
			urb.mutex.Lock()
			_, ok := urb.delivered[keys[i]]
			urb.mutex.Unlock()
			if !ok && urb.canDeliver(keys[i]) {
				urb.mutex.Lock()
				urb.pending[keys[i]] = false
				urb.delivered[keys[i]] = true
				urb.PushIndMessageToChannel(unformatKey(keys[i]))
				urb.mutex.Unlock()
			}
		}
	}
}

func (urb *MajorityAckUniformReliableBroadcast) bebDeliver(message messages.IndMessage) {
	indMsg := unformatKey(message.Message())
	fmt.Printf("URB FROM %s MSG %s\n", indMsg.From(), indMsg.Message())
	key := message.Message()
	urb.mutex.Lock()
	count, okAck := urb.ack[key]

	if okAck {
		urb.ack[key] = count + 1
	} else {
		urb.ack[key] = 1
	}
	_, ok := urb.pending[key]
	urb.mutex.Unlock()

	if !ok {
		reqMsg := *messages.NewReqMessage(urb.targetIpAddresses, message.Message())
		urb.mutex.Lock()
		urb.pending[key] = true
		urb.beb.PushReqMessageToChannel(reqMsg)
		urb.mutex.Unlock()
	}
}

func (urb *MajorityAckUniformReliableBroadcast) canDeliver(key string) bool {
	urb.mutex.Lock()
	can := urb.ack[key] > urb.numberOfProcesses/2
	urb.mutex.Unlock()
	return can
}

func formatKey(ip string, message string) string {
	return ip + ";" + message
}

func unformatKey(key string) messages.IndMessage {
	s := strings.Split(key, ";")
	indMsg := *messages.NewIndMessage(s[0], s[1])
	return indMsg
}
