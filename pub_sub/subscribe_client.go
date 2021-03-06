package main

import (
	"fmt"
	"strings"
	"strconv"
	"net"
	"log"
	"encoding/gob"
	"github.com/tungct/go-libs/messqueue"
	"os"
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

	// create subscribe message send to server
	mess := messqueue.CreateMessage(messqueue.SubscribeStatus, topicName)

	// create nil message return after receive from server
	mes := messqueue.CreateMessage(messqueue.NilMessageStatus, topicName)
	// init encoder to send data to server
	encoder := gob.NewEncoder(conn)
	encoder.Encode(mess)

	// init decoder to read data from server response
	dec := gob.NewDecoder(conn)

	for {

		messRes := &messqueue.Message{}
		err = dec.Decode(messRes)
		if err != nil{
			log.Fatal(err)
			break
		}
		err = encoder.Encode(mes)
		if err != nil{
			log.Fatal(err)
			break
		}
		fmt.Println("Received message : ", messRes, "from topic " + topicName)
	}
	defer conn.Close()
}

func main() {
	var (
		ip   = "127.0.0.1"
		port = 8080
	)
	topicName := os.Args[1]
	fmt.Println("subscribe client");
	GetMess(ip, port, topicName)
}
