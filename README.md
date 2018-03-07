# go-libs

- WorkerPool & RuleEngine : Server nhận Request message (POST), đẩy message vào MessageQueue chờ xử lý, WorkerPool thực hiện xử lý message trong queue và tuân theo Rule_Engine để chờ có các action tương ứng
- Publish & Subscribe : (TCP/IP Socket) Server nhận message từ publisher và đẩy vào topic tương ứng, subscriber gửi yêu cầu đến server lấy phần tử trong topic mong muốn

## 1. Yêu cầu
- Go 1.9 hoặc thấp hơn, ubuntu 16.04

## 2. Hướng dẫn

### 2.1 WorkerPool & RuleEngine
Cấu trúc MessageQueue và WorkerPool 

![architecture introduction diagram](image/msq.png)

#### 2.1.1 Định nghĩa MessageQueue
messqueue/messqueue.go:

```
type Message struct {
	Status int // status field to check rule
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

```

#### 2.1.2 Định nghĩa WorkerPool

workerpool/workerpool.go

```
var MaxLenWorker int = 10
// Worker pool
var Worker chan(int)

func CallWorker(idWoker int){
	message := <-messqueue.Queue
	rule_engine.RuleSys(idWoker, message)

	//return worker to pool
	Worker <- idWoker
}

```

#### 2.1.3 Cấu trúc Rule_Engine

rule_engine/rule_engine.go

```
func RuleSys(idWoker int, message messqueue.Message){
	if message.Status == 1{
		fmt.Println("Status 1")
		WriteAppendFile(idWoker, message)
	}else if message.Status == 2{
		fmt.Println("Status 2")
		WriteNewFile(idWoker, message)
	}else {
		fmt.Println("Status not support")
	}
}

func WriteAppendFile(idWorker int, message messqueue.Message) bool{
	fmt.Println("Worker ", idWorker, " execute Message write append file")

	// check exits file output
	if _, err := os.Stat("output.json"); err == nil {
		f, _ := os.OpenFile("output.json", os.O_APPEND|os.O_WRONLY, 0600)
		defer f.Close()
		rs, _ := json.Marshal(message)
		_, er := f.Write(rs)
		if er != nil {
			panic(er)
		}
		return true
	}else {
		jsonData, _  := json.Marshal(message)
		ioutil.WriteFile("output.json", jsonData, 0600)
		return true
	}
	return false
}

func WriteNewFile(idWorker int, message messqueue.Message) bool{
	fmt.Println("Worker ", idWorker ," execute Message write new file")

	// name file = timeNowUnix + .json
	newFile := strconv.FormatInt(time.Now().UnixNano(),10) + ".json"
	jsonData, _  := json.Marshal(message)
	err := ioutil.WriteFile(newFile, jsonData, 0600)
	if err != nil{
		panic(err)
		return false
	}
	return true
}

```

#### 2.1.4 Server

server/server.go

```
func RecMessage(rw http.ResponseWriter, request *http.Request) {
	decoder := json.NewDecoder(request.Body)

	var message go_messqueue.Message
	err := decoder.Decode(&message)

	if err != nil {
		panic(err)
	}
	go_messqueue.PutMessage(message)

	fmt.Println(message.Content)
}


func main() {
	messqueue.Queue = make(chan messqueue.Message, messqueue.MaxLenQueue)
	workerpool.Worker = make(chan int, workerpool.MaxLenWorker)

	for id := 0 ; id < workerpool.MaxLenWorker ; id ++{
		workerpool.Worker <-id
	}

	// Worker execute message in pool, write to disk
	go func() {
		for {
			w := <- workerpool.Worker
			workerpool.CallWorker(w)
		}
	}()
	fmt.Println("Server run at port 8000")
	http.HandleFunc("/message", RecMessage)
	http.ListenAndServe(":8000", nil)
}

```

#### 2.1.5 Chạy chương trình

Để chạy chương trình, vào thư mục $GOPATH/src/github.com/tungct/go-libs

```
$ go run /server/server.go

Server run at port 8000
```
Call REST API với phương thức POST

```
curl -X POST -d '{"status":1,"content":"test"}' http://127.0.0.1:8000/message
```
#### 2.1.6 Test Performance bằng go-wrk

- https://github.com/tsliwowicz/go-wrk

Để test, vào thư mục $GOPATH/bin


