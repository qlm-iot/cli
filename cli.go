package main

import (
	"flag"
	"fmt"
	"github.com/qlm-iot/qlm/df"
	"io/ioutil"
	"net/http"
	"net/url"
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
func createServerConnection(address string, send, receive *chan []byte){
	go httpserverconnector(address, send, receive)
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

	createServerConnection(address, &send, &receive)

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
			interval := flag.Arg(3)
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
	return []byte(`<?xml version="1.0" encoding="UTF-8"?>
<qlm:qlmEnvelope version="1.0" ttl="10">
<qlm:read>
<qlm:msg>
</qlm:msg>
</qlm:read>
</qlm:qlmEnvelope>`)
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
	return []byte(fmt.Sprintf(`<?xml version="1.0" encoding="UTF-8"?>
<qlm:qlmEnvelope xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance"
xmlns:qlm="QLMmi.xsd" xsi:schemaLocation="QLMmi.xsd QLMmi.xsd" version="1.0"
ttl="10">
<qlm:read msgformat="QLM_mf.xsd">
<qlm:msg xmlns="QLMdf.xsd" xsi:schemaLocation="QLMdf.xsd QLMdf.xsd">
%s
</qlm:msg>
</qlm:read>
</qlm:qlmEnvelope>`, createQLMMessage(id, name)))
}

func createSubscriptionRequest(id, name, interval string) []byte{
	return []byte(fmt.Sprintf(`<?xml version="1.0" encoding="UTF-8"?>
<qlm:qlmEnvelope xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance"
xmlns:qlm="QLMmi.xsd" xsi:schemaLocation="QLMmi.xsd QLMmi.xsd" version="1.0"
ttl="10">
<qlm:read msgformat="QLM_mf.xsd" interval="%s">
<qlm:msg xmlns="QLMdf.xsd" xsi:schemaLocation="QLMdf.xsd QLMdf.xsd">
%s
</qlm:msg>
</qlm:read>
</qlm:qlmEnvelope>`, interval, createQLMMessage(id, name)))
}

func createReadSubscriptionRequest(requestId string) []byte{
	return []byte(fmt.Sprintf(`<?xml version="1.0" encoding="UTF-8"?>
<qlm:qlmEnvelope xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance"
xmlns:qlm="QLMmi.xsd" xsi:schemaLocation="QLMmi.xsd QLMmi.xsd" version="1.0"
ttl="10">
<qlm:read msgformat="QLM_mf.xsd">
<qlm:requestId>%s</qlm:requestId>
</qlm:read>
</qlm:qlmEnvelope>`, requestId))
}

func createCancelSubscriptionRequest(requestId string) []byte{
	return []byte(fmt.Sprintf(`<?xml version="1.0" encoding="UTF-8"?>
<qlm:qlmEnvelope xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance"
xmlns:qlm="QLMmi.xsd" xsi:schemaLocation="QLMmi.xsd QLMmi.xsd" version="1.0"
ttl="10">
<qlm:cancel>
<qlm:requestId>%s</qlm:requestId>
</qlm:cancel>
</qlm:qlmEnvelope>`, requestId))
}

func createWriteRequest(id, name, value string) []byte{
	return []byte(fmt.Sprintf(`<?xml version="1.0" encoding="UTF-8"?>
<qlm:qlmEnvelope xmlns:qlm="QLMmi.xsd"
xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance"
xsi:schemaLocation="QLMmi.xsd QLMmi.xsd" version="1.0" ttl="-1">
<qlm:write msgformat="QLMdf" targetType="device">
<qlm:msg xmlns="QLMdf.xsd" xsi:schemaLocation="QLMdf.xsd QLMdf.xsd">
%s
</qlm:msg>
</qlm:write>
</qlm:qlmEnvelope>`, createQLMMessageWithValue(id, name, value)))
}
