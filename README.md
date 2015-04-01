# QLM(Quantum Lifecycle Mechanism) Command Line Client (cli)

Qlm-iot/cli is a command line program that allows communication with
server that is running [qlm-iot/core](https://github.com/qlm-iot/core) package. Current implementation
supports communication over WebSocket and HTTP layers. The implementation is capable of
sending messages to the specified server and reading response. Response from server is
written to the program standard output stream without interpreting its content in any way.

This project was done for the T-106.5700 course at Aalto University.

## Code structure
The implementation has two main tasks: Handling connection over specified
protocol and the construction of valid QLM messages.

### QLM message building
QLM messages are constructed using [qlm-iot/qlm](https://github.com/qlm-iot/qlm) packet. Each QLM message is XML-encapsulated
message that contains QLM Messaging Interface structure containing QLM Data
Format structure. Construction of messages is handled in messages.go source
code file.

### Connection handling
All connection related code is separated using golang channels. Each connection
is represented as two channels (one for sending requests and other for receiving
responses). This abstraction allows adding of new communication protocols that
QLM messages can be transmitted on top of.

This implementation supports communication over WebSocket and HTTP layers.
HTTP communication is done according to QLM Messaging Interface document using
HTTP POST requests. Websocket communication is our custom protocol that is only
supported by our [qlm-iot/core](https://github.com/qlm-iot/core) server.

All connection handling related code is located in connections.go source code file.


## Compilation
```
> go get github.com/qlm-iot/cli
> cd $GOPATH/src/github.com/qlm-iot/cli/
> go build
> ./cli
Unknown command.
Usage:
cli [--server http://localhost/qlm/] test
cli [--server http://localhost/qlm/] read id name
cli [--server http://localhost/qlm/] write id name value
cli [--server http://localhost/qlm/] order id name interval
cli [--server http://localhost/qlm/] order-get req_id
cli [--server http://localhost/qlm/] order-cancel req_id
```

## Examples
Following examples assume that you have [qlm-iot/core](https://github.com/qlm-iot/core) QLM server running at localhost:8000.

Test connection (sends empty request):
```
> ./cli --server http://localhost:8000/qlm/ read oven_12 temperature
<omiEnvelope version="1.0" ttl="0">
    <response>
        <result>
            <return returnCode="200"></return>
        </result>
    </response>
</omiEnvelope>
```
Read value `temperature` from node `oven_12` over http:
```
> ./cli --server http://localhost:8000/qlm/ read oven_12 temperature
<omiEnvelope version="1.0" ttl="0">
    <response>
        <result>
            <return returnCode="200"></return>
            <msg><Objects xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance" xsi:noNamespaceSchemaLocation="odf.xsd">
    <Object>
        <id>oven_12</id>
        <InfoItem name="temperature">
            <value unixTime="1427893342">82</value>
        </InfoItem>
    </Object>
</Objects></msg>
        </result>
    </response>
</omiEnvelope>
```
All requests can easily be done over WebSocket connections by changing the
connection protocol. For example running same request over WebSocket protocol:
```
> ./cli --server ws://localhost:8000/qlm/ read oven_12 temperature
```
Write value `100` to variable `target_temperature` in node `oven_12`:
```
> ./cli --server http://localhost:8000/qlm/ write oven_12 target_temperature 100
<omiEnvelope version="1.0" ttl="0">
    <response>
        <result>
            <return returnCode="200"></return>
        </result>
    </response>
</omiEnvelope>
```
Create subscription for variable `temperature` value for every `10` seconds in node `oven_12`:
```
> ./cli --server http://localhost:8000/qlm/ order oven_12 temperature 10
<omiEnvelope version="1.0" ttl="0">
    <response>
        <result>
            <return returnCode="200"></return>
            <requestId>REQ0000007</requestId>
        </result>
    </response>
</omiEnvelope>
```
Read values from subscription with id `REQ0000007`:
```
> ./cli --server http://localhost:8000/qlm/ order-get REQ0000007
<omiEnvelope version="1.0" ttl="0">
    <response>
        <result>
            <return returnCode="200"></return>
            <msg><Objects xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance" xsi:noNamespaceSchemaLocation="odf.xsd">
    <Object>
        <id>oven_12</id>
        <InfoItem name="temperature">
            <value unixTime="1427893495">83</value>
        </InfoItem>
    </Object>
    <Object>
        <id>oven_12</id>
        <InfoItem name="temperature">
            <value unixTime="1427893508">83</value>
        </InfoItem>
    </Object>
    <Object>
        <id>oven_12</id>
        <InfoItem name="temperature">
            <value unixTime="1427893517">84</value>
        </InfoItem>
    </Object>
</Objects></msg>
        </result>
    </response>
</omiEnvelope>
```
Cancel subscription with id `REQ0000007`:
```
> ./cli --server http://localhost:8000/qlm/ order-cancel REQ0000007
<omiEnvelope version="1.0" ttl="0">
    <response>
        <result>
            <return returnCode="200"></return>
        </result>
    </response>
</omiEnvelope>
```

## Future work
 - It might be useful to have a way to explicit define timeouts for the execution
of QLM requests.
 - New communication layers that QLM messages can be transmitted on top of.
 - The program return code should be negative in case there is an error while
 handling request.

## Limitations
Current implementation doesn't support writing multiple values to single node
in single request. Server interprets missing values as removed values, so all
other values get removed from that node.

Our server implementation doesn't support call-back based subscriptions and
this same limitation also also present in this implementation.
