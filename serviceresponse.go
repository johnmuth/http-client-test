package main

import (
	"encoding/json"
)

type ServiceResponse struct {
	Qux string `json:"qux"`
	UUID string `json:"uuid"`
}

func (sr ServiceResponse) String() string {
	out, _ := json.Marshal(sr)
	return string(out)
}
