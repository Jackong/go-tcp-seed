/**
 * User: Jackong
 * Date: 13-11-24
 * Time: 上午11:27
 */
package global

import (
	"github.com/Jackong/go-tcp-seed/pb"
	"code.google.com/p/goprotobuf/proto"
)

const (
	HEADER_LENGTH = 9
)
type Handler interface {
	Handle(*pb.Request, *pb.Response, *Connection) error
}

type HandlerFunc func(*pb.Request, *pb.Response, *Connection) error
func (this HandlerFunc) Handle(req *pb.Request, resp *pb.Response, conn *Connection) error {
	return this(req, resp, conn)
}
type AfterFunc func(*pb.Response)
type WrapHandler struct {
	Handler
	beforeHandlers []HandlerFunc
	afterHandlers []AfterFunc
}

func (this *WrapHandler) Handle(req *pb.Request, resp *pb.Response, conn *Connection) error {
	for _, beforeHandler := range this.beforeHandlers {
		err := beforeHandler(req, resp, conn)
		if err != nil {
			return err
		}
		if resp.Code.String() != pb.Code_OK.String() {
			return nil
		}
	}

	err := this.Handler.Handle(req, resp, conn)
	if err != nil {
		return err
	}

	for _, afterFunc := range this.afterHandlers {
		afterFunc(resp)
	}
	return nil
}

func (this *WrapHandler) Before(beforeHandlers ...HandlerFunc) *WrapHandler{
	this.beforeHandlers = append(this.beforeHandlers, beforeHandlers...)
	return this
}

func (this *WrapHandler) After(afterHandlers ...AfterFunc) *WrapHandler {
	this.afterHandlers = append(this.afterHandlers, afterHandlers...)
	return this
}

var (
	handlers map[pb.Module] Handler
)

func init() {
	handlers = make(map[pb.Module] Handler)
}

func Register(module pb.Module, handler Handler) *WrapHandler{
	wrapHandler := &WrapHandler{Handler: handler}
	handlers[module] = wrapHandler
	return wrapHandler
}

func Handle(request *pb.Request, response *pb.Response, conn *Connection) error {
	module := request.GetModule()
	handler, ok := handlers[module]
	if !ok {
		response = &pb.Response{Code: pb.Code_MODULE_NOT_EXIST.Enum()}
		Log.Alertf("%v request not exist module %v with %v", conn.Conn.RemoteAddr(), module, request)
		return nil
	}
	return handler.Handle(request, response, conn)
}

func HandleHeader(buf []byte) ([]byte, error) {
	header := &pb.Header{}
	err := proto.Unmarshal(buf, header)
	if err != nil {
		return nil, err
	}
	return make([]byte, header.GetLength()), err
}

func handleRequest(buf []byte, conn *Connection) (resp []byte, err error) {
	request := &pb.Request{}
	err = proto.Unmarshal(buf, request)
	if err != nil {
		return resp, err
	}
	response := &pb.Response{}
	err = Handle(request, response, conn)
	if err != nil {
		Log.Alert(err)
		response.Code = pb.Code_BAD_REQUEST.Enum()
	}
	response.Module = request.Module
	resp, err = proto.Marshal(response)
	if err != nil {
		Log.Alert(err)
		newResponse := &pb.Response{Code: pb.Code_BAD_RESPONSE.Enum(), Module: request.Module}
		resp, _ = proto.Marshal(newResponse)
	}
	return resp , nil
}

func GetHeader(protocol []byte) []byte {
	header := &pb.Header{Length: proto.Uint64(uint64(len(protocol)))}
	headerBuf, _ := proto.Marshal(header)
	return headerBuf
}
