package main

import (
	"flag"
	"fmt"
	"strconv"
)

func main() {
	var receive chan []byte
	var send chan []byte

	send = make(chan []byte)
	receive = make(chan []byte)

	var address string
	flag.StringVar(&address, "server", "http://localhost/qlm/", "Server address")

	flag.Parse()

	if err := createServerConnection(address, &send, &receive); err != nil {
		fmt.Println(err)
		return
	}

	command := flag.Arg(0)
	switch command {
	case "test":
		{
			send <- createEmptyReadRequest()
		}
	case "read":
		{
			id := flag.Arg(1)
			name := flag.Arg(2)
			send <- createReadRequest(id, name)
		}
	case "write":
		{
			id := flag.Arg(1)
			name := flag.Arg(2)
			value := flag.Arg(3)
			send <- createWriteRequest(id, name, value)
		}
	case "order":
		{
			id := flag.Arg(1)
			name := flag.Arg(2)
			interval, _ := strconv.ParseFloat(flag.Arg(3), 32)
			send <- createSubscriptionRequest(id, name, interval)
		}
	case "order-get":
		{
			requestId := flag.Arg(1)
			send <- createReadSubscriptionRequest(requestId)
		}
	case "order-cancel":
		{
			requestId := flag.Arg(1)
			send <- createCancelSubscriptionRequest(requestId)
		}
	default:
		{
			fmt.Println("Unknown command.")
			fmt.Println("Usage:")
			fmt.Println("cli [--server http://localhost/qlm/] test")
			fmt.Println("cli [--server http://localhost/qlm/] read id name")
			fmt.Println("cli [--server http://localhost/qlm/] write id name value")
			fmt.Println("cli [--server http://localhost/qlm/] order id name interval")
			fmt.Println("cli [--server http://localhost/qlm/] order-get req_id")
			fmt.Println("cli [--server http://localhost/qlm/] order-cancel req_id")
			return
		}
	}

	msg := <-receive
	fmt.Println(string(msg))
}
