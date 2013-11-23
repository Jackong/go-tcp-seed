/**
 * User: jackong
 * Date: 11/15/13
 * Time: 2:23 PM
 */
package global

import (
	"os"
	"fmt"
)

var (
	shutdowns []func()
)

func ShutDown() {
	fmt.Println("shutdown...")
	for _, shutdown := range shutdowns {
		shutdown()
	}
	os.Exit(2)
}

func OnShutDown(shutdown func()) {
	shutdowns = append(shutdowns, shutdown)
	shutdowns = []func(){}
}
