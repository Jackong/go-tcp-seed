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

type Handler interface {
	Handle(*pb.Request) *pb.Response
}

var (
	handlers map[pb.Module] Handler
)

func init() {
	handlers = make(map[pb.Module] Handler)
}

func Register(module pb.Module, handler Handler) {
	handlers[module] = handler
}

func Handle(request *pb.Request) (resp *pb.Response) {
	module := request.GetModule()
	handler, ok := handlers[module]
	if !ok {
		resp = new(pb.Response)
		Log.Alert("could not found this module", module)
		resp.Code = pb.Code_MODULE_NOT_EXIST.Enum()
		return
	}
	return handler.Handle(request)
}

func HandleHeader(buf []byte) ([]byte, error) {
	header := &pb.Header{}
	err := proto.Unmarshal(buf, header)
	if err != nil {
		return nil, err
	}
	err = checkSum(header, header.GetCheckSum())
	if err != nil {
		return nil, err
	}
	return make([]byte, header.GetLength()), err
}

func handleRequest(buf []byte) (resp []byte, err error) {
	request := &pb.Request{}
	err = proto.Unmarshal(buf, request)
	if err != nil {
		return resp, err
	}
	err = checkSum(request, request.GetCheckSum())
	if err != nil {
		return resp, err
	}
	response := Handle(request)
	response.Module = request.Module
	setCheckSum(response)
	resp, err = proto.Marshal(response)
	if err != nil {
		Log.Alert(err)
		newResponse := &pb.Response{Code: pb.Code_BAD_RESPONSE.Enum(), Module: request.Module}
		setCheckSum(newResponse)
		resp, _ = proto.Marshal(newResponse)
	}
	header := &pb.Header{Length: proto.Uint64(uint64(len(resp))), CheckSum: proto.Uint32(264)}
	headerBuf, _ := proto.Marshal(header)
	return append(headerBuf, resp...) , nil
}

func checkSum(msg proto.Message, checkSum uint32) error {
	return nil
}

func setCheckSum(response *pb.Response) {
	response.CheckSum = proto.Uint32(123456)
}
