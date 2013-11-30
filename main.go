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
	AttachFunc(pb.Module_SIGN_UP, func(req *pb.Request, res *pb.Response, conn *Connection) error {
			fmt.Println("sign up...")
			sign := req.GetSign()
			SignIn(conn.Id, sign.GetEmail())
			res.Code = pb.Code_OK.Enum()
			return nil
		}).Before(func(req *pb.Request, res *pb.Response, conn *Connection) error {
		fmt.Println("aha!")
		return nil
	})
	SetUp()
}

type signUp struct {

}

func (this *signUp)Handle(req *pb.Request, res *pb.Response, conn *Connection) error {
	fmt.Println("sign up..")
	sign := req.GetSign()
	SignIn(conn.Id, sign.GetEmail())
	res.Code = pb.Code_OK.Enum()
	return nil
}
