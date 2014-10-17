package main

import (
	"flag"
	"fmt"
)


func mockserver(sendPtr, receivePtr *chan []byte){
	send := *sendPtr
	receive := *receivePtr
	for {
		select {
		case raw_msg := <-send:
			msg := string(raw_msg)
			if msg == string(createEmptyReadRequest()){
				receive <- []byte(`<?xml version="1.0" encoding="UTF-8"?><qlmEnvelope version="0.2" ttl="0"><response><result><return returnCode="200"></return></result></response></qlmEnvelope>`)
			}else if msg == string(createReadRequest("SmartFridge1", "PowerConsumption")){
				receive <- []byte(`<qlm:qlmEnvelope xmlns:qlm="QLMmi.xsd" xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance" xsi:schemaLocation="QLMmi.xsd QLMmi.xsd" version="1.0" ttl="10"><qlm:response><qlm:result msgformat="QLMdf"><qlm:return returnCode="200"></qlm:return><qlm:msg xmlns="QLMdf.xsd" xsi:schemaLocation="QLMdf.xsd QLMdf.xsd"><Objects><Object><id>SmartFridge1</id><InfoItem name="PowerConsumption"><value type="xs:int" unixTime="5453563">43</value></InfoItem></Object></Objects></qlm:msg></qlm:result></qlm:response>`)
			}else if msg == string(createSubscriptionRequest("SmartFridge1", "PowerConsumption", "-1")){
				receive <- []byte(`<qlm:qlmEnvelope xmlns:qlm="QLMmi.xsd" xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance" xsi:schemaLocation="QLMmi.xsd QLMmi.xsd" version="1.0" ttl="0"><qlm:response><qlm:result><qlm:return returnCode="200"></qlm:return><qlm:requestId>REQ1</qlm:requestId></qlm:result></qlm:response></qlm:qlmEnvelope>`)
			}else if msg == string(createReadSubscriptionRequest("REQ1")){
				receive <- []byte(`<qlm:qlmEnvelope xmlns:qlm="QLMmi.xsd" xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance" xsi:schemaLocation="QLMmi.xsd QLMmi.xsd" version="1.0" ttl="10"><qlm:response><qlm:result msgformat="QLMdf"><qlm:return returnCode="200"></qlm:return><qlm:requestId>REQ1</qlm:requestId><qlm:msg xmlns="QLMdf.xsd" xsi:schemaLocation="QLMdf.xsd QLMdf.xsd"><Objects><Object><id>SmartFridge1</id><InfoItem name="PowerConsumption"><value type="xs:int" unixTime="5453563">43</value><value type="xs:int" unixTime="5453584">47</value></InfoItem></Object></Objects></qlm:msg></qlm:result></qlm:response>`)
			}else if msg == string(createWriteRequest("SmartFridge1", "FridgeTemperatureSetpoint", "6")){
				receive <- []byte(`<qlm:qlmEnvelope xmlns:qlm="QLMmi.xsd" xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance" xsi:schemaLocation="QLMmi.xsd QLMmi.xsd" version="1.0" ttl="0"><qlm:response><qlm:result><qlm:return returnCode="200"></qlm:return></qlm:result></qlm:response></qlm:qlmEnvelope>`)
			}else if msg == string(createCancelSubscriptionRequest("REQ1")){
				receive <- []byte(`<qlm:qlmEnvelope xmlns:qlm="QLMmi.xsd" xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance" xsi:schemaLocation="QLMmi.xsd QLMmi.xsd" version="1.0" ttl="0"><qlm:response><qlm:result><qlm:return returnCode="200"></qlm:return></qlm:result></qlm:response></qlm:qlmEnvelope>`)
			}else{
				// Response 404
				receive <- []byte(`<qlm:qlmEnvelope xmlns:qlm="QLMmi.xsd" xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance" xsi:schemaLocation="QLMmi.xsd QLMmi.xsd" version="1.0" ttl="0"><qlm:response><qlm:result><qlm:return returnCode="404"></qlm:return></qlm:result></qlm:response></qlm:qlmEnvelope>`)
			}
		}
	}
}
func createServerConnection(send, receive *chan []byte){
	go mockserver(send, receive)
}
/*
Usage
cli test
cli read id name
cli write id name value
cli order id name interval
cli order-get req_id
cli order-cancel req_id

TODO:
cli read-meta path ?


Supported commands:

cli test
cli read SmartFridge1 PowerConsumption
cli order SmartFridge1 PowerConsumption -1
cli order-get REQ1
cli order-cancel REQ1
cli write SmartFridge1 FridgeTemperatureSetpoint 6

*/

func main() {
	var receive chan []byte
	var send chan []byte

	send = make(chan []byte)
	receive = make(chan []byte)

	flag.Parse()

	createServerConnection(&send, &receive)

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
	fmt.Println("Received: ", string(msg))
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

func createReadRequest(id, name string) []byte{
	return []byte(fmt.Sprintf(`<?xml version="1.0" encoding="UTF-8"?>
<qlm:qlmEnvelope xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance"
xmlns:qlm="QLMmi.xsd" xsi:schemaLocation="QLMmi.xsd QLMmi.xsd" version="1.0"
ttl="10">
<qlm:read msgformat="QLM_mf.xsd">
<qlm:msg xmlns="QLMdf.xsd" xsi:schemaLocation="QLMdf.xsd QLMdf.xsd">
<Objects>
<Object>
<id>%s</id>
<InfoItem name="%s">
</InfoItem>
</Object>
</Objects>
</qlm:msg>
</qlm:read>
</qlm:qlmEnvelope>`, id, name))
}

func createSubscriptionRequest(id, name, interval string) []byte{
	return []byte(fmt.Sprintf(`<?xml version="1.0" encoding="UTF-8"?>
<qlm:qlmEnvelope xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance"
xmlns:qlm="QLMmi.xsd" xsi:schemaLocation="QLMmi.xsd QLMmi.xsd" version="1.0"
ttl="10">
<qlm:read msgformat="QLM_mf.xsd" interval="%s">
<qlm:msg xmlns="QLMdf.xsd" xsi:schemaLocation="QLMdf.xsd QLMdf.xsd">
<Objects>
<Object>
<id>%s</id>
<InfoItem name="%s">
</InfoItem>
</Object>
</Objects>
</qlm:msg>
</qlm:read>
</qlm:qlmEnvelope>`, interval, id, name))
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
<Objects>
<Object>
<id>%s</id>
<InfoItem name="%s"><value>%s</value></InfoItem>
</Object>
</Objects>
</qlm:msg>
</qlm:write>
</qlm:qlmEnvelope>`, id, name, value))
}