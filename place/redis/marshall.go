package redis

import "encoding/json"

var (
	JsonMarshal = func(v interface{}) (string, error) {
		bb, err := json.Marshal(v)
		return string(bb), err
	}

	JsonUnmarshal = func(s string, v interface{}) error {
		err := json.Unmarshal([]byte(s), &v)
		return err
	}
)