```
./go-wrk -M POST -d 5 -body "{\"status\":1,\"content\": \"test\"}" http://127.0.0.1:8000/message
```

Kết quả chạy trên máy tính RAM 4 Gigabyte :

- 1 worker trong workerpool (MaxLenWorker=1)

```
Running 5s test @ http://127.0.0.1:8000/message
  10 goroutine(s) running concurrently
86281 requests in 4.918361348s, 8.15MB read
Requests/sec:		17542.63
Transfer/sec:		1.66MB
Avg Req Time:		570.039µs
Fastest Request:	63.42µs
Slowest Request:	21.826247ms
Number of Errors:	0
```

- 2 worker trong workerpool (MaxLenWorker=2)

```
Running 5s test @ http://127.0.0.1:8000/message
  10 goroutine(s) running concurrently
87194 requests in 4.920370263s, 8.23MB read
Requests/sec:		17721.02
Transfer/sec:		1.67MB
Avg Req Time:		564.301µs
Fastest Request:	61.499µs
Slowest Request:	24.348407ms
Number of Errors:	0
```

- 10 worker trong workerpool (MaxLenWorker=10)

```
Running 5s test @ http://127.0.0.1:8000/message
  10 goroutine(s) running concurrently
87531 requests in 4.914277797s, 8.26MB read
Requests/sec:		17811.57
Transfer/sec:		1.68MB
Avg Req Time:		561.432µs
Fastest Request:	58.769µs
Slowest Request:	18.407424ms
Number of Errors:	0

```

- 20 worker trong workerpool (MaxLenWorker=20)

```
Running 5s test @ http://127.0.0.1:8000/message
  10 goroutine(s) running concurrently
87047 requests in 4.922004624s, 8.22MB read
Requests/sec:		17685.27
Transfer/sec:		1.67MB
Avg Req Time:		565.442µs
Fastest Request:	66.411µs
Slowest Request:	12.811064ms
Number of Errors:	0

```
### 2.2 Publish & Subscribe 
Kiến trúc chung Publish & Subscribe

![architecture introduction diagram](image/publish_subscribe.jpg)

Kiến trúc trong chương trình

![architecture introduction diagram](image/pubsub.png)

#### 2.2.1 Định nghĩa Message 

Message được chuyển qua socket giữa server và client nhằm khởi tạo connect, publish, subscribe hay báo lỗi với trường status tương ứng

messqueue/messqueue.go:

```
const InitConnectStatus = 1
const PublishStatus = 2
const SubscribeStatus = 3
const NilMessageStatus = -1
const LenTopic = 20

var MaxLenQueue int = 600

type Message struct {
	Status int // status field to check rule
	Content   string
}

func CreateMessage(status int, content string) Message{
	var message Message
	message.Status = status
	message.Content = content
	return message
}

func PutMessageToTopic(message Message, queue MessQueue){
	if len(queue) < 10{
		queue <- message
		fmt.Println("Lenght Queue : ", len(queue))
	}else {
		fmt.Println("Full Queue")
	}
}


```

#### 2.2.2 Định nghĩa topic 

Message sẽ được server đưa vào các Topic tương ứng

topic/topic.go:

```
type Topic struct {
	Name string
	MessQueue messqueue.MessQueue
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

```

#### 2.2.3 Server 

pub_sub/server.go

Sử dụng goroutine để xử lý độc lập các connect từ client 

```
...
for {
		conn, err := ln.Accept() // this blocks until connection or error
		if err != nil {
			// handle error
			continue
		}
		go HandleConnection(conn) // a goroutine handles conn so that the loop can accept other connections
	}
...

```

Xử lý các connect từ client, nếu message từ publish client, xác nhận topic client mong muốn, nhận message và đẩy vào topic đó

```
...
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
	}
...

```

Xử lý connect đến từ subscribe client, nhận các thông tin như topic muốn subscribe và xếp subscribe client vào danh sách
các client cùng subscribe vào topic đó

```
...

// struct of one connect message with a subscriber
type sub struct {
	encode *gob.Encoder // write message to subscriber
	decode *gob.Decoder // read message to subscriber
}
...
sub := sub{encoder, dec}
subscribers[topicName] = append(subscribers[topicName], sub)
...

```

Đối với gửi message đến các client subscribe, xử dụng goroutine để làm worker xử lý riêng biệt

