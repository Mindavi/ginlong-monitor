package main

import (
	"io/ioutil"
	"log"
	"net"
	"time"
)

const (
	host      = ""
	port      = "9999"
	conn_type = "tcp"
)

func main() {
	server, err := net.Listen(conn_type, host+":"+port)
	if err != nil {
		log.Fatal(err)
	}
	defer server.Close()

	for {
		conn, err := server.Accept()
		log.Print("accepted connection")
		if err != nil {
			log.Fatal(err)
		}
		go handleRequest(conn)
	}
}

func handleRequest(conn net.Conn) {
	buf := make([]byte, 512)
	requestLen, err := conn.Read(buf)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Received %d bytes", requestLen)
	log.Print(buf[:requestLen])
	err = ioutil.WriteFile(time.Now().Format(time.RFC3339), buf[:requestLen], 0644)
	if err != nil {
		log.Fatal(err)
	}
	conn.Close()
}
