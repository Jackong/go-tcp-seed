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
	Register(pb.Module_SIGN_UP, &signUp{})
	SetUp()
}

type signUp struct {

}

func (this *signUp)Handle(*pb.Request) (res *pb.Response) {
	fmt.Println("sign up..")
	res = new(pb.Response)
	res.Code = pb.Code_OK.Enum()
	return res
}
