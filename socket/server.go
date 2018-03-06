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
var topicNames chan string

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
			topicNames <- topicName
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
				fmt.Println(mess)

				// send message success to client
				conn.Write([]byte("Success"))

				topic.PrintTopic(Topics)
			}
			defer conn.Close()
		}
		// status 3 : Subscribe message
	}else if mess.Status == 3{
		topicName := mess.Content
		indexTopic := topic.GetIndexTopic(topicName, Topics)
		if indexTopic == -1 {
			fmt.Println("22222222")
			tp := topic.InitTopic(topicName, messqueue.LenTopic)
			Topics = append(Topics, tp)
			fmt.Println("333333")
			fmt.Println(len(topicNames))
			topicNames <- topicName
			fmt.Println("1111")

		}
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

func Subscribe(topicName string){
	message := &messqueue.Message{}
	indexTopic := topic.GetIndexTopic(topicName, Topics)

	if len(subs) >0 {
		mess := <-Topics[indexTopic].MessQueue
		for i := 0;i<len(subs);i++{
			er := subs[i].encode.Encode(mess)
			if er != nil {
				fmt.Println(er)
				fmt.Println("Connect fail, return message to topic")
				if len(subs) == 0{
					topic.PublishToTopic(Topics[indexTopic], mess.(messqueue.Message))
				}
				if len(subs) >0{
					subs = removeItemArray(subs, i)
					i = i-1
				}

				continue
			}
			er = subs[i].decode.Decode(message)
			if er != nil {
				fmt.Println(er)
				fmt.Println("Connect fail, return message to topic")

				if len(subs) >0{
					subs = removeItemArray(subs, i)
					i = i-1
				}
				if len(subs) == 0{
					topic.PublishToTopic(Topics[indexTopic], mess.(messqueue.Message))
				}
				continue
			}
			topic.PrintTopic(Topics)
			fmt.Println("end send 1 mess")
		}
	}
	fmt.Println("return topic name to topicNames")
	topicNames <- topicName
}

func main() {
	fmt.Println("Server listion at port 8080");
	ln, err := net.Listen("tcp", ":8080")
	if err != nil {
		// handle error
	}

	sbs = make(chan sub, 10)
	wrk = make(chan int, 10)
	topicNames = make(chan string, 10)

	for id := 0;id < 1;id ++{
		wrk <- id
	}
	go func() {
		for {
			topicName := <- topicNames
			fmt.Println("Sub from topic ", topicName)
			Subscribe(topicName)

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