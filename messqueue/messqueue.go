package messqueue

import (
	"fmt"
)

type Message struct {
	Status int // status field to check rule
	Content   string
}

const InitConnectStatus = 1
const PublishStatus = 2
const SubscribeStatus = 3
const NilMessageStatus = -1

var MaxLenQueue int = 600

// Message Queue
type MessQueue chan(interface{})
var Queue MessQueue

func CreateMessage(status int, content string) Message{
	var message Message
	message.Status = status
	message.Content = content
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
