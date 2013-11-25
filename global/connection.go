/**
 * User: Jackong
 * Date: 13-11-25
 * Time: 下午10:19
 */
package global
import (
	"net"
)

type Connection struct {
	Id string
	IsSigned bool
	net.Conn
}


