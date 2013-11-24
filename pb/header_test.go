/**
 * User: Jackong
 * Date: 13-11-24
 * Time: 下午10:40
 */
package pb

import (
	"testing"
	"code.google.com/p/goprotobuf/proto"
)

func TestHeaderSize(t *testing.T) {
	header := &Header{Length: proto.Uint64(222)}
	buf, err := proto.Marshal(header)
	if err != nil {
		t.Error(err)
	}
	t.Log(len(buf))
}
