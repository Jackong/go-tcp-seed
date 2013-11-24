/**
 * User: Jackong
 * Date: 13-11-24
 * Time: 下午5:13
 */
package global

import (
	"net"
	"io"
	"fmt"
)

func SetUp() {
	listener, err := net.Listen("tcp", Project.String("server", "addr"))
	if err != nil {
		Log.Fatal(err)
	}

	for {
		conn, err := listener.Accept()
		if err != nil {
			Log.Alert(err)
			continue
		}
		go handleClient(conn)
	}

	ShutDown()
}

func handleClient(conn net.Conn) {
	defer conn.Close()

	for {
		fmt.Println("read")
		var header [14]byte
		_, err := io.ReadFull(conn, header[0:])
		if err != nil {
			Log.Alert(err)
			return
		}
		request, err := HandleHeader(header[0:])
		if err != nil {
			Log.Alert(err)
			return
		}
		fmt.Println("read2")

		_, err = io.ReadFull(conn, request[0:])
		if err != nil {
			Log.Alert(err)
			return
		}
		resp, err := handleRequest(request)
		if err != nil {
			Log.Alert(err)
			return
		}
		fmt.Println("write")
		_, err = conn.Write(resp[0:])
		if err != nil {
			Log.Alert(err)
			return
		}
	}
}
