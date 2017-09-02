package main

import (
	"encoding/json"
)

type ServiceRequest struct {
	Foo string `json:"foo,omitempty"`
	Bar string `json:"bar,omitempty"`
}

func (req ServiceRequest) String() string {
	out, _ := json.Marshal(req)
	return string(out)
}
