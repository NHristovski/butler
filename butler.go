package main

import (
	"butler/config"
	"butler/config-parser"
	"butler/httputils"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"strconv"
	"strings"
)

func check(e error) {
	if e != nil {
		log.Panic(e)
	}
}

func logIfDebug(message string, debug bool) {
	if debug {
		log.Println(message)
	}
}

func HandleHTTPConnection(connection net.Conn, requestId string, cfg config.Config) {
	log.Println("| " + requestId + " | | INFO | A connection is established.")

	firstLine, rest, err := httputils.Read(connection, requestId, cfg.Logging.Debug.Enabled)
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

	logIfDebug("| "+requestId+" | | DEBUG | Full path: "+(cfg.Server.RootDir+filePath), cfg.Logging.Debug.Enabled)

	content, err := ioutil.ReadFile(cfg.Server.RootDir + filePath)

	if err != nil {
		handleNotFound(connection, requestId, cfg, filePath, err)
		return
	}

	_, err = httputils.Write(connection, "HTTP/1.1 200 OK\n", requestId, cfg.Logging.Debug.Enabled)
	check(err)
	_, err = httputils.Write(connection, fmt.Sprintf("Content-length: %d \n", len(content)), requestId, cfg.Logging.Debug.Enabled)
	check(err)

	if strings.HasSuffix(filePath, "ico") {
		logIfDebug("| "+requestId+" | | DEBUG | "+filePath+" ends with ico", cfg.Logging.Debug.Enabled)
		_, err = httputils.Write(connection, "Content-Type: image/webp\n", requestId, cfg.Logging.Debug.Enabled)
		check(err)
	} else if strings.HasSuffix(filePath, "css") {
		logIfDebug("| "+requestId+" | | DEBUG | "+filePath+" ends with css", cfg.Logging.Debug.Enabled)
		_, err = httputils.Write(connection, "Content-Type: text/css\n", requestId, cfg.Logging.Debug.Enabled)
		check(err)
	} else if strings.HasSuffix(filePath, "js") {
		logIfDebug("| "+requestId+" | | DEBUG | "+filePath+" ends with js", cfg.Logging.Debug.Enabled)
		_, err = httputils.Write(connection, "Content-Type: text/javascript\n", requestId, cfg.Logging.Debug.Enabled)
		check(err)
	} else if strings.HasSuffix(filePath, "html") {
		logIfDebug("| "+requestId+" | | DEBUG | "+filePath+" ends with html", cfg.Logging.Debug.Enabled)
		_, err = httputils.Write(connection, "Content-Type: text/html\n", requestId, cfg.Logging.Debug.Enabled)
		check(err)
	} else {
		logIfDebug("| "+requestId+" | | DEBUG | "+filePath+" Setting CONTENT TYPE to */*", cfg.Logging.Debug.Enabled)
		_, err = httputils.Write(connection, "Content-Type: */*\n", requestId, cfg.Logging.Debug.Enabled)
		check(err)
	}
	_, err = httputils.Write(connection, "\n", requestId, cfg.Logging.Debug.Enabled)
	check(err)
	_, err = httputils.Write(connection, string(content), requestId, cfg.Logging.Debug.Enabled)
	check(err)

	log.Println("| " + requestId + " | | INFO | Write done")

	check(connection.Close())

	log.Println("| " + requestId + " | | INFO | Connection closed")
}

func handleNotFound(connection net.Conn, requestId string, cfg config.Config, filePath string, err error) {
	content := "<html><b>THE FILE IS NOT PRESENT ON THE SERVER</b></html>"

	log.Printf("| "+requestId+" | | ERROR | Failed to open file %s . Message: %s", filePath, err.Error())
	_, err = httputils.Write(connection, "HTTP/1.1 404 NOT FOUND\n", requestId, cfg.Logging.Debug.Enabled)
	check(err)
	_, err = httputils.Write(connection, fmt.Sprintf("Content-length: %d \n", len(content)), requestId, cfg.Logging.Debug.Enabled)
	check(err)
	_, err = httputils.Write(connection, "Content-Type: text/html\n", requestId, cfg.Logging.Debug.Enabled)
	check(err)
	_, err = httputils.Write(connection, "\n", requestId, cfg.Logging.Debug.Enabled)
	check(err)
	_, err = httputils.Write(connection, content, requestId, cfg.Logging.Debug.Enabled)
	check(err)

	check(connection.Close())
	log.Println("| " + requestId + " | | INFO | Connection closed")
}

func main() {

	cfg := config.Config{}
	err := config_parser.InitConfig(&cfg, "butler.yml")

	if err != nil {
		log.Println("Failed to load configuration")
		log.Fatal(err)
	}

	socket, err := net.Listen(cfg.Server.NetworkType, cfg.Server.Port)

	if err != nil {
		log.Println("The web server cannot start listening on" +
			"port " + cfg.Server.Port + " for " + cfg.Server.NetworkType + " network traffic")
		log.Fatal(err)
	}

	log.Println("| NO_REQ | | INFO | The web server started successfully on port " + cfg.Server.Port + " for " + cfg.Server.NetworkType + " network traffic")
	log.Println("| NO_REQ | | INFO | The root directory of the web server is " + cfg.Server.RootDir)
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
		go HandleHTTPConnection(connection, id, cfg)
	}
}
