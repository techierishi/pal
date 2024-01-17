package svcm

type StatusResponse struct {
	Status  bool   `json:"status"`
	Version string `json:"version"`
	Commit  string `json:"commit"`
}
