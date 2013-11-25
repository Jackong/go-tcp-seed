/**
 * User: Jackong
 * Date: 13-11-25
 * Time: 下午8:26
 */
package global

type manager struct {
	connections map[string]*Connection
}
var (
	Anonymous *manager
	Signed *manager
)

func init() {
	Anonymous = &manager{connections: make(map[string]*Connection)}
	Signed = &manager{connections: make(map[string]*Connection)}
	OnShutDown(func() {
		Anonymous.CloseAll()
		Signed.CloseAll()
	})
}

func (this *manager) Get(id string) *Connection {
	return this.connections[id]
}

func (this *manager) Put(id string, conn *Connection) {
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
	conn.IsSigned = true
	Signed.Put(sid, conn)
	return true
}

func SignOut(sid, aid string) bool {
	conn := Signed.Get(sid)
	if conn == nil {
		return false
	}
	conn.IsSigned = false
	Anonymous.Put(aid, conn)
	return true
}
