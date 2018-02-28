package main

import (
	"fmt"
	"log"
	"net"
	"encoding/gob"
	"github.com/tungct/go-libs/messqueue"
	"strings"
	"strconv"
	"time"
)

// convert bytes array to string
func BytesToString(data []byte) string {
	return string(data[:])
}

// init connection to server
func InitConn(conn *net.TCPConn, encoder *gob.Encoder) bool{
	// Init connect to server
	mess := messqueue.CreateMessage(messqueue.InitConnectStatus, "Init")
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
func SendMess(ip string, port int){
	addr := strings.Join([]string{ip, strconv.Itoa(port)}, ":")
	tcpAddr, err := net.ResolveTCPAddr("tcp", addr)
	conn, err := net.DialTCP("tcp", nil, tcpAddr)

	if err != nil {
		log.Fatal("Connection error", err)
	}

	// Publish 3 message
	encoder := gob.NewEncoder(conn)

	// init connect to server
	init := InitConn(conn, encoder)

	// if server accept connect, publish message
	if init == true {
		for i := 0; i < 3; i++ {
			mess := messqueue.CreateMessage(messqueue.PublishStatus, strconv.Itoa(i))
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
	fmt.Println("start client");
	SendMess(ip, port)
}
