package main

import (
	"fmt"
	"net"
	"encoding/gob"
	"reflect"
)

type P struct {
	M, N int64
}

func handleConnection(conn net.Conn) {
	dec := gob.NewDecoder(conn)
	p := &P{}
	dec.Decode(p)
	fmt.Println(reflect.TypeOf(p).Kind())
	fmt.Printf("Received : %+v", p.M);
	conn.Close()
}

func main() {
	fmt.Println("Server listion at port 8080");
	ln, err := net.Listen("tcp", ":8080")
	if err != nil {
		// handle error
	}
	for {
		conn, err := ln.Accept() // this blocks until connection or error
		if err != nil {
			// handle error
			continue
		}
		go handleConnection(conn) // a goroutine handles conn so that the loop can accept other connections
	}
}