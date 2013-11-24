/**
 * User: Jackong
 * Date: 13-11-24
 * Time: 下午5:13
 */
package global

import (
	"net"
	"io"
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
		request := handleRead(conn, make([]byte, HEADER_LENGTH), HandleHeader)
		response := handleRead(conn, request, handleRequest)
		header := GetHeader(response)
		_, err := conn.Write(append(header, response...))
		if err != nil {
			Log.Alert(err)
			return
		}
	}
}

func handleRead(conn net.Conn, buf []byte, handlerFunc func([]byte)([]byte, error)) []byte {
	_, err := io.ReadFull(conn, buf)
	if err != nil {
		Log.Alert(err)
		return nil
	}
	result, err := handlerFunc(buf)
	if err != nil {
		Log.Alert(err)
		return nil
	}
	return result
}
