/**
 * User: jackong
 * Date: 11/15/13
 * Time: 9:50 AM
 */
package input

import (
	"testing"
	"net/http"
	"net/url"
	"morning-dairy/err"
)

var (
	req *http.Request
)

func init() {
	req = &http.Request{Form: url.Values{
		"key": []string{"value"},
		"key2": []string{"1234"},
	}}
}
func TestDefault(t *testing.T) {
	defo := Default(req, "key", "default")
	if defo != "value" {
		t.Error("could not get the value passed")
	}
	defo = Default(req, "notExist", "default")
	if defo != "default" {
		t.Error("could not get the default value")
	}
}

func TestPattern(t *testing.T) {
	value := Pattern(req, "key", "")
	if value != "value" {
		t.Error("could not get value by empty pattern")
	}

	value2 := Pattern(req, "key2", "^[1-4]{4}$")
	if value2 != "1234" {
		t.Error("could not get value by pattern")
	}
}

func TestPatternButNotFound(t *testing.T) {
	defer func() {
		if e := recover(); e != nil {
			switch e.(type) {
			case err.AccessError:
				return
			default:
				t.Error("should be AccessError panic", e)
			}
		}
	}()
	Pattern(req, "key2", "[a-z]")
	t.Error("should be panic")
}


func TestRequired(t *testing.T) {
    value := Required(req, "key")
	if value != "value" {
		t.Error("could not get the required value")
	}
}

func TestRequiredButNotFound(t *testing.T) {
	defer func () {
	    if e := recover(); e != nil {
			switch e.(type) {
			case err.AccessError:
				return
			default:
				t.Error("should be AccessError panic", e)
			}
		}
	}()

	Required(req, "key3")
	t.Error("should be panic")
}

func TestPatternButDefault(t *testing.T) {
    value := Get(req, "key3", "[0-4]", "4")
	if value != "4" {
		t.Error("could not get the default value by pattern")
	}
}
