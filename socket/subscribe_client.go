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

func GetMess(ip string, port int, topicName int){
	addr := strings.Join([]string{ip, strconv.Itoa(port)}, ":")
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		log.Fatal("Connection error", err)
	}
	encoder := gob.NewEncoder(conn)
	mess := messqueue.SubscribeMessage(topicName)
	encoder.Encode(mess)
	dec := gob.NewDecoder(conn)
	messRes := &messqueue.Message{}
	dec.Decode(messRes)
	conn.Close()
	fmt.Println(messRes);
}

func main() {
	var (
		ip   = "127.0.0.1"
		port = 8080
	)
	fmt.Println("start client");
	GetMess(ip, port, 2)
}
