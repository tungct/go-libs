package main

import (
	"fmt"
	"net"
	"encoding/gob"
	"github.com/tungct/go-libs/messqueue"
	"github.com/tungct/go-libs/topic"
	"log"
)

var	 Topics [] topic.Topic

const lenTopic = 10

//func HandleConn(conn net.Conn){
//	for {
//		dec := gob.NewDecoder(conn)
//		mess := &messqueue.Message{}
//
//		// keep listen message from this connect client
//		err := dec.Decode(mess)
//
//		if err != nil{
//			break
//		}
//
//		fmt.Printf("Received : %+v \n", mess);
//
//		// status 1 : Init connect
//		if mess.Status == 1 {
//			conn.Write([]byte("OK"))
//			defer conn.Close()
//
//			// status 2 : Publish message
//		} else if mess.Status == 2 {
//			topicName := topic.RuleTopic(*mess)
//			indexTopic := topic.GetIndexTopic(topicName, Topics)
//
//			// if topicName is not in topics, create new topic
//			if indexTopic == -1 {
//				topic := Topic.InitTopic(topicName, lenTopic)
//				Topic.PublishToTopic(topic, *mess)
//				Topics = append(Topics, topic)
//			} else {
//				Topic.PublishToTopic(Topics[indexTopic], *mess)
//			}
//			conn.Write([]byte("Success"))
//			defer conn.Close()
//
//			// status 3 : Subscribe message
//		} else if mess.Status == 3 {
//			var messResponse messqueue.Message
//			//topicName, _ := strconv.Atoi(mess.Content)
//			topicName := Topic.RuleTopic(*mess)
//			for {
//				indexTopic := Topic.GetIndexTopic(topicName, Topics)
//
//				// if exits topicName in topics
//				if indexTopic != -1 {
//					// if not message in topic
//					if len(Topics[indexTopic].MessQueue) != 0 {
//						_, messResponse = Topic.Subscribe(Topics[indexTopic])
//						encoder := gob.NewEncoder(conn)
//						er := encoder.Encode(messResponse)
//						if er != nil{
//							break
//						}
//						defer conn.Close()
//					} else {
//						messResponse = messqueue.CreateMessage(messqueue.NilMessageStatus, "Not message in topic")
//						encoder := gob.NewEncoder(conn)
//						er := encoder.Encode(messResponse)
//						if er != nil{
//							break
//						}
//						defer conn.Close()
//					}
//				} else {
//					messResponse = messqueue.CreateMessage(messqueue.NilMessageStatus, "Not exits topic")
//					encoder := gob.NewEncoder(conn)
//					er := encoder.Encode(messResponse)
//					if er != nil{
//						break
//					}
//					defer conn.Close()
//				}
//				fmt.Println("Subscribe Topic ", topicName)
//			}
//		}
//		Topic.PrintTopic(Topics)
//	}
//	return
//}

// Server handle connection from client
func HandleConnection(conn net.Conn) {

	dec := gob.NewDecoder(conn)
	mess := &messqueue.Message{}
	dec.Decode(mess)
	fmt.Printf("Received : %+v \n", mess);

	// Status 2 : publish message
	// if client is publisher
	if mess.Status == 2 {
		topicName := topic.RuleTopic(*mess)
		indexTopic := topic.GetIndexTopic(topicName, Topics)

		// if topicName is not in topics, create new topic
		if indexTopic == -1 {
			tp := topic.InitTopic(topicName, lenTopic)
			topic.PublishToTopic(tp, *mess)
			Topics = append(Topics, tp)
		} else {
			topic.PublishToTopic(Topics[indexTopic], *mess)
		}
		conn.Write([]byte("Success"))
		topic.PrintTopic(Topics)

		// keep listen message from this client
		for {
			//dec = gob.NewDecoder(conn)
			err := dec.Decode(mess)

			//if client closed connect, break keep listen
			if err != nil {
				break
			}
			topicName := topic.RuleTopic(*mess)
			indexTopic := topic.GetIndexTopic(topicName, Topics)

			// if topicName is not in topics, create new topic
			if indexTopic == -1 {
				tp := topic.InitTopic(topicName, lenTopic)
				topic.PublishToTopic(tp, *mess)
				Topics = append(Topics, tp)
			} else {
				topic.PublishToTopic(Topics[indexTopic], *mess)
			}

			// send message success to client
			conn.Write([]byte("Success"))

			topic.PrintTopic(Topics)
			defer conn.Close()
		}
		// status 3 : Subscribe message
	}else if mess.Status == 3{
		var messResponse messqueue.Message

		topicName := topic.RuleTopic(*mess)
		indexTopic := topic.GetIndexTopic(topicName, Topics)

		// init encode to send message to client
		encoder := gob.NewEncoder(conn)
		for {
			if len(Topics[indexTopic].MessQueue) != 0 {
				_, messResponse = topic.Subscribe(Topics[indexTopic])
				fmt.Println("Message send to subscriber : ", messResponse)
				er := encoder.Encode(messResponse)
				if er != nil{
					log.Fatal(er)
					break
				}
				fmt.Println("Subscribe Topic ", topicName)
				topic.PrintTopic(Topics)
			} else {
				continue
			}

		}
		defer conn.Close()
	}

	return
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