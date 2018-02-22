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

func GetMess(ip string, port int, topicName string){
	addr := strings.Join([]string{ip, strconv.Itoa(port)}, ":")
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		log.Fatal("Connection error", err)
	}
	encoder := gob.NewEncoder(conn)

	// subscribe topic with topicName
	fmt.Println("Subscribe topic : ", topicName)
	mess := messqueue.CreateMessage(messqueue.SubscribeStatus, topicName)
	encoder.Encode(mess)
	dec := gob.NewDecoder(conn)
	messRes := &messqueue.Message{}
	dec.Decode(messRes)
	conn.Close()
	fmt.Println("Received message : ", messRes);
}

func main() {
	var (
		ip   = "127.0.0.1"
		port = 8080
	)
	fmt.Println("subscribe client");
	GetMess(ip, port, "Message")
}
