/**
 * User: Jackong
 * Date: 13-11-24
 * Time: 上午11:27
 */
package global

import (
	"github.com/Jackong/go-tcp-seed/pb"
)

type Handler interface {
	Handle(protocol *pb.Request) *pb.Response
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

func Handle(module pb.Module, protocol *pb.Request) (resp *pb.Response) {
	handler, ok := handlers[module]
	if !ok {
		Log.Alert("could not found this module", module)
		resp.Code = pb.Code_MODULE_NOT_EXIST
		return
	}
	return handler.Handle(protocol)
}
