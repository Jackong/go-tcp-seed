/**
 * User: jackong
 * Date: 11/12/13
 * Time: 5:01 PM
 */
package output

import (
	"net/http"
	"encoding/json"
	"fmt"
)

type Json struct {
	http.ResponseWriter
	ret map[string]interface {}
}

func NewJson(writer http.ResponseWriter) *Json{
	return &Json{ResponseWriter: writer, ret: make(map[string]interface {})}
}

func (this *Json) Puts(ret []interface {}) {
	length := len(ret)
	length -= (length % 2)
	for index := 0; index < length; index += 2 {
		this.ret[ret[index].(string)] = ret[index + 1]
	}
}

func (this *Json) Render(status int, msg string) (err error) {
	this.ret["status"] = status
	this.ret["msg"] = msg

	var data []byte
	data, err = json.Marshal(this.ret)
	if err != nil {
		return err
	}
	//this.ResponseWriter.Header().Set("Content-Type", "applicatoin/json; charset=utf-8")
	fmt.Fprint(this.ResponseWriter, string(data))
	return
}
