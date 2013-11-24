/**
 * User: Jackong
 * Date: 13-11-24
 * Time: 下午9:05
 */
package global
import (
	"net"
	"io"
	"github.com/Jackong/go-tcp-seed/pb"
	"code.google.com/p/goprotobuf/proto"
	"fmt"
)

type client interface {
	Request(*pb.Request) (*pb.Response)
}

type cln struct {
	net.Conn
}

func NewClient() client {
	conn, err := net.Dial("tcp", "localhost" + Project.String("server", "addr"))
	if err != nil {
		Log.Fatal(err)
		return nil
	}
	return &cln{Conn: conn}
}
func (this *cln) Request(req *pb.Request) (res *pb.Response) {
	reqBuf, _ := proto.Marshal(req)
	length := proto.Uint64(uint64(len(reqBuf)))
	header := &pb.Header{Length: length}
	hBuf, _ := proto.Marshal(header)
	fmt.Println("write")
	this.Conn.Write(append(hBuf, reqBuf...))
	rb := resBuf(this.Conn)
	fmt.Println("read")
	io.ReadFull(this.Conn, rb)
	res = new(pb.Response)
	err := proto.Unmarshal(rb, res)
	if err != nil {
		Log.Fatal(err)
	}
	return res
}

func resBuf(conn net.Conn) []byte {
	var header [HEADER_LENGTH]byte
	_, err := io.ReadFull(conn, header[0:])
	if err != nil {
		Log.Alert(err)
		return nil
	}
	response, err := HandleHeader(header[0:])
	if err != nil {
		Log.Alert(err)
		return nil
	}
	return response
}
