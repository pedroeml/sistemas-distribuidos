package perfect

import (
	"./messages"
	"fmt"
	"net"
)

type Link struct {
	indChannel chan messages.IndMessage
	reqChannel chan messages.ReqMessage
	isRunning  bool
	cache      map[string] net.Conn
}

func (ppl *Link) Init() *Link {
	if ppl.IsRunning() {
		return nil
	}

	ppl.indChannel = make(chan messages.IndMessage)
	ppl.reqChannel = make(chan messages.ReqMessage)
	ppl.cache = make(map[string] net.Conn)
	ppl.isRunning = true

	return ppl
}

func NewPerfectLink() *Link {
	return new(Link).Init()
}

func (ppl *Link) IsRunning() bool {
	return ppl.isRunning
}

func (ppl *Link) PushIndMessageToChannel(message messages.IndMessage) {
	ppl.indChannel <- message
}

func (ppl *Link) PopIndMessageFromChannel() messages.IndMessage {
	msg := <- ppl.indChannel
	return msg
}

func (ppl *Link) PushReqMessageToChannel(message messages.ReqMessage) {
	ppl.reqChannel <- message
}

func (ppl *Link) Start(address string) {
	go ppl.EstablishConnection(address)
	go ppl.KeepSending()
}

func (ppl *Link) EstablishConnection(address string) {
	fmt.Printf("DEBUG: Starting to listen to %s\n", address)
	ln, err := net.Listen("tcp4", address)

	if err != nil {
		fmt.Printf("DEBUG: Error on establishing connection to %s\n", address)
	} else {
		for {
			ppl.AcceptConnection(ln)
		}
	}
}

func (ppl *Link) AcceptConnection(ln net.Listener) {
	conn, err := ln.Accept()

	go ppl.Listen(conn, err)
}

func (ppl *Link) Listen(conn net.Conn, err error) {
	for {
		if err != nil {
			fmt.Printf("DEBUG: Error on listening accept\n")
			continue
		}

		buf := make([]byte, 1024)
		len, _ := conn.Read(buf)
		
		if len > 0 {
			content := make([]byte, len)
			copy(content, buf)

			msg := *messages.NewIndMessage(conn.RemoteAddr().String(), string(content))
			fmt.Printf("DEBUG: Received message \"%s\" from %s\n", msg.Message(), msg.From())
			ppl.indChannel <- msg	
		}
	}
}

func (ppl *Link) KeepSending() {
	for {
		message := <- ppl.reqChannel
		ppl.Send(message)
	}
}

func (ppl *Link) Send(message messages.ReqMessage) {
	conn, ok := ppl.cache[message.To()]

	if ok {
		fmt.Printf("DEBUG: Reusing connection to %s\n", message.To())
	} else {
		fmt.Printf("DEBUG: Opening connection to %s\n", message.To())
		conn = ppl.OpenConnection(message)
	}

	fmt.Printf("DEBUG: Sending message \"%s\" to %s\n", message.Message(), message.To())

	fmt.Fprintf(conn, message.Message())
}

func (ppl *Link) OpenConnection(message messages.ReqMessage) net.Conn {
	fmt.Printf("DEBUG: Starting dialing to %s\n", message.To())
	conn, err := net.Dial("tcp4", message.To())

	if err != nil {
		fmt.Printf("DEBUG: Error on dialing\n")
		return nil
	} else {
		fmt.Printf("DEBUG: Caching dialed connection\n")
		ppl.cache[message.To()] = conn
	}

	return conn
}
