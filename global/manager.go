/**
 * User: Jackong
 * Date: 13-11-25
 * Time: 下午8:26
 */
package global

import (
	"net"
)

type manager struct {
	connections map[string]net.Conn
}
var (
	Anonymous *manager
	Signed *manager
)

func init() {
	Anonymous = &manager{connections: make(map[string]net.Conn)}
	Signed = &manager{connections: make(map[string]net.Conn)}
	OnShutDown(func() {
		Anonymous.CloseAll()
		Signed.CloseAll()
	})
}

func (this *manager) Get(id string) net.Conn {
	return this.connections[id]
}

func (this *manager) Put(id string, conn net.Conn) {
	this.Close(id)
	this.connections[id] = conn
}

func (this *manager) Del(id string) {
	delete(this.connections, id)
}

func (this *manager) CloseAll() {
	for id, conn := range this.connections {
		conn.Close()
		this.Del(id)
	}
}

func (this *manager) Close(id string) {
	if _, ok := this.connections[id]; ok {
		this.connections[id].Close()
		this.Del(id)
	}
}

func SignIn(aid, sid string) bool {
	conn := Anonymous.Get(aid)
	if conn == nil {
		return false
	}
	Signed.Put(sid, conn)
	return true
}

