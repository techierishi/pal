package svcm

import (
	"net/http"
	"time"
)

func NewHttpClient() *http.Client {
	return &http.Client{
		Timeout: 500 * time.Millisecond,
	}
}
