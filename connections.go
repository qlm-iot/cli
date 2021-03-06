package main

import (
	"errors"
	"github.com/gorilla/websocket"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
)

func httpServerConnector(address string, sendPtr, receivePtr *chan []byte) {
	send := *sendPtr
	receive := *receivePtr
	for {
		select {
		case rawMsg := <-send:
			msg := string(rawMsg)
			data := url.Values{}
			data.Set("msg", msg)
			response, err := http.PostForm(address, data)
			if err == nil {
				defer response.Body.Close()
				content, err := ioutil.ReadAll(response.Body)
				if err == nil {
					receive <- content
				} else {
					receive <- []byte(err.Error())
				}
			} else {
				receive <- []byte(err.Error())
			}
		}
	}
}
func wsServerConnector(address string, sendPtr, receivePtr *chan []byte) {
	send := *sendPtr
	receive := *receivePtr
	for {
		select {
		case rawMsg := <-send:
			var h http.Header

			conn, _, err := websocket.DefaultDialer.Dial(address, h)
			if err == nil {
				if err := conn.WriteMessage(websocket.BinaryMessage, rawMsg); err != nil {
					receive <- []byte(err.Error())
				}
				_, content, err := conn.ReadMessage()
				if err == nil {
					receive <- content
				} else {
					receive <- []byte(err.Error())
				}
			} else {
				receive <- []byte(err.Error())
			}
		}
	}
}
func createServerConnection(address string, send, receive *chan []byte) error {
	if strings.HasPrefix(address, "http://") {
		go httpServerConnector(address, send, receive)
	} else if strings.HasPrefix(address, "ws://") {
		go wsServerConnector(address, send, receive)
	} else {
		return errors.New("Unsupported server protocol. Supported protocols are http and ws.")
	}
	return nil
}
