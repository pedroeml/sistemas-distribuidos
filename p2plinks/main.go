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
	addresses := utils.ParseAddressFile(os.Args[2])

	var target string

	if id == 0 {
		target = addresses[1]
	} else if id == 1 {
		target = addresses[0]
	} else {
		fmt.Printf("PROCESS %d: Only 2 addresses supported for now!", id)
		return
	}

	fmt.Printf("PROCESS %d: from %s is messaging to %s\n", id, addresses[id], target)

	ch := make(chan int)

	ppl := *perfect.NewPerfectLink()
	ppl.Start(addresses[id])

	go func() {
		for {
			fmt.Println("Type the message:")
			reader := bufio.NewReader(os.Stdin)
			msg, _ := reader.ReadString('\n')
			msg = msg[:len(msg) - 1]
			reqMsg := *messages.NewReqMessage(target, msg)
			fmt.Printf("PROCESS %d: sent message \"%s\" to %s\n", id, reqMsg.Message(), target)
			ppl.PushReqMessageToChannel(reqMsg)
		}
	}()

	go func() {
		for {
			msg := ppl.PopIndMessageFromChannel()
			fmt.Printf("PROCESS %d: got message \"%s\" from %s\n", id, msg.Message(), target)
		}
	}()

	<- ch
}
