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

	var message go_messqueue.Message
	err := decoder.Decode(&message)

	if err != nil {
		panic(err)
	}
	go_messqueue.PutMessage(message)
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
	fmt.Println("Server run at port 8000")
	http.HandleFunc("/message", RecMessage)
	http.ListenAndServe(":8000", nil)
}
