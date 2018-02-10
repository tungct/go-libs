package main

import (
	"fmt"
	"net"
	"encoding/gob"
	"github.com/tungct/go-libs/messqueue"
	"github.com/tungct/go-libs/Topic"
)

var Topics [] Topic.Topic


func HandleConnection(conn net.Conn) {
	dec := gob.NewDecoder(conn)
	mess := &messqueue.Message{}
	dec.Decode(mess)

	indexTopic := Topic.GetIndexTopic(mess.Status, Topics)

	if indexTopic == -1{
		topic := Topic.InitTopic(mess.Status, 10)
		Topics = append(Topics, topic)
	}else{
		Topics[indexTopic].MessQueue <- mess
	}
	fmt.Println(len(Topics))
	conn.Write([]byte("OK"))
	fmt.Printf("Received : %+v", mess);
	conn.Close()
}

func main() {
	fmt.Println("Server listion at port 8080");
	ln, err := net.Listen("tcp", ":8080")
	if err != nil {
		// handle error
	}
	for {
		conn, err := ln.Accept() // this blocks until connection or error
		if err != nil {
			// handle error
			continue
		}
		go HandleConnection(conn) // a goroutine handles conn so that the loop can accept other connections
	}
}