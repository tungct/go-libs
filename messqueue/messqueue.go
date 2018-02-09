package messqueue

import (
	"fmt"
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
	message.Content = "test"
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