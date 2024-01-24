package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"github.com/nicholas-sokolov/test_task/config"
	"github.com/nicholas-sokolov/test_task/message"
	"log"
	"net"
	"time"
)

func main() {
	config.LoadEnv("../.env")

	address := config.GetAddress()

	tcpServer, err := net.ResolveTCPAddr("tcp", address)
	if err != nil {
		log.Fatalln("error to resolve address:", err)
	}

	conn, err := net.DialTCP("tcp", nil, tcpServer)
	if err != nil {
		log.Fatalln("unable to connect to server:", err)
	}
	defer conn.Close()

	for {
		err = handleConnection(conn)
		if err != nil {
			log.Fatalln("error to handle connection:", err)
		}
		time.Sleep(10 * time.Second)
	}

}

func handleConnection(conn *net.TCPConn) error {
	data, err := message.NewRequestChallenge()
	if err != nil {
		return fmt.Errorf("failed to get challenge message: %v", err)
	}

	_, err = conn.Write(data)
	if err != nil {
		return fmt.Errorf("error to send a request: %v", err)
	}

	r := bufio.NewReader(conn)

	for i := 1; i <= 2; i++ {
		payload, err := r.ReadString('\n')
		if err != nil {
			return fmt.Errorf("error to read data: %v", err)
		}

		log.Println("message received:", payload)

		var m message.Message

		err = json.Unmarshal([]byte(payload), &m)
		if err != nil {
			return fmt.Errorf("failed to unmarshal data: %v", err)
		}

		data, err = m.ProcessClientMessage()
		if err != nil {
			return fmt.Errorf("error to process message: %v", err)
		}

		n, err := conn.Write(data)
		if err != nil {
			return fmt.Errorf("error to response hashcash to server: %v", err)
		}

		log.Printf("wrote %d data to server", n)

	}

	return nil

}
