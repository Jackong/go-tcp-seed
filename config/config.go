/**
 * User: jackong
 * Date: 11/5/13
 * Time: 5:22 PM
 */
package config

import (
	"io/ioutil"
	"encoding/json"
	"fmt"
	"os"
)

type Config map[string]interface {}

func NewConfig(file string) Config {
	config := make(Config)
	data, err := ioutil.ReadFile(file)
	if err != nil {
		fmt.Println(err)
		os.Exit(2)
	}
	json.Unmarshal(data, &config)
	return config
}

func (this Config) Get(keys ... string) (val interface {}) {
	tmp := this
	for _, key := range keys {
		val = tmp[key]
		if val == nil {
			return nil
		}
		switch val.(type) {
		case map[string] interface {}:
			tmp = val.(map[string] interface {})
		default:
			return val
		}
	}
	return val
}

func (this Config) String(keys ... string) string {
	val := this.Get(keys...)
	if val == nil {
		return ""
	}
	return val.(string)
}

func (this Config) Slice(keys ... string) (array []string) {
	val := this.Get(keys...)
	if val == nil {
		return
	}
	for _, element := range val.([]interface {}) {
		array = append(array, element.(string))
	}
	return
}
