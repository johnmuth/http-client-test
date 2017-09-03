package main

import (
	"encoding/json"
)

type ServiceRequest struct {
	UUID string `json:"uuid,omitempty"`
}

func (req ServiceRequest) String() string {
	out, _ := json.Marshal(req)
	return string(out)
}
