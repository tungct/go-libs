package go_messqueue

import (
	"fmt"
)

type Message struct {
	Content   string
}

var MaxLenQueue int = 600
// Message Queue
var Queue chan(Message)

// Recv message, push to message queue
func PutMessage(message Message) {
	if len(Queue) < MaxLenQueue{
		Queue <- message
		fmt.Println("Lenght Queue : ", len(Queue))
	}else {
		fmt.Println("Full Queue")
	}
}