/**
 * user: jackong
 * Date: 11/13/13
 * Time: 3:31 PM
 */
package model

import (
	. "morning-dairy/global"
	"github.com/Jackong/db"
	_ "github.com/Jackong/db/mysql"
)

type user struct {
	db.Collection
}

var (
	User user
)

func init() {
	User = user{Collection("user")}
}

func (this user) Exist(name string) bool {
	total, err := this.Count(db.Cond{"name": name})
	if err != nil {
		Log.Alert(err)
		return true
	}
	if total > 0 {
		return true
	}
	return false
}

func (this user) Create(name, password string) bool {
	ids, err := this.Append(db.Item{
		"name": name,
		"password" : password,
	})
	if err != nil {
		Log.Alert(err)
		return false
	}
	if len(ids) < 1 {
		return false
	}
	return true
}

func (this user) Password(name string) string {
	item, err := this.Collection.Find(db.Cond{
		"name": name,
	})
	if err != nil {
		Log.Warn(err)
		return ""
	}
	return item["password"].(string)
}
