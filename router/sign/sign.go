/**
 * User: jackong
 * Date: 11/12/13
 * Time: 8:17 PM
 */
package sign

import (
	"net/http"
	"morning-dairy/io/input"
	."morning-dairy/global"
	"morning-dairy/service"
	"morning-dairy/io/output"
)


const (
	RE_EMAIL = `(?i)[A-Z0-9._%+-]+@(?:[A-Z0-9-]+\.)+[A-Z]{2,6}`
	RE_PASSWORD = `[0-9a-f]{32}`
)

func init() {
	Router.HandleFunc("/sign/up", signUp).Methods("POST", "GET")
	Router.HandleFunc("/sign/in", signIn).Methods("POST", "GET")
	Router.HandleFunc("/sign/out", signOut).Methods("POST", "GET")
}

func signUp(writer http.ResponseWriter, req * http.Request) {
	name, password := nameAndPassword(req)
	if service.User.Exist(name) {
		output.Puts(writer, "code", CODE_EXIST)
		return
	}
	code := CODE_FAIL
	if service.User.Create(name, password) {
		code = CODE_OK
	}
	service.User.SignIn(writer, name, password)
	output.Puts(writer, "code", code)
}

func signIn(writer http.ResponseWriter, req *http.Request) {
	name, password := nameAndPassword(req)
	code := CODE_FAIL
	if service.User.CanSignIn(name, password) {
		code = CODE_OK
	}
	service.User.SignIn(writer, name, password)
	output.Puts(writer, "code", code)
}

func signOut(writer http.ResponseWriter, req *http.Request) {
	service.User.SignOut(writer, req)
	output.Puts(writer, "code", CODE_OK)
}

func nameAndPassword(req *http.Request) (name, password string){
	name = input.Pattern(req, "name", RE_EMAIL)
	password = input.Pattern(req, "password", RE_PASSWORD)
	return
}
