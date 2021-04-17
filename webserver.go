package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net"
	"strings"
)

func Read(conn net.Conn) (string, string, error) {
	reader := bufio.NewReader(conn)
	var buffer bytes.Buffer
	new_line_bytes := []byte("\n")

	counter := 0
	var firstLine string

	for {
		bytes, isPrefix, err := reader.ReadLine()
		if err != nil {
			log.Println("There was an error")
			if err == io.EOF {
				log.Println("EOF")
				break
			}
			return "", "", err
		}
		if len(bytes) == 0 && (isPrefix || !isPrefix) {
			break
		}
		log.Printf("CONN_IN <<< \"%s\"", string(bytes))

		if counter == 0 {
			firstLine = string(bytes)
		} else {
			buffer.Write(bytes)
			buffer.Write(new_line_bytes)
		}
		counter = counter + 1
	}
	return firstLine, buffer.String(), nil
}

func Write(conn net.Conn, content string) (int, error) {
	writer := bufio.NewWriter(conn)
	if (len(content) < 100) {
		log.Printf("CONN_OUT >>> %s", content)
	}
	number, err := writer.WriteString(content)
	if err == nil {
		err = writer.Flush()
	}
	return number, err
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func HandleHTTPConnection(connection net.Conn) error {
	log.Println("A connection is established")

	firstLine, _, err := Read(connection)

	parts := strings.Split(firstLine, " ")

	filePath := parts[1]

	if err != nil {
		return err
	}

	if filePath == "/favicon.ico" {

		log.Println("Request for the favicon.ico")

		favicon, err := ioutil.ReadFile("C:\\Users\\hp\\dev\\GO\\webserver\\favicon.ico")

		check(err)

		Write(connection, "HTTP/1.1 200 OK\n")
		Write(connection, fmt.Sprintf("Content-length: %d \n", len(favicon)))
		Write(connection, "Content-Type: image/webp\n")
		Write(connection, "\n")
		Write(connection, string(favicon))

		connection.Close()

		return nil
	} else {
		log.Println("Requested file: " + filePath)

		content, err := ioutil.ReadFile("C:\\Users\\hp\\dev\\GO\\webserver\\" + filePath)

		if err != nil {
			content := "<html><b>THE FILE IS NOT PRESENT ON THE SERVER</b></html>"

			log.Printf("| ERROR | Failed to open file %s . Message: %s", filePath, err.Error())
			Write(connection, "HTTP/1.1 404 NOT FOUND\n")
			Write(connection, fmt.Sprintf("Content-length: %d \n", len(content)))
			Write(connection, "Content-Type: text/html\n")
			Write(connection, "\n")
			Write(connection, content)

			connection.Close()

			return nil
		}

		Write(connection, "HTTP/1.1 200 OK\n")
		Write(connection, fmt.Sprintf("Content-length: %d \n", len(content)))
		if strings.HasSuffix(filePath, "css") {
			log.Println("STRING" + filePath + " ends with css")
			Write(connection, "Content-Type: text/css\n")
		}else if strings.HasSuffix(filePath, "js") {
			log.Println("STRING" + filePath + " ends with js")
			Write(connection, "Content-Type: text/javascript\n")
		}else{
			log.Println("STRING" + filePath + " ends with html")
			Write(connection, "Content-Type: text/html\n")
		}
		Write(connection, "\n")
		Write(connection, string(content))

		log.Println("Write done")
		connection.Close()
		log.Println("Connection closed")
		return nil
	}
}

func main() {
	// Listen only for ip v4 over tcp network
	const networkType = "tcp4"
	// Listen on port 4444
	const port = ":4444"

	socket, err := net.Listen(networkType, port)

	if err != nil {
		log.Println("The web server cannot start listening on" +
			"port " + port + " for " + networkType + " network traffic")
		log.Fatal(err)
	}

	log.Println("| INFO | The web server started successfully on port " + port + " for " + networkType + " network traffic")
	log.Println("| INFO | The web server is waiting for a connection")

	for {
		connection, err := socket.Accept()

		if err != nil {
			log.Printf("| ERROR | Failed to accept connection from the socket\n")
			log.Printf("| ERROR | %s\n", err)
		}

		// starts goroutine for handling http connection
		go HandleHTTPConnection(connection)
	}
}
