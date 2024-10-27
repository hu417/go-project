package utils

import (
	"encoding/json"
	"fmt"
)

func JsonMarshal(v interface{}) string {
	claimsJSON, err := json.Marshal(v)
	if err != nil {
		panic(fmt.Errorf("claims json marshal error: %v", err))

	}
	return string(claimsJSON)
}

func JsonUnMarshal(str string, v interface{}) (m interface{}, err error) {
	if err := json.Unmarshal([]byte(str), v); err != nil {
		return nil, err
	}
	return v, nil
}
