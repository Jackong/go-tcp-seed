/**
 * User: jackong
 * Date: 10/17/13
 * Time: 6:24 PM
 */
package main

import (
	"fmt"
	"net"
	. "github.com/Jackong/go-tcp-seed/global"
	"os"
)


func main() {
	listener, err := net.Listen("tcp", Project.String("server", "addr"))
	if err != nil {
		fmt.Println(err)
		os.Exit(2)
	}

	for {
		conn, err := listener.Accept()
		if err != nil {
			continue
		}
		go handleClient(conn)
	}

	ShutDown()
}

func handleClient(conn net.Conn) {
	defer conn.Close()

	var buf [512]byte
	for {
		n, err := conn.Read(buf[0:])
		if err != nil {
			return
		}

		_, err2 := conn.Write(buf[0:n])
		if err2 != nil {
			return
		}
	}
}
