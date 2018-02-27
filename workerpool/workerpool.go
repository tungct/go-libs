package workerpool

import (
	"github.com/tungct/go-libs/messqueue"
	"github.com/tungct/go-libs/rule_engine"
	"fmt"
	"time"
)

// lenght worker in pool
var MaxLenWorker int = 1
// Worker pool
var Worker chan(int)

// call a worker in workerpool to execute a message by rule
func CallWorker(idWoker int){
	message := <-messqueue.Queue
	t1:= time.Now().UnixNano()
	rule_engine.RuleSys(idWoker, message.(messqueue.Message))

	//return worker to pool
	Worker <- idWoker
	t2 := time.Now().UnixNano()
	t := t2-t1
	fmt.Println(t / 1000)
}

//func WriteToDisk(id int) bool{
//	message := <-go_messqueue.Queue
//	fmt.Println("Worker ", id, "execute Message")
//
//	// check exits file output
//	if _, err := os.Stat("output.json"); err == nil {
//		f, _ := os.OpenFile("output.json", os.O_APPEND|os.O_WRONLY, 0600)
//		defer f.Close()
//		rs, _ := json.Marshal(message)
//		_, er := f.Write(rs)
//		if er != nil {
//			panic(er)
//		}
//		return true
//	}else {
//		jsonData, _  := json.Marshal(message)
//		ioutil.WriteFile("output.json", jsonData, 0600)
//		return true
//	}
//	return false
//}