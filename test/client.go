package main

import (
	"log"
	"net"
	"strconv"
	"strings"
	"github.com/tungct/go-libs/messqueue"
)

const (
	messageInit       = "Hello"
	messagePublish   = "Job"
	StopCharacter = "\r\n\r\n"
)

func InitConnect(ip string, port int){
	addr := strings.Join([]string{ip, strconv.Itoa(port)}, ":")

	conn, err := net.Dial("tcp", addr)
	defer conn.Close()

	if err != nil {
		log.Fatalln(err)
	}
	conn.Write([]byte(messageInit))
	conn.Write([]byte(StopCharacter))
	log.Printf("Send: %s", messageInit)
	buff := make([]byte, 1024)
	n, _ := conn.Read(buff)
	log.Printf("Receive: %s", buff[:n])
}

func SendMessage(message messqueue.Message, ip string, port int){
	addr := strings.Join([]string{ip, strconv.Itoa(port)}, ":")

	conn, err := net.Dial("tcp", addr)
	defer conn.Close()

	if err != nil {
		log.Fatalln(err)
	}
	conn.Write([]byte(message))
	conn.Write([]byte(StopCharacter))
	log.Printf("Send: %s", messageInit)
	buff := make([]byte, 1024)
	n, _ := conn.Read(buff)
	log.Printf("Receive: %s", buff[:n])
}

func SocketClient(ip string, port int) {
	addr := strings.Join([]string{ip, strconv.Itoa(port)}, ":")

	conn, err := net.Dial("tcp", addr)
	defer conn.Close()

	if err != nil {
		log.Fatalln(err)
	}
	conn.Write([]byte(messageInit))
	conn.Write([]byte(StopCharacter))
	log.Printf("Send: %s", messageInit)
	buff := make([]byte, 1024)
	n, _ := conn.Read(buff)
	log.Printf("Receive: %s", buff[:n])

}

func main() {

	var (
		ip   = "127.0.0.1"
		port = 8000
	)
	message := messqueue.InitMessage()

	SendMessage(message, ip, port)

}