package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"log"

	"github.com/nicholas-sokolov/test_task/config"
	"github.com/nicholas-sokolov/test_task/message"
	"net"
)

func main() {
	config.LoadEnv("../.env")

	address := config.GetAddress()

	l, err := net.Listen("tcp", address)
	if err != nil {
		log.Fatalln("failed to run server:", err)
	}
	defer l.Close()

	fmt.Println("server is listening", address)

	for {
		conn, err := l.Accept()
		if err != nil {
			log.Println("failed to accept connection:", err)
			continue
		}

		go handleConn(conn)
	}
}

func handleConn(conn net.Conn) {
	defer conn.Close()

	addr := conn.RemoteAddr()
	log.Println("Serving connection:", addr)

	r := bufio.NewReader(conn)

	for {
		msg, err := r.ReadBytes('\n')
		if err != nil {
			if err != io.EOF {
				log.Println("error to read data:", err)
			}
			break
		}

		msg, err = getResponse(msg, addr.String())
		if err != nil {
			log.Println("error to get response", err)
			break
		}

		_, err = conn.Write(msg)
		if err != nil {
			log.Println("failed to response data to connection", err)
			break
		}
	}

	log.Println("closing connection:", addr)

}

func getResponse(data []byte, address string) ([]byte, error) {
	var msg message.Message

	err := json.Unmarshal(data, &msg)
	if err != nil {
		return nil, fmt.Errorf("error to unmarshal message: %v", err)
	}

	payload, err := msg.ProcessServerMessage(address)
	if err != nil {
		return nil, err
	}

	return payload, nil
}
