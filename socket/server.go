package main

import (
	"fmt"
	"net"
	"encoding/gob"
	"github.com/tungct/go-libs/messqueue"
	"github.com/tungct/go-libs/topic"
	"sync"
)

type sub struct {
	encode *gob.Encoder
	decode *gob.Decoder
}
var mutex = &sync.Mutex{}
var	 Topics [] topic.Topic
var sbs chan sub
var subs [] sub
var wrk chan int

// Server handle connection from client
func HandleConnection(conn net.Conn) {

	dec := gob.NewDecoder(conn)
	mess := &messqueue.Message{}
	dec.Decode(mess)
	fmt.Printf("Received : %+v \n", mess);

	// Status 1 : init connect from client
	// if client is publisher
	if mess.Status == 1 {
		conn.Write([]byte("OK"))
		topicName := mess.Content
		indexTopic := topic.GetIndexTopic(topicName, Topics)

		// if topicName is not in topics, create new topic
		if indexTopic == -1 {
			tp := topic.InitTopic(topicName, messqueue.LenTopic)
			Topics = append(Topics, tp)
			indexTopic = topic.GetIndexTopic(topicName, Topics)
		}
		// keep listen message from this client
		for {
			//dec = gob.NewDecoder(conn)
			err := dec.Decode(mess)

			//if client closed connect, break keep listen
			if err != nil {
				break
			}

			// Status 2 : if message is publish
			if mess.Status == 2 {

				topic.PublishToTopic(Topics[indexTopic], *mess)

				// send message success to client
				conn.Write([]byte("Success"))

				topic.PrintTopic(Topics)
			}
			defer conn.Close()
		}
		// status 3 : Subscribe message
	}else if mess.Status == 3{
		encoder := gob.NewEncoder(conn)
		mutex.Lock()
		sub := sub{encoder, dec}
		subs = append(subs, sub)
		mutex.Unlock()

		//mutex.Lock()
		//encoder := gob.NewEncoder(conn)
		//subs = append(subs, encoder)
		//mutex.Unlock()

		//var messResponse messqueue.Message
		//
		//topicName := mess.Content
		//
		//fmt.Println("Subscribe Topic ", topicName)
		//indexTopic := topic.GetIndexTopic(topicName, Topics)
		//if indexTopic == -1 {
		//	tp := topic.InitTopic(topicName, messqueue.LenTopic)
		//	Topics = append(Topics, tp)
		//	indexTopic = topic.GetIndexTopic(topicName, Topics)
		//}
		//
		//mutex.Lock()
		//if sub[topicName] >= 1{
		//	sub[topicName] = sub[topicName] + 1
		//}else {
		//	sub[topicName] = 1
		//}
		//mutex.Unlock()
		//
		//// init encode to send message to client
		//encoder := gob.NewEncoder(conn)
		//for {
		//	if len(Topics[indexTopic].MessQueue) != 0 {
		//		_, messResponse = topic.Subscribe(Topics[indexTopic])
		//		fmt.Println("Message send to subscriber : ", messResponse)
		//		er := encoder.Encode(messResponse)
		//		if er != nil{
		//			fmt.Println("Connect fail, return message to topic")
		//			topic.PublishToTopic(Topics[indexTopic], messResponse)
		//			//log.Fatal(er)
		//			break
		//		}
		//		er = dec.Decode(mess)
		//		if er != nil{
		//			fmt.Println("Connect fail, return message to topic")
		//			topic.PublishToTopic(Topics[indexTopic], messResponse)
		//			//log.Fatal(er)
		//			break
		//		}
		//
		//		// listen message return from client after send success
		//
		//
		//		topic.PrintTopic(Topics)
		//	} else {
		//		continue
		//	}
		//
		//}
		//defer conn.Close()
	}
}

func removeItemArray(s []sub, index int) []sub {
	return append(s[:index], s[index+1:]...)
}

func Subscribe(w int, topicName string){
	message := &messqueue.Message{}
	indexTopic := topic.GetIndexTopic(topicName, Topics)
	if indexTopic == -1 {
		tp := topic.InitTopic(topicName, messqueue.LenTopic)
		Topics = append(Topics, tp)
		indexTopic = topic.GetIndexTopic(topicName, Topics)
	}
	if len(subs) >0 {
		mess := <-Topics[indexTopic].MessQueue
		for i := 0;i<len(subs);i++{
			fmt.Println(len(subs))
			er := subs[i].encode.Encode(mess)
			if er != nil {
				fmt.Println(er)
				fmt.Println("Connect fail, return message to ")
				if len(subs) >0{
					subs = removeItemArray(subs, i)
				}
				if len(subs) == 0{
					fmt.Println("3")
					topic.PublishToTopic(Topics[indexTopic], mess.(messqueue.Message))
				}
				break
			}
			er = subs[i].decode.Decode(message)
			fmt.Println(len(subs))
			if er != nil {
				fmt.Println(er)
				fmt.Println("Connect fail, return message to topic")
				if len(subs) >0{
					subs = removeItemArray(subs, i)
				}
				if len(subs) == 0{
					fmt.Println("3")
					topic.PublishToTopic(Topics[indexTopic], mess.(messqueue.Message))
				}
				break
			}
			topic.PrintTopic(Topics)
		}
	}
	wrk <- w
}

func main() {
	fmt.Println("Server listion at port 8080");
	ln, err := net.Listen("tcp", ":8080")
	if err != nil {
		// handle error
	}

	sbs = make(chan sub, 10)
	wrk = make(chan int, 10)

	for id := 1;id < 10;id ++{
		wrk <- id
	}
	go func() {
		for {
			w := <-wrk
			Subscribe(w, "1")
		}
	}()

	for {
		conn, err := ln.Accept() // this blocks until connection or error
		if err != nil {
			// handle error
			continue
		}
		go HandleConnection(conn) // a goroutine handles conn so that the loop can accept other connections
	}

}