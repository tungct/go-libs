package topic

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

// print all info topic
func PrintTopic(listTopic []Topic) {
	fmt.Println("List Topics : ")
	for _, tp := range listTopic{
		InfoTopic(tp)
	}
	fmt.Println("---------------------")
}

// create new topic
func InitTopic(name string, len int) Topic{
	var topic Topic
	topic.Name = name
	topic.MessQueue = messqueue.InitQueue(len)
	return topic
}

// print info of a topic
func InfoTopic(topic Topic){
	fmt.Println("Name topic : ", topic.Name)
	fmt.Println("Lenght messQueue of topic : ", len(topic.MessQueue))
}

// get index of a topic in topics array
func GetIndexTopic(name string, listTopic []Topic) int {
	for i, tp := range listTopic {
		if tp.Name == name {
			return i
		}
	}
	// if topicName not in topics, return index -1
	return -1
}

// publish a message to a topic
func PublishToTopic(topic Topic, message messqueue.Message){
	messqueue.PutMessageToTopic(message, topic.MessQueue, topic.Name)
}

// subscribe a topic, return a message client need
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