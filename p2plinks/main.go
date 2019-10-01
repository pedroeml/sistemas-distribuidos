package main

import (
	"./broadcast"
	"./broadcast/messages"
	"./utils"
	"bufio"
	"fmt"
	"os"
	"strconv"
)

func main() {
	id, _ := strconv.Atoi(os.Args[1])
	ipAddresses := utils.ParseAddressFile(os.Args[2])

	if id >= len(ipAddresses) {
		fmt.Printf("ID out of bounds! Max ID allowed %d", len(ipAddresses) - 1)
		return
	}

	sourceIpAddress := ipAddresses[id]
	targetIpAddresses := make([]string, len(ipAddresses) -1, len(ipAddresses) -1)
	index := 0

	fmt.Printf("PROCESS %d: from %s is messaging to ", id, sourceIpAddress)
	for i := 0; i < len(ipAddresses); i++ {
		if i != id {
			targetIpAddresses[index] = ipAddresses[i]
			fmt.Printf("%s ", targetIpAddresses[index])
			index++
		}
	}
	fmt.Printf("\n")

	ch := make(chan int)

	urb := broadcast.NewMajorityAckUniformReliableBroadcast(sourceIpAddress, targetIpAddresses, false)
	urb.Start()

	go sendKeyboardInputMessage(urb, id, targetIpAddresses)
	go listenReceivingMessages(urb, id)

	<- ch
}

func sendKeyboardInputMessage(urb *broadcast.MajorityAckUniformReliableBroadcast, id int, targetIpAddresses []string) {
	for {
		fmt.Println("Type the message:")
		reader := bufio.NewReader(os.Stdin)
		msg, _ := reader.ReadString('\n')
		msg = msg[:len(msg) - 1]

		if len(msg) == 0 {
			continue
		}

		reqMsg := *messages.NewReqMessage(targetIpAddresses, msg)
		fmt.Printf("PROCESS %d: sent message \"%s\" to ", id, reqMsg.Message())
		for i := 0; i < len(targetIpAddresses); i++ {
			fmt.Printf("%s ", targetIpAddresses[i])
		}
		fmt.Printf("\n")

		urb.PushReqMessageToChannel(reqMsg)
	}
}

func listenReceivingMessages(urb *broadcast.MajorityAckUniformReliableBroadcast, id int) {
	for {
		msg := urb.PopIndMessageFromChannel()
		fmt.Printf("PROCESS %d: got message \"%s\" from %s\n", id, msg.Message(), msg.From())
	}
}
