package main

import (
	"./perfect"
	"./perfect/messages"
	"./utils"
	"bufio"
	"fmt"
	"os"
	"strconv"
)

func main() {
	id, _ := strconv.Atoi(os.Args[1])
	ipAddresses := utils.ParseAddressFile(os.Args[2])

	var sourceIpAddress string
	var targetIpAddress string

	if id == 0 {
		sourceIpAddress = ipAddresses[0]
		targetIpAddress = ipAddresses[1]
	} else if id == 1 {
		sourceIpAddress = ipAddresses[1]
		targetIpAddress = ipAddresses[0]
	} else {
		fmt.Printf("PROCESS %d: Only 2 addresses supported for now!", id)
		return
	}

	fmt.Printf("PROCESS %d: from %s is messaging to %s\n", id, sourceIpAddress, targetIpAddress)

	ch := make(chan int)

	ppl := perfect.NewPerfectLink()
	ppl.Start(sourceIpAddress)

	go sendKeyboardInputMessage(ppl, id, targetIpAddress)
	go listenReceivingMessages(ppl, id, targetIpAddress)

	<- ch
}

func sendKeyboardInputMessage(ppl *perfect.Link, id int, targetIpAddress string) {
	for {
		fmt.Println("Type the message:")
		reader := bufio.NewReader(os.Stdin)
		msg, _ := reader.ReadString('\n')
		msg = msg[:len(msg) - 1]

		if len(msg) == 0 {
			continue
		}

		reqMsg := *messages.NewReqMessage(targetIpAddress, msg)
		fmt.Printf("PROCESS %d: sent message \"%s\" to %s\n", id, reqMsg.Message(), targetIpAddress)
		ppl.PushReqMessageToChannel(reqMsg)
	}
}

func listenReceivingMessages(ppl *perfect.Link, id int, targetIpAddress string) {
	for {
		msg := ppl.PopIndMessageFromChannel()
		fmt.Printf("PROCESS %d: got message \"%s\" from %s\n", id, msg.Message(), targetIpAddress)
	}
}