```
...

// all name of topics
var topicNames chan string

// all subscriber in one topic
var subscribers map[string] []sub

...

go func() {
		for {
			topicName := <- topicNames
			Subscribe(topicName)
		}

}()
...

```

Thực hiện gửi message trong 1 topic đến subscribe client 

```
// worker to subscribe message to all subscriber this topic
func Subscribe(topicName string){
	...
}

```

Gửi message trong topic đến từng client subscribe vào topic đó

```
// get subscribers in this topic
allSub := subscribers[topicName]

...

mess := <-Topics[indexTopic].MessQueue
		for i := 0;i<len(allSub);i++{
		    ...
		}

```

Trường hợp connect giữa server và client bị ngắt, xóa connect của client trong danh sách các client sub vào topic
và nếu là client duy nhất sub vào topic, khi gửi lỗi cần return message đó vào topic để không bị mất dữ liệu

```
if er != nil {
				fmt.Println("Connect fail, return message to topic")
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

```

Server listen at port 8080

#### 2.2.4 Publisher 

Khởi tạo connect, tên topic muốn publish tới và gửi các message đến server

pub_sub/publish_client.go:

```
// convert bytes array to string
func BytesToString(data []byte) string {
	return string(data[:])
}

// init connection to server
func InitConn(conn *net.TCPConn, encoder *gob.Encoder, topicName string) bool{
	// Init connect to server
	mess := messqueue.CreateMessage(messqueue.InitConnectStatus, topicName)
	encoder.Encode(mess)
	buff := make([]byte, 1024)
	n, _ := conn.Read(buff)
	log.Printf("Receive: %s", buff[:n])
	if BytesToString(buff[:n]) != "OK"{
		return false
	}
	return true
}

// send publish message to server
func SendMess(ip string, port int, topicName string){
	fmt.Println("Publish to topic " + topicName)
	addr := strings.Join([]string{ip, strconv.Itoa(port)}, ":")
	tcpAddr, err := net.ResolveTCPAddr("tcp", addr)
	conn, err := net.DialTCP("tcp", nil, tcpAddr)

	if err != nil {
		log.Fatal("Connection error", err)
	}

	// Publish 3 message
	encoder := gob.NewEncoder(conn)

	// init connect to server
	init := InitConn(conn, encoder, topicName)

	// if server accept connect, publish message
	if init == true {
		for i := 0; i < 10; i++ {
			content := rand.Intn(100)
			mess := messqueue.CreateMessage(messqueue.PublishStatus, strconv.Itoa(content))
			encoder.Encode(mess)
			buff := make([]byte, 1024)
			n, _ := conn.Read(buff)
			log.Printf("Receive: %s", buff[:n])
			time.Sleep(1 * time.Second)
		}
	}
	conn.Close()
	fmt.Println("done");
}



func main() {
	var (
		ip   = "127.0.0.1"
		port = 8080
	)
	topicName := os.Args[1]
	fmt.Println("start client");
	SendMess(ip, port, topicName)
}

```

#### 2.2.5 Subscriber 

Khởi tạo message đến server chứa tên topic cần subscribe, chờ để nhận các message chứa trong topic đó

pub_sub/subscribe_client.go

```
// send subscribe message to server and get message response
func GetMess(ip string, port int, topicName string){
	addr := strings.Join([]string{ip, strconv.Itoa(port)}, ":")
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		log.Fatal("Connection error", err)
	}

	// subscribe topic with topicName
	fmt.Println("Subscribe topic : ", topicName)

	// create subscribe message send to server
	mess := messqueue.CreateMessage(messqueue.SubscribeStatus, topicName)

	// create nil message return after receive from server
	mes := messqueue.CreateMessage(messqueue.NilMessageStatus, topicName)
	// init encoder to send data to server
	encoder := gob.NewEncoder(conn)
	encoder.Encode(mess)

	// init decoder to read data from server response
	dec := gob.NewDecoder(conn)

	for {

		messRes := &messqueue.Message{}
		err = dec.Decode(messRes)
		if err != nil{
			log.Fatal(err)
			break
		}
		err = encoder.Encode(mes)
		if err != nil{
			log.Fatal(err)
			break
		}
		fmt.Println("Received message : ", messRes, "from topic " + topicName)
	}
	defer conn.Close()
}

func main() {
	var (
		ip   = "127.0.0.1"
		port = 8080
	)
	topicName := os.Args[1]
	fmt.Println("subscribe client");
	GetMess(ip, port, topicName)
}

```































