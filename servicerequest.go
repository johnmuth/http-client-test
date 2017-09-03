package main

import (
	"encoding/json"
)

type ServiceRequest struct {
	RequestID string `json:"requestid,omitempty"`
}

func (req ServiceRequest) String() string {
	out, _ := json.Marshal(req)
	return string(out)
}
