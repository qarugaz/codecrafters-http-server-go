package main

import (
	"fmt"
	"net"
	"os"
	"strings"
)

func main() {
	fmt.Println("Logs from your program will appear here!")

	l, err := net.Listen("tcp", "0.0.0.0:4221")
	if err != nil {
		fmt.Println("Failed to bind to port 4221")
		os.Exit(1)
	}
	defer l.Close()

	conn, err := l.Accept()
	if err != nil {
		fmt.Println("Error accepting connection:", err)
		os.Exit(1)
	}
	defer conn.Close()

	buffer := make([]byte, 1024)
	_, err = conn.Read(buffer)
	if err != nil {
		fmt.Println("Error reading request:", err)
		return
	}

	requestLine := strings.Split(string(buffer), "\r\n")[0]
	requestParts := strings.Split(requestLine, " ")

	if len(requestParts) < 2 {
		fmt.Println("Invalid request")
		return
	}

	path := requestParts[1]
	response := ""

	if path == "/" {
		response = "HTTP/1.1 200 OK\r\n\r\n"
	} else if parts := strings.SplitN(path, "/", 3); len(parts) >= 2 && parts[1] == "echo" {
		response = "HTTP/1.1 200 OK\r\nContent-Type: text/plain\r\nContent-Length: " + fmt.Sprint(len(parts[2])) + "\r\n\r\n" + parts[2]
	} else {
		response = "HTTP/1.1 404 Not Found\r\n\r\n"
	}
	_, err = conn.Write([]byte(response))

	if err != nil {
		fmt.Println("Error writing to connection:", err)
	}
}
