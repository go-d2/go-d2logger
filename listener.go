package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
	"strings"
)

const (
	NETWORK = "tcp"
	ADDRESS = "localhost"
	PORT    = "4444"
)

var logfile *os.File
var buffer string

func listen() {
	var err error
	logfile, err = os.OpenFile("dofus.html", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal("> Can't open log file: ", err)
	}

	defer logfile.Close()

	l, err := net.Listen(NETWORK, ADDRESS+":"+PORT)
	if err != nil {
		log.Fatal("> Can't start listener: ", err)
	}

	defer l.Close()

	fmt.Println("> Listening on", ADDRESS+":"+PORT)

	for {
		conn, err := l.Accept()
		if err != nil {
			log.Fatal(err)
		}

		go handle(conn)
	}
}

func handle(conn net.Conn) {
	fmt.Println("> New connection from", conn.LocalAddr().String())
	defer close(conn)

	reader := bufio.NewReader(conn)

	for {
		data, err := reader.ReadString('\n')
		if err != nil {
			return
		}

		buffer += data

		if strings.Contains(data, "]]></showMessage>") {
			logs := buffer
			buffer = ""

			logs = strings.ReplaceAll(logs, "!SOS<showMessage key=\"trace\"><![CDATA[", "<li class=\"trace\">")
			logs = strings.ReplaceAll(logs, "!SOS<showMessage key=\"debug\"><![CDATA[", "<li class=\"debug\">")
			logs = strings.ReplaceAll(logs, "!SOS<showMessage key=\"info\"><![CDATA[", "<li class=\"info\">")
			logs = strings.ReplaceAll(logs, "!SOS<showMessage key=\"warning\"><![CDATA[", "<li class=\"warning\">")
			logs = strings.ReplaceAll(logs, "!SOS<showMessage key=\"error\"><![CDATA[", "<li class=\"error\">")
			logs = strings.ReplaceAll(logs, "!SOS<showMessage key=\"fatal\"><![CDATA[", "<li class=\"fatal\">")
			logs = strings.ReplaceAll(logs, "]]></showMessage>", "</li>")

			logs = strings.ReplaceAll(logs, "!SOS<showFoldMessage key=\"trace\"><title><![CDATA[", "<li class=\"trace\">")
			logs = strings.ReplaceAll(logs, "!SOS<showFoldMessage key=\"debug\"><title><![CDATA[", "<li class=\"debug\">")
			logs = strings.ReplaceAll(logs, "!SOS<showFoldMessage key=\"info\"><title><![CDATA[", "<li class=\"info\">")
			logs = strings.ReplaceAll(logs, "!SOS<showFoldMessage key=\"warning\"><title><![CDATA[", "<li class=\"warning\">")
			logs = strings.ReplaceAll(logs, "!SOS<showFoldMessage key=\"error\"><title><![CDATA[", "<li class=\"error\">")
			logs = strings.ReplaceAll(logs, "!SOS<showFoldMessage key=\"fatal\"><title><![CDATA[", "<li class=\"fatal\">")

			logs = strings.ReplaceAll(logs, "</title><message><![CDATA[", "<ul><li>")
			logs = strings.ReplaceAll(logs, "</message></showFoldMessage>", "</li></ul></li>")

			n, err := logfile.WriteString(logs)
			if err != nil {
				log.Fatal("> Can't write into log file: ", err)
			}

			fmt.Println("> Wrote", n, "bytes to log file")
		}
	}
}

func close(conn net.Conn) {
	fmt.Println("> Client disconnected", conn.LocalAddr().String())
	conn.Close()
}
