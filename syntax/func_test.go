/**
 * User: jackong
 * Date: 11/13/13
 * Time: 11:25 AM
 */
package syntax

import "testing"

type base interface {
	doSt() bool
}

type funcType func(base) bool

type child struct {

}

func (this *child) doSt() bool {
    return true
}

func doSt(f funcType) bool {
	c := &child{}
	return f(c)
}

func TestBase(t *testing.T) {
	doSt(func(b base) bool{
		return b.doSt()
	})

/*	//build failed
	doSt(func(c *child) bool{
	    return c.doSt()
	})*/
}
