/**
 * User: jackong
 * Date: 10/17/13
 * Time: 6:24 PM
 */
package main

import (
	"fmt"
	"github.com/braintree/manners"
	. "github.com/Jackong/go-web-seed/global"
)


func main() {
	err := manners.ListenAndServe(Project.String("server", "addr"), Router)
	if	err != nil {
		fmt.Println(err)
	}
	ShutDown()
}
