package chaincode

import "encoding/json"

func ParseJSON[T ITYPES](data string) (T, error) {
	var emptyObject T
	var result T

	if err := json.Unmarshal([]byte(data), &result); err != nil {
		return emptyObject, err
	}
	if err := result.Validate(); err != nil {
		return emptyObject, err
	}

	return result, nil
}
