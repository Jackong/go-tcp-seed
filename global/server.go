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
		connection := &Connection{Id: conn.RemoteAddr().String(), Conn: conn}
		Anonymous.Put(connection.Id, connection)
		go handleClient(connection)
	}

	ShutDown()
}

func handleClient(conn *Connection) {
	defer func() {
		if conn.IsSigned {
			Signed.Close(conn.Id)
		} else {
			Anonymous.Close(conn.Id)
		}
		if e := recover(); e != nil {
			Log.Alert(e)
		}
	}()

	for {
		header := make([]byte, HEADER_LENGTH)
		HandleRead(conn, header)
		request, err := HandleHeader(header)
		if err != nil {
			panic(err)
		}
		HandleRead(conn, request)
		response, err := handleRequest(request, conn)
		if err != nil {
			panic(err)
		}
		HandleWrite(conn, response)
	}
}

func HandleRead(conn net.Conn, buf []byte) {
	_, err := io.ReadFull(conn, buf)
	if err != nil {
		panic(err)
		return
	}
}

func HandleWrite(conn net.Conn, buf []byte) {
	header := GetHeader(buf)
	_, err := conn.Write(append(header, buf...))
	if err != nil {
		panic(err)
	}
}
