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
	Handle(*pb.Request, *Connection) *pb.Response
}

type HandlerFunc func(*pb.Request, *Connection) *pb.Response
func (this HandlerFunc) Handle(req *pb.Request, conn *Connection) *pb.Response {
	return this(req, conn)
}
type AfterFunc func(*pb.Response)
type WrapHandler struct {
	Handler
	beforeHandlers []HandlerFunc
	afterHandlers []AfterFunc
}

func (this *WrapHandler) Handle(req *pb.Request, conn *Connection) *pb.Response {
	for _, beforeHandler := range this.beforeHandlers {
		if resp := beforeHandler(req, conn); resp != nil {
			return resp
		}
	}

	resp := this.Handler.Handle(req, conn)

	for _, afterFunc := range this.afterHandlers {
		afterFunc(resp)
	}
	return resp
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

func Handle(request *pb.Request, conn *Connection) (resp *pb.Response) {
	module := request.GetModule()
	handler, ok := handlers[module]
	if !ok {
		resp = new(pb.Response)
		Log.Alert("could not found this module", module)
		resp.Code = pb.Code_MODULE_NOT_EXIST.Enum()
		return
	}
	return handler.Handle(request, conn)
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
	response := Handle(request, conn)
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
