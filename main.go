/**
 * User: jackong
 * Date: 10/17/13
 * Time: 6:24 PM
 */
package main

import (
	. "github.com/Jackong/go-tcp-seed/global"
	"github.com/Jackong/go-tcp-seed/pb"
	"fmt"
)


func main() {
	Register(pb.Module_SIGN_UP, &signUp{}).Before(func(req *pb.Request, conn *Connection) (res *pb.Response) {
		fmt.Println("aha!")
		return &pb.Response{Code: pb.Code_BAD_REQUEST.Enum()}
	})
	SetUp()
}

type signUp struct {

}

func (this *signUp)Handle(req *pb.Request, conn *Connection) (res *pb.Response) {
	fmt.Println("sign up..")
	sign := req.GetSign()
	SignIn(conn.Id, sign.GetEmail())
	res = &pb.Response{Code: pb.Code_OK.Enum()}
	return res
}
