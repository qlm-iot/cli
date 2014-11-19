package main

import (
	"flag"
	"fmt"
	"strconv"
)

/*
Usage
cli [--server http://localhost/qlm/] test
cli [--server http://localhost/qlm/] read id name
cli [--server http://localhost/qlm/] write id name value
cli [--server http://localhost/qlm/] order id name interval
cli [--server http://localhost/qlm/] order-get req_id
cli [--server http://localhost/qlm/] order-cancel req_id
*/

func main() {
	var receive chan []byte
	var send chan []byte

	send = make(chan []byte)
	receive = make(chan []byte)

	var address string
	flag.StringVar(&address, "server", "http://localhost/qlm/", "Server address")

	flag.Parse()

	if !createServerConnection(address, &send, &receive) {
		fmt.Println("Unsupported server protocol")
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
			fmt.Println("Unknown command")
			return
		}
	}

	msg := <-receive
	fmt.Println(string(msg))
}
