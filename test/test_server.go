package main

import (
	"bufio"
	"io"
	"log"
	"net"
	"strings"
)

func LinijkaServer() {
	l, err := net.Listen("tcp", ":4001")
	if err != nil {
		log.Fatal(err)
	}
	defer l.Close()
	for {
		// Wait for a connection.
		conn, err := l.Accept()
		if err != nil {
			log.Fatal(err)
		}
		// Handle the connection in a new goroutine.
		// The loop then returns to accepting, so that
		// multiple connections may be served concurrently.
		go func(c net.Conn) {
			defer conn.Close()

			buf := make([]byte, 1024)
			r := bufio.NewReader(conn)
			w := bufio.NewWriter(conn)

		CONN:
			for {
			COMMAND:
				for {
					n, err := r.Read(buf)
					data := string(buf[:n])

					switch err {
					case io.EOF:
						log.Println("EOF")
						break CONN
					case nil:
						log.Println("Receive:", data)
						if isTransportOver(data) {
							break COMMAND
						}

					default:
						log.Fatalf("Receive data failed:%s", err)
						return
					}

				}
				w.Write([]byte("<OK!>\r\n"))
				w.Flush()
				log.Printf("Send: %s", "<OK!>\r\n")
			}
			log.Println("Ending connection")
		}(conn)
	}
}

func isTransportOver(data string) (over bool) {
	over = strings.HasSuffix(data, "\r\n")
	return
}

func main() {
	log.Println("Waiting for connection")
	LinijkaServer()
	log.Println("Shutting down")
}
