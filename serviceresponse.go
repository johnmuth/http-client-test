package main

import (
	"encoding/json"
)

type ServiceResponse struct {
	Qux       string `json:"qux"`
	RequestID string `json:"requestid"`
}

func (sr ServiceResponse) String() string {
	out, _ := json.Marshal(sr)
	return string(out)
}
