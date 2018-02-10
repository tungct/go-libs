package messqueue

import (
	"fmt"
	"strconv"
)

type Message struct {
	Status int // status field to check rule
	Content   string
}

var MaxLenQueue int = 600
// Message Queue
type MessQueue chan(interface{})
var Queue MessQueue

func InitMessage() Message{
	var message Message
	message.Status = 1
	message.Content = "Init"
	return message
}

func PublishMessage() Message{
	var message Message
	message.Status = 2
	message.Content = "Message"
	return message
}

func SubscribeMessage(topicName int) Message{
	var message Message
	message.Status = 3
	message.Content = strconv.Itoa(topicName)
	return message
}

func InitQueue(len int)(msQ chan interface{}){
	MessQueue := make(chan interface{}, len)
	return MessQueue
}

// Recv message, push to message queue
func PutMessage(message Message) {
	if len(Queue) < MaxLenQueue{
		Queue <- message
		fmt.Println("Lenght Queue : ", len(Queue))
	}else {
		fmt.Println("Full Queue")
	}
}

func PutMessageToTopic(message Message, queue MessQueue){
	if len(queue) < 10{
		queue <- message
		fmt.Println("Lenght Queue : ", len(queue))
	}else {
		fmt.Println("Full Queue")
	}
}