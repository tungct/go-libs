# go-msqueue-worker

Server nhận Request message (POST), đẩy message vào MessageQueue chờ xử lý, WorkerPool thực hiện xử lý message trong queue và ghi ra file

## 1. Yêu cầu
- Go 1.9 hoặc thấp hơn, ubuntu 16.04

## 2. Hướng dẫn

Cấu trúc MessageQueue và WorkerPool 

![architecture introduction diagram](image/msq.png)

### 2.1 Định nghĩa MessageQueue
go-messqueue.go:

```
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
```

### 2.2 Định nghĩa WorkerPool

workerpool

```
var MaxLenWorker int = 10
// Worker pool
var Worker chan(int)

func WriteToDisk(id int) bool{
	Worker <-id
	message := <-go_messqueue.Queue
	fmt.Println("Worker ", id, "execute Message")

	// check exits file output
	if _, err := os.Stat("output.json"); err == nil {
		f, _ := os.OpenFile("output.json", os.O_APPEND|os.O_WRONLY, 0600)
		rs, _ := json.Marshal(message)
		if _, err := f.Write(rs); err != nil {
			panic(err)
		}
		return true
	}else {
		jsonData, _  := json.Marshal(message)
		ioutil.WriteFile("output.json", jsonData, 0600)
		return true
	}
	return false
}
```

### 2.3 Server

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
	go_messqueue.Queue = make(chan go_messqueue.Message, go_messqueue.MaxLenQueue)
	go_workerpool.Worker = make(chan int, go_workerpool.MaxLenWorker)

	for id := 0 ; id < go_workerpool.MaxLenWorker ; id ++{
		go_workerpool.Worker <-id
	}

	// Worker execute message in pool, write to disk
	go func() {
		for {
			w := <- go_workerpool.Worker
			go_workerpool.WriteToDisk(w)
		}
	}()


	router := mux.NewRouter()
	router.HandleFunc("/message", RecMessage).Methods("POST")
	log.Fatal(http.ListenAndServe(":8000", router))
}
```

### 2.4 Chạy chương trình

Để chạy chương trình, vào thư mục $GOPATH/src/github.com/tungct/go-libs

```
$ go run /server/server.go

Server run at port 8000
```
Call REST API với phương thức POST

```
curl -X POST -d "{\"content\": \"test\"}" http://localhost:8000/message
```
### 2.5 Test Performance bằng go-wrk

- https://github.com/tsliwowicz/go-wrk

Để test, vào thư mục $GOPATH/bin


```
./go-wrk -M POST -d 5 -body "{\"content\": \"test\"}" http://127.0.0.1:8000/message
```

Kết quả :

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
