package Topic

import (
	"github.com/tungct/go-libs/messqueue"
)

type Topic struct {
	Name int
	MessQueue messqueue.MessQueue
}

func InitTopic(name int, len int) Topic{
	var topic Topic
	topic.Name = name
	topic.MessQueue = messqueue.InitQueue(len)
	return topic
}

func GetIndexTopic(name int, listTopic []Topic) int {
	for i, tp := range listTopic {
		if tp.Name == name {
			return i
		}
	}
	return -1
}

func PublishToTopic(topic Topic, message messqueue.Message){
	messqueue.PutMessageToTopic(message, topic.MessQueue)
}

func Subscribe(topic Topic) messqueue.Message{
	message := <- topic.MessQueue
	mess := message.(messqueue.Message)
	return mess
}