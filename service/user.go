/**
 * User: jackong
 * Date: 11/12/13
 * Time: 9:23 PM
 */
package service
import (
	"morning-dairy/model"
	"github.com/gorilla/securecookie"
	"net/http"
	. "morning-dairy/global"
	"morning-dairy/err"
	"morning-dairy/io/output"
)

type user struct {
	secureCookie *securecookie.SecureCookie
}

const (
	COOKIE_SESSION = "gsession"
)

var (
	User user
)

func init() {
    User = user{secureCookie: securecookie.New([]byte(Project.String("cookie", "hash")), []byte(Project.String("cookie", "block")))}
}

func (this user) Exist(name string) bool {
	return model.User.Exist(name)
}


func (this user) Create(name, password string) bool {
	return model.User.Create(name, password)
}


func (this user) CanSignIn(name, password string) bool {
    return model.User.Password(name) == password
}

func (this user) SignIn(writer http.ResponseWriter, name, password string) {
	if encoded, err := this.secureCookie.Encode(COOKIE_SESSION, name); err == nil {
		cookie := &http.Cookie{
			Name:  COOKIE_SESSION,
			Value: encoded,
			Path:  "/",
		}
		http.SetCookie(writer, cookie)
	}
}

func (this user) IsSignIn(req *http.Request) (ok bool, name string) {
	if cookie, err := req.Cookie(COOKIE_SESSION); err == nil {
		var name string
		if err = this.secureCookie.Decode(COOKIE_SESSION, cookie.Value, &name); err == nil {
			if name != "" {
				return true, name
			}
		}
	}
	return false, ""
}

func (this user) SignOut(writer http.ResponseWriter, req *http.Request) {
	if cookie, err := req.Cookie(COOKIE_SESSION); err == nil {
		cookie.MaxAge = -1
		cookie.Expires = Time()
		http.SetCookie(writer, cookie)
	}
}

func Before2Sign(writer http.ResponseWriter, req *http.Request) (e error) {
	if ok, name := User.IsSignIn(req); ok {
		cookie := &http.Cookie{
			Name: "name",
			Value: name,
			Path: "/",
		}
		req.AddCookie(cookie)
		return
	}
	output.Puts(writer, "code", CODE_SIGN_IN_REQUIRED)
	return err.AccessError{Status: http.StatusBadRequest, Msg: "required sign in"}
}
