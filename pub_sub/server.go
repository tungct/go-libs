package main

import (
	"fmt"
	"net"
	"encoding/gob"
	"github.com/tungct/go-libs/messqueue"
	"github.com/tungct/go-libs/topic"
	"sync"

)

// struct of one connect message with a subscriber
type sub struct {
	encode *gob.Encoder // write message to subscriber
	decode *gob.Decoder // read message to subscriber
}

// all topics in server
var Topics [] topic.Topic

// lock a data in concurrency
var mutex = &sync.Mutex{}

// all name of topics
var topicNames chan string

// all subscriber in one topic
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

				// send message success to client
				conn.Write([]byte("Success"))

				topic.PrintTopic(Topics)
			}
			defer conn.Close()
		}
		// status 3 : Subscribe message
	}else if mess.Status == 3{

		// get name topic want to subscribe
		topicName := mess.Content

		// get index topic in topics array
		indexTopic := topic.GetIndexTopic(topicName, Topics)

		// if this topic not exits, init new topic
		if indexTopic == -1 {
			mutex.Lock()
			var subs [] sub
			subscribers[topicName] = subs
			mutex.Unlock()
			tp := topic.InitTopic(topicName, messqueue.LenTopic)
			Topics = append(Topics, tp)
			topicNames <- topicName
		}

		// send success message to client
		encoder := gob.NewEncoder(conn)
		mutex.Lock()
		sub := sub{encoder, dec}
		subscribers[topicName] = append(subscribers[topicName], sub)
		mutex.Unlock()
	}
}

// delete one item in array
func removeItemArray(s []sub, index int) []sub {
	return append(s[:index], s[index+1:]...)
}

// worker to subscribe message to all subscriber this topic
func Subscribe(topicName string){
	message := &messqueue.Message{}
	indexTopic := topic.GetIndexTopic(topicName, Topics)

	mutex.Lock()
	// get subscribers in this topic
	allSub := subscribers[topicName]
	mutex.Unlock()

	if len(allSub) >0 && len(Topics[indexTopic].MessQueue) > 0{
		mess := <-Topics[indexTopic].MessQueue
		for i := 0;i<len(allSub);i++{
			fmt.Println("Subscribe message ", mess)
			er := allSub[i].encode.Encode(mess)

			// if send message to client fail
			if er != nil {
				fmt.Println("Connect fail, return message to topic")

				// if only 1 subscriber of this topic
				if len(allSub) == 0{
					topic.PublishToTopic(Topics[indexTopic], mess.(messqueue.Message))
				}
				if len(allSub) >0{
					// delete conn of client
					allSub = removeItemArray(allSub, i)
					i = i-1
				}
				continue
			}

			er = allSub[i].decode.Decode(message)
			// if receive fail message return from client
			if er != nil {
				fmt.Println("Connect fail, return message to topic")

				if len(allSub) >0{
					// delete conn of client
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

	topicNames = make(chan string)
	subscribers = make(map[string] []sub)

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