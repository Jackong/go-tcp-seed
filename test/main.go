/**
 * User: Jackong
 * Date: 13-11-24
 * Time: 上午10:15
 */
package main

import (
	"net"
	"fmt"
	. "github.com/Jackong/go-tcp-seed/global"
	"os"
)

func main() {
	fmt.Println("conn")
	conn, err := net.Dial("tcp", "localhost" + Project.String("server", "addr"))
	if err != nil {
		fmt.Println(err)
		os.Exit(2)
	}
	fmt.Println("write")
	fmt.Fprintf(conn, "hello")
	fmt.Println("read")
	var buf [512]byte
	n, err := conn.Read(buf[0:])
	if err != nil {
		fmt.Println(err)
		os.Exit(2)
	}
	res := string(buf[0:n])
	fmt.Println(res)
	if res != "hello" {
		fmt.Println("req != res")
		os.Exit(2)
	}
}
