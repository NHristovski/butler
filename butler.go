package main

import (
	"bufio"
	"butler/config"
	"butler/config-parser"
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net"
	"strconv"
	"strings"
)

func Read(conn net.Conn, requestId string) (string, string, error) {
	reader := bufio.NewReader(conn)
	var buffer bytes.Buffer
	new_line_bytes := []byte("\n")

	counter := 0
	var firstLine string

	for {
		bytes, isPrefix, err := reader.ReadLine()
		if err != nil {
			if err == io.EOF {
				log.Println("| " + requestId + " | | INFO | EOF")
				break
			}
			log.Println("| " + requestId + " | | ERROR | There was an error")
			return "", "", err
		}
		if len(bytes) == 0 && (isPrefix || !isPrefix) {
			break
		}
		log.Printf("| "+requestId+" | | INFO | CONN_IN <<< \"%s\"", string(bytes))

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

func Write(conn net.Conn, content string, requestId string) (int, error) {
	writer := bufio.NewWriter(conn)
	if len(content) < 100 {
		log.Printf("| "+requestId+" | | INFO | CONN_OUT >>> %s", content)
	}
	number, err := writer.WriteString(content)
	if err == nil {
		err = writer.Flush()
	}
	return number, err
}

func check(e error) {
	if e != nil {
		log.Panic(e)
	}
}

func HandleHTTPConnection(connection net.Conn, requestId string) {
	log.Println("| " + requestId + " | | INFO | A connection is established.")

	firstLine, rest, err := Read(connection, requestId)
	// JUST SO THAT I DON'T FORGET TO USE THE REST OF THE REQUEST
	if false {
		log.Println("The rest of the request: " + rest)
	}

	if firstLine == "" {
		log.Println("| " + requestId + " | | ERROR | Empty HTTP Request")
		return
	}

	if err != nil {
		log.Panic(err)
	}

	parts := strings.Split(firstLine, " ")

	filePath := parts[1]

	log.Println("| " + requestId + " | | INFO | Requested file: " + filePath)

	content, err := ioutil.ReadFile("./" + filePath)

	if err != nil {
		content := "<html><b>THE FILE IS NOT PRESENT ON THE SERVER</b></html>"

		log.Printf("| "+requestId+" | | ERROR | Failed to open file %s . Message: %s", filePath, err.Error())
		Write(connection, "HTTP/1.1 404 NOT FOUND\n", requestId)
		Write(connection, fmt.Sprintf("Content-length: %d \n", len(content)), requestId)
		Write(connection, "Content-Type: text/html\n", requestId)
		Write(connection, "\n", requestId)
		Write(connection, content, requestId)

		connection.Close()
	}

	Write(connection, "HTTP/1.1 200 OK\n", requestId)
	Write(connection, fmt.Sprintf("Content-length: %d \n", len(content)), requestId)

	//log.Println("| " + requestId + " | | INFO | " + filePath + " ends with ico")
	//
	//favicon, err := ioutil.ReadFile("favicon.ico")
	//
	//check(err)
	//
	//_, err = Write(connection, "HTTP/1.1 200 OK\n", requestId)
	//check(err)
	//_, err = Write(connection, fmt.Sprintf("Content-length: %d \n", len(favicon)), requestId)
	//check(err)
	//_, err = Write(connection, "Content-Type: image/webp\n", requestId)
	//check(err)
	//_, err = Write(connection, "\n", requestId)
	//check(err)
	//_, err = Write(connection, string(favicon), requestId)
	//check(err)
	//
	//err = connection.Close()
	//check(err)
	if strings.HasSuffix(filePath, "ico") {
		log.Println("| " + requestId + " | | INFO | " + filePath + " ends with ico")
		_, err = Write(connection, "Content-Type: image/webp\n", requestId)
		check(err)
	} else if strings.HasSuffix(filePath, "css") {
		log.Println("| " + requestId + " | | INFO | " + filePath + " ends with css")
		Write(connection, "Content-Type: text/css\n", requestId)
	} else if strings.HasSuffix(filePath, "js") {
		log.Println("| " + requestId + " | | INFO | " + filePath + " ends with js")
		Write(connection, "Content-Type: text/javascript\n", requestId)
	} else if strings.HasSuffix(filePath, "html") {
		log.Println("| " + requestId + " | | INFO | " + filePath + " ends with html")
		Write(connection, "Content-Type: text/html\n", requestId)
	} else {
		log.Println("| " + requestId + " | | INFO | " + filePath + " Setting CONTENT TYPE to */*")
		Write(connection, "Content-Type: */*\n", requestId)
	}
	Write(connection, "\n", requestId)
	Write(connection, string(content), requestId)

	log.Println("| " + requestId + " | | INFO | Write done")
	connection.Close()
	log.Println("| " + requestId + " | | INFO | Connection closed")
}

func main() {

	cfg := config.Config{}
	err := config_parser.InitConfig(&cfg, "butler.yml")

	socket, err := net.Listen(cfg.Server.NetworkType, cfg.Server.Port)

	if err != nil {
		log.Println("The web server cannot start listening on" +
			"port " + cfg.Server.Port + " for " + cfg.Server.NetworkType + " network traffic")
		log.Fatal(err)
	}

	log.Println("| NO_REQ | | INFO | The web server started successfully on port " + cfg.Server.Port + " for " + cfg.Server.NetworkType + " network traffic")
	log.Println("| NO_REQ | | INFO | The web server is waiting for a connection")

	var requestNumber uint64 = 0
	for {
		connection, err := socket.Accept()

		if err != nil {
			log.Printf("| ERROR | Failed to accept connection from the socket\n")
			log.Printf("| ERROR | %s\n", err)
		}

		requestNumber++
		id := strconv.FormatUint(requestNumber, 10)

		// starts goroutine for handling http connection
		go HandleHTTPConnection(connection, id)
	}
}
