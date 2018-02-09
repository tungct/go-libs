package Topic

import (
	"github.com/tungct/go-libs/messqueue"
)

type Topic struct {
	Name string
	MessQueue messqueue.MessQueue
}

func InitTopic(name string, len int) Topic{
	var topic Topic
	topic.Name = name
	topic.MessQueue = messqueue.InitQueue(len)
	return topic
}

func PublishToTopic(topic Topic, message messqueue.Message){
	messqueue.PutMessage(message)
}

func Subscribe(topic Topic) messqueue.Message{
	message := <- topic.MessQueue
	message = message.(messqueue.Message)
	return message
}