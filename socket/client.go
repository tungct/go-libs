package main

import (
	"fmt"
	"log"
	"net"
	"encoding/gob"
	"github.com/tungct/go-libs/messqueue"
	"strings"
	"strconv"
)

func BytesToString(data []byte) string {
	return string(data[:])
}

func InitConn(ip string, port int) bool{
	addr := strings.Join([]string{ip, strconv.Itoa(port)}, ":")
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		log.Fatal("Connection error", err)
		return false
	}
	encoder := gob.NewEncoder(conn)
	mess := messqueue.InitMessage()
	encoder.Encode(mess)
	buff := make([]byte, 1024)
	n, _ := conn.Read(buff)
	log.Printf("Receive: %s", buff[:n])
	if BytesToString(buff[:n]) != "OK"{
		return false
	}
	conn.Close()
	fmt.Println("done");
	return true
}

func SendMess(ip string, port int){
	addr := strings.Join([]string{ip, strconv.Itoa(port)}, ":")
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		log.Fatal("Connection error", err)
	}
	encoder := gob.NewEncoder(conn)
	mess := messqueue.TestMessage()
	encoder.Encode(mess)
	buff := make([]byte, 1024)
	n, _ := conn.Read(buff)
	log.Printf("Receive: %s", buff[:n])
	conn.Close()
	fmt.Println("done");
}

func main() {
	var (
		ip   = "127.0.0.1"
		port = 8080
	)
	fmt.Println("start client");
	succ := InitConn(ip, port)
	if succ == true{
		SendMess(ip, port)
	}
}
