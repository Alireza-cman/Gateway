package models

type Device struct {
	Name   string `json:"name"`
	Status string `json:"status"`
	Time   int64  `json:latest`
}
