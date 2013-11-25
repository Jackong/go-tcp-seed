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
	defer func() {
		conn.Close()
		if e := recover(); e != nil {
			Log.Alert(e)
		}
	}()

	for {
		request := HandleRead(conn, make([]byte, HEADER_LENGTH), HandleHeader)
		response := HandleRead(conn, request, handleRequest)
		HandleWrite(conn, response)
	}
}

func HandleRead(conn net.Conn, buf []byte, handlerFunc func([]byte)([]byte, error)) ([]byte){
	_, err := io.ReadFull(conn, buf)
	if err != nil {
		panic(err)
		return nil
	}
	result, err := handlerFunc(buf)
	if err != nil {
		panic(err)
		return nil
	}
	return result
}

func HandleWrite(conn net.Conn, buf []byte) {
	header := GetHeader(buf)
	_, err := conn.Write(append(header, buf...))
	if err != nil {
		panic(err)
	}
}
