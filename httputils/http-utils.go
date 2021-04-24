package httputils

import (
	"bufio"
	"bytes"
	"io"
	"log"
	"net"
)

func Read(conn net.Conn, requestId string, debug bool) (string, string, error) {
	reader := bufio.NewReader(conn)
	var buffer bytes.Buffer
	new_line_bytes := []byte("\n")

	counter := 0
	var firstLine string

	for {
		bytes, isPrefix, err := reader.ReadLine()
		if err != nil {
			if err == io.EOF {
				if debug {
					log.Println("| " + requestId + " | | DEBUG | EOF")
				}
				break
			}
			log.Println("| " + requestId + " | | ERROR | There was an error")
			return "", "", err
		}
		if len(bytes) == 0 && (isPrefix || !isPrefix) {
			break
		}

		if debug {
			log.Printf("| "+requestId+" | | DEBUG | CONN_IN <<< \"%s\"", string(bytes))
		}

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

func Write(conn net.Conn, content string, requestId string, debug bool) (int, error) {
	writer := bufio.NewWriter(conn)
	if debug {
		log.Printf("| "+requestId+" | | DEBUG | CONN_OUT >>> %s", content)
	}
	number, err := writer.WriteString(content)
	if err == nil {
		err = writer.Flush()
	}
	return number, err
}
