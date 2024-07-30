package main

import (
	"bufio"
	"fmt"
	"net"
	"net/http"
	"os"
	"strings"
)

func handleConnection(conn net.Conn) {
	defer func(conn net.Conn) {
		err := conn.Close()
		if err != nil {
			fmt.Println("Error closing connection:", err)
			return
		}
	}(conn)

	for {
		request, err := http.ReadRequest(bufio.NewReader(conn))
		if err != nil {
			fmt.Println("Error reading request:", err)
			return
		}

		path := request.URL.Path

		response := ""

		if path == "/" {
			response = "HTTP/1.1 200 OK\r\n\r\n"
		} else if parts := strings.SplitN(path, "/", 3); len(parts) >= 2 {
			if parts[1] == "echo" {
				response = "HTTP/1.1 200 OK\r\nContent-Type: text/plain\r\nContent-Length: " + fmt.Sprint(len(parts[2])) + "\r\n\r\n" + parts[2]
			} else if parts[1] == "user-agent" {
				response = "HTTP/1.1 200 OK\r\nContent-Type: text/plain\r\nContent-Length: " + fmt.Sprint(len(request.UserAgent())) + "\r\n\r\n" + request.UserAgent()
			} else {
				response = "HTTP/1.1 404 Not Found\r\n\r\n"
			}
		}

		_, err = conn.Write([]byte(response))

		if err != nil {
			fmt.Println("Error writing to connection:", err)
		}
	}
}

func main() {
	fmt.Println("Logs from your program will appear here!")

	l, err := net.Listen("tcp", "0.0.0.0:4221")
	if err != nil {
		fmt.Println("Failed to bind to port 4221")
		os.Exit(1)
	}

	for {
		conn, err := l.Accept()
		if err != nil {
			fmt.Println("Error accepting connection:", err)
			os.Exit(1)
		}
		// concurrent connections
		go handleConnection(conn)
	}
}
