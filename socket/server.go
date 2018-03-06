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

var Topics [] topic.Topic
var mutex = &sync.Mutex{}
var topicNames chan string
var subscribers map[string] []sub

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
			mutex.Lock()
			var subs [] sub
			subscribers[topicName] = subs
			mutex.Unlock()
			tp := topic.InitTopic(topicName, messqueue.LenTopic)
			Topics = append(Topics, tp)
			fmt.Println(len(topicNames))
			topicNames <- topicName

		}
		encoder := gob.NewEncoder(conn)
		mutex.Lock()
		sub := sub{encoder, dec}
		subscribers[topicName] = append(subscribers[topicName], sub)
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
	mutex.Lock()
	allSub := subscribers[topicName]
	mutex.Unlock()

	if len(allSub) >0 && len(Topics[indexTopic].MessQueue) > 0{
		mess := <-Topics[indexTopic].MessQueue
		for i := 0;i<len(allSub);i++{
			er := allSub[i].encode.Encode(mess)
			if er != nil {
				fmt.Println("Connect fail, return message to topic")
				if len(allSub) == 0{
					topic.PublishToTopic(Topics[indexTopic], mess.(messqueue.Message))
				}
				if len(allSub) >0{
					allSub = removeItemArray(allSub, i)
					i = i-1
				}

				continue
			}
			er = allSub[i].decode.Decode(message)
			if er != nil {
				fmt.Println("Connect fail, return message to topic")

				if len(allSub) >0{
					allSub = removeItemArray(allSub, i)
					i = i-1
				}
				if len(allSub) == 0{
					topic.PublishToTopic(Topics[indexTopic], mess.(messqueue.Message))
				}
				continue
			}
			topic.PrintTopic(Topics)
		}
	}
	mutex.Lock()
	subscribers[topicName] = allSub
	mutex.Unlock()
	topicNames <- topicName
}

func main() {
	fmt.Println("Server listion at port 8080");
	ln, err := net.Listen("tcp", ":8080")
	if err != nil {
		// handle error
	}

	topicNames = make(chan string, 10)
	subscribers = make(map[string] []sub, 10)

	go func() {
		for {
			topicName := <- topicNames
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