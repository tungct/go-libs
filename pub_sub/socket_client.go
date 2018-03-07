package main

import (
	"log"
	"net"
	"fmt"
	"time"
)

const (
	message       = "Ping"
	StopCharacter = "\r\n\r\n"
)

func SocketClient(ip string, port int) {
	conn, _ := net.Dial("tcp", "127.0.0.1:3333")

	err := conn.(*net.TCPConn).SetKeepAlive(true)
	if err != nil {
		fmt.Println(err)
		return
	}

	err = conn.(*net.TCPConn).SetKeepAlivePeriod(30 * time.Second)
	if err != nil {
		fmt.Println(err)
		return
	}

	defer conn.Close()

	if err != nil {
		log.Fatalln(err)
	}

	conn.Write([]byte(message))
	conn.Write([]byte(StopCharacter))
	log.Printf("Send: %s", message)

	buff := make([]byte, 1024)
	n, _ := conn.Read(buff)
	log.Printf("Receive: %s", buff[:n])
	conn.Write([]byte(message))
	conn.Write([]byte(StopCharacter))
	log.Printf("Send: %s", message)

	buff = make([]byte, 1024)
	n, _ = conn.Read(buff)
	log.Printf("Receive: %s", buff[:n])

}

func main() {

	var (
		ip   = "127.0.0.1"
		port = 3333
	)

	SocketClient(ip, port)

}