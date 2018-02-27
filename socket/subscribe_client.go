package main

import (
	"fmt"
	"strings"
	"strconv"
	"net"
	"log"
	"encoding/gob"
	"github.com/tungct/go-libs/messqueue"
)

// send subscribe message to server and get message response
func GetMess(ip string, port int, topicName string){
	addr := strings.Join([]string{ip, strconv.Itoa(port)}, ":")
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		log.Fatal("Connection error", err)
	}

	// subscribe topic with topicName
	fmt.Println("Subscribe topic : ", topicName)
	mess := messqueue.CreateMessage(messqueue.SubscribeStatus, topicName)
	encoder := gob.NewEncoder(conn)
	encoder.Encode(mess)

	for {
		messRes := &messqueue.Message{}
		dec := gob.NewDecoder(conn)
		err = dec.Decode(messRes)
		if err != nil{
			log.Fatal(err)
			break
		}
		fmt.Println("Received message : ", messRes)
	}
	defer conn.Close()
}

func main() {
	var (
		ip   = "127.0.0.1"
		port = 8080
	)
	fmt.Println("subscribe client");
	GetMess(ip, port, "other")
}
