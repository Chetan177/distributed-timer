package main

import (
	"bytes"
	"encoding/json"
	"github.com/go-stomp/stomp/v3"
	"log"
	"net"
	"net/http"

	"dtimer/model"
)

const queueName = "/queue/timer"
const activeMQURL = "localhost:61613"

var activeMQ *stomp.Conn
var options []func(*stomp.Conn) error = []func(*stomp.Conn) error{
	stomp.ConnOpt.Login("admin", "admin"),
}

func main() {

	netConn, err := net.Dial("tcp", activeMQURL)
	if err != nil {
		log.Fatalln("activemq dial fail", err)
	}
	activeMQ, err = stomp.Connect(netConn, options...)
	if err != nil {
		log.Fatalln("activemq dial fail", err)
	}

	sub, err := activeMQ.Subscribe(queueName, stomp.AckAuto)
	if err != nil {
		log.Fatalln("cannot subscribe to", queueName, err.Error())
		return
	}

	log.Println(" Debug: Consumer started")
	for {
		msg := <-sub.C
		data := &model.QueueData{}
		err := json.Unmarshal(msg.Body, data)
		if err != nil {
			log.Println("Error: ", err)
		}

		log.Println("Debug: sending callback for timerID: ", data.TimerID)
		sendCallback(data, msg.Body)

	}

}

func sendCallback(data *model.QueueData, body []byte) {
	bodyReader := bytes.NewReader(body)
	_, err := http.Post(data.CallbackURL, "application/json", bodyReader)
	if err != nil {
		log.Println("Error: ", err)
	}
}
