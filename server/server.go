package main

import (
	//"github.com/gorilla/mux"
	//"log"
	"net/http"
	"github.com/tungct/go-libs/messqueue"
	"encoding/json"
	"fmt"
	"github.com/tungct/go-libs/workerpool"
)

func RecMessage(rw http.ResponseWriter, request *http.Request) {
	decoder := json.NewDecoder(request.Body)

	var message messqueue.Message
	err := decoder.Decode(&message)

	if err != nil {
		panic(err)
	}
	messqueue.PutMessage(message)
}


func main() {
	messqueue.Queue = messqueue.InitQueue(messqueue.MaxLenQueue)
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
