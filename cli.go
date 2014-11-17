package main

import (
	"flag"
	"fmt"
	"github.com/gorilla/websocket"
	"github.com/qlm-iot/qlm/df"
	"github.com/qlm-iot/qlm/mi"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

func httpserverconnector(address string, sendPtr, receivePtr *chan []byte){
	send := *sendPtr
	receive := *receivePtr
	for {
		select {
		case raw_msg := <-send:
			msg := string(raw_msg)
			data := url.Values{}
			data.Set("msg", msg)
			response, err := http.PostForm(address, data)
			if err == nil{
				defer response.Body.Close()
				content, err := ioutil.ReadAll(response.Body)
				if err == nil{
					receive <- content
				}else{
					receive <- []byte(err.Error())
				}
			}else{
				receive <- []byte(err.Error())
			}
		}
	}
}
func wsServerConnector(address string, sendPtr, receivePtr *chan []byte){
	send := *sendPtr
	receive := *receivePtr
	for {
		select {
		case raw_msg := <-send:
			var h http.Header

			conn, _, err := websocket.DefaultDialer.Dial(address, h)
			if err == nil{
				if err := conn.WriteMessage(websocket.BinaryMessage, raw_msg); err != nil {
					receive <- []byte(err.Error())
				}
				_, content, err := conn.ReadMessage()
				if err == nil {
					receive <- content
				}else{
					receive <- []byte(err.Error())
				}
			}else{
				receive <- []byte(err.Error())
			}
		}
	}
}
func createServerConnection(address string, send, receive *chan []byte) bool{
	if strings.HasPrefix(address, "http://"){
		go httpserverconnector(address, send, receive)
	}else if strings.HasPrefix(address, "ws://"){
		go wsServerConnector(address, send, receive)
	}else{
		return false
	}
	return true
}
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

	if !createServerConnection(address, &send, &receive){
		fmt.Println("Unsupported server protocol")
		return
	}

	command := flag.Arg(0)
	switch command {
		case "test": {
			send <- createEmptyReadRequest()
		}
		case "read": {
			id := flag.Arg(1)
			name := flag.Arg(2)
			send <- createReadRequest(id, name)
		}
		case "write":{
			id := flag.Arg(1)
			name := flag.Arg(2)
			value := flag.Arg(3)
			send <- createWriteRequest(id, name, value)
		}
		case "order": {
			id := flag.Arg(1)
			name := flag.Arg(2)
			interval, _ := strconv.ParseFloat(flag.Arg(3), 32)
			send <- createSubscriptionRequest(id, name, interval)
		}
		case "order-get": {
			requestId := flag.Arg(1)
			send <- createReadSubscriptionRequest(requestId)
		}
		case "order-cancel": {
			requestId := flag.Arg(1)
			send <- createCancelSubscriptionRequest(requestId)
		}
		default: {
			fmt.Println("Unknown command")
			return
		}
	}

	msg := <-receive
	fmt.Println(string(msg))
}

func createEmptyReadRequest() []byte{
	ret, _ := mi.Marshal(mi.QlmEnvelope{
		Version: "1.0",
		Ttl:     -1,
		Read:    &mi.ReadRequest{},
	})
	return ret
}

func createQLMMessage(id, name string) string{
	objects := df.Objects{
		Objects: []df.Object{
			df.Object{
				Id:   &df.QLMID{Text: id},
				InfoItems: []df.InfoItem{
					df.InfoItem{
						Name: name,
					},
				},
			},
		},
	}
	data, _ := df.Marshal(objects)
	return (string)(data)
}

func createQLMMessageWithValue(id, name, value string) string{
	objects := df.Objects{
		Objects: []df.Object{
			df.Object{
				Id:   &df.QLMID{Text: id},
				InfoItems: []df.InfoItem{
					df.InfoItem{
						Name: name,
						Values: []df.Value{
							df.Value{
								Text: value,
							},
						},
					},
				},
			},
		},
	}
	data, _ := df.Marshal(objects)
	return (string)(data)
}

func createReadRequest(id, name string) []byte{
	ret, _ := mi.Marshal(mi.QlmEnvelope{
		Version: "1.0",
		Ttl:     -1,
		Read: &mi.ReadRequest{
			MsgFormat:  "QLMdf",
			Message:    &mi.Message{
				Data: createQLMMessage(id, name),
			},
		},
	})
	return ret
}

func createSubscriptionRequest(id, name string, interval float64) []byte{
	ret, _ := mi.Marshal(mi.QlmEnvelope{
		Version: "1.0",
		Ttl:     -1,
		Read: &mi.ReadRequest{
			MsgFormat:  "QLMdf",
			Interval:   interval,
			Message:    &mi.Message{
				Data: createQLMMessage(id, name),
			},
		},
	})
	return ret
}

func createReadSubscriptionRequest(requestId string) []byte{
	ret, _ := mi.Marshal(mi.QlmEnvelope{
		Version: "1.0",
		Ttl:     -1,
		Read: &mi.ReadRequest{
			MsgFormat:  "QLMdf",
			RequestIds: []mi.Id{
				mi.Id{Text: requestId},
			},
		},
	})
	return ret
}

func createCancelSubscriptionRequest(requestId string) []byte{
	ret, _ := mi.Marshal(mi.QlmEnvelope{
		Version: "1.0",
		Ttl:     -1,
		Cancel:  &mi.CancelRequest{
			RequestIds: []mi.Id{
				mi.Id{Text: requestId},
			},
		},
	})
	return ret
}

func createWriteRequest(id, name, value string) []byte{
	ret, _ := mi.Marshal(mi.QlmEnvelope{
		Version: "1.0",
		Ttl:     -1,
		Write:   &mi.WriteRequest{
			MsgFormat:  "QLMdf",
			TargetType: "device",
			Message: &mi.Message{
				Data: createQLMMessageWithValue(id, name, value),
			},
		},
	})
	return ret
}
