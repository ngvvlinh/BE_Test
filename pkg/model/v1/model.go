package umodel

import (
	"encoding/json"
)

type (
	Message struct {
		T int             `json:"t"`
		D json.RawMessage `json:"d"`
		E []int           `json:"e,omitempty"`
		S int             `json:"s"`
	}

	Data struct {
		ID   int    `json:"i"`
		T    int    `json:"t"`
		Val1 string `json:"s"`
		Val2 string `json:"s2,omitempty"`
		Val3 string `json:"s3,omitempty"`
		Num1 int    `json:"n"`
		Num2 int    `json:"n2,omitempty"`
		Num3 int    `json:"n3,omitempty"`
	}
)

func UnmarshalData(data []byte) (Data, error) {
	var s Data
	err := json.Unmarshal(data, &s)
	return s, err
}
