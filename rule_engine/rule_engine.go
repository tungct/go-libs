package rule_engine

import (
	"github.com/tungct/go-libs/messqueue"
	"fmt"
	"os"
	"encoding/json"
	"io/ioutil"
	"strconv"
	"time"
)

// Rule to classifer message by status
func RuleSys(idWoker int, message messqueue.Message){
	if message.Status == 1{
		fmt.Println("1")
		WriteAppendFile(idWoker, message)
	}else if message.Status == 2{
		fmt.Println("2")
		WriteNewFile(idWoker, message)
	}else {
		fmt.Println("3")
	}
}

// write message to file by append method
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

// write message to new json file
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

