package Topic

import (
	"github.com/tungct/go-libs/messqueue"
	"strings"
	"fmt"
)

type Topic struct {
	Name string
	MessQueue messqueue.MessQueue
}

// classifer message to many topic
func RuleTopic(mess messqueue.Message) string{
	var topicName string
	if strings.Contains(mess.Content, "Message"){
		topicName = "Message"
	}else {
		topicName = "other"
	}
	return topicName
}

func PrintTopic(listTopic []Topic) {
	for _, tp := range listTopic{
		fmt.Println(tp.Name)
	}
}

func InitTopic(name string, len int) Topic{
	var topic Topic
	topic.Name = name
	topic.MessQueue = messqueue.InitQueue(len)
	return topic
}

func GetIndexTopic(name string, listTopic []Topic) int {
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

func Subscribe(topic Topic) (succ bool, ms messqueue.Message){
	var message interface{}
	var er bool = false
	if len(topic.MessQueue )> 0{
		message = <- topic.MessQueue
		er = true
	}
	mess := message.(messqueue.Message)
	return er, mess
}