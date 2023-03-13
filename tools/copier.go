package tools

import (
	jsoniter "github.com/json-iterator/go"
	"github.com/mitchellh/mapstructure"
)

// Copier copy by json
func Copier(in, out interface{}) (err error) {
	var (
		b []byte
	)

	if b, err = jsoniter.Marshal(in); err != nil {
		return
	}
	return jsoniter.Unmarshal(b, out)
}

// MustDecode must decode
func MustDecode(in, out interface{}) {
	if err := mapstructure.Decode(in, out); err != nil {
		panic(err)
	}
	return
}
