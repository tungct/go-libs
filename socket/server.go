package main

import (
	"fmt"
	"net"
	"encoding/gob"
	"github.com/tungct/go-libs/messqueue"
	"github.com/tungct/go-libs/Topic"
	"strconv"
)

var Topics [] Topic.Topic


func HandleConnection(conn net.Conn) {
	dec := gob.NewDecoder(conn)
	mess := &messqueue.Message{}
	dec.Decode(mess)

	// status 1 : Init connect
	if mess.Status == 1{
		conn.Write([]byte("OK"))
		fmt.Printf("Received : %+v", mess);
		conn.Close()

	// status 2 : Publish message
	}else if mess.Status == 2{
		indexTopic := Topic.GetIndexTopic(mess.Status, Topics)

		if indexTopic == -1{
			topic := Topic.InitTopic(mess.Status, 10)
			Topic.PublishToTopic(topic, *mess)
			Topics = append(Topics, topic)
		}else{
			Topic.PublishToTopic(Topics[indexTopic], *mess)
		}
		fmt.Println(len(Topics))
		conn.Write([]byte("Success"))
		fmt.Printf("Received : %+v", mess);
		conn.Close()

	// status 3 : Subscribe message
	}else if mess.Status == 3{
		var messResponse messqueue.Message
		topicName, _ := strconv.Atoi(mess.Content)
		indexTopic := Topic.GetIndexTopic(topicName, Topics)
		_, messResponse = Topic.Subscribe(Topics[indexTopic])
		encoder := gob.NewEncoder(conn)
		encoder.Encode(messResponse)
		conn.Close()
		return
	}
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