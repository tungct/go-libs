package go_workerpool

import (
	"fmt"
	"os"
	"encoding/json"
	"io/ioutil"
	"github.com/tungct/go-libs/messqueue"
)

var MaxLenWorker int = 10
// Worker pool
var Worker chan(int)

func WriteToDisk(id int) bool{
	message := <-go_messqueue.Queue
	fmt.Println("Worker ", id, "execute Message")

	// check exits file output
	if _, err := os.Stat("output.json"); err == nil {
		f, _ := os.OpenFile("output.json", os.O_APPEND|os.O_WRONLY, 0600)
		defer f.Close()
		rs, _ := json.Marshal(message)
		_, er := f.Write(rs)
		if er != nil {
			panic(er)
		}
		Worker <-id
		return true
	}else {
		jsonData, _  := json.Marshal(message)
		ioutil.WriteFile("output.json", jsonData, 0600)
		Worker <-id
		return true
	}
	Worker <-id
	return false
}