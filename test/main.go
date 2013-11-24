/**
 * User: Jackong
 * Date: 13-11-24
 * Time: 上午10:15
 */
package main

import (
	. "github.com/Jackong/go-tcp-seed/global"
	"github.com/Jackong/go-tcp-seed/pb"
	"code.google.com/p/goprotobuf/proto"
	"fmt"
)

func main() {
	cln := NewClient()
	response := cln.Request(&pb.Request{
		Module: pb.Module_SIGN_UP.Enum(),
		Sign: &pb.Sign{Email: proto.String("email"), Password: proto.String("password")},
	})
	fmt.Println(response)
}
