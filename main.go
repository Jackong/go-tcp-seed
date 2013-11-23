/**
 * User: jackong
 * Date: 10/17/13
 * Time: 6:24 PM
 */
package main

import (
	"fmt"
	"github.com/braintree/manners"
	. "morning-dairy/global"
	_ "morning-dairy/router/sign"
)


func main() {
	err := manners.ListenAndServe(Project.String("server", "addr"), Router)
	if	err != nil {
		fmt.Println(err)
	}
	ShutDown()
}
