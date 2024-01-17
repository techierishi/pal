package svcm

import (
	"encoding/json"
	"net/http"

	"log"
)

type statusHandler struct {
}

func (h *statusHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	resp := map[string]string{
		"Status": "true",
	}
	jsn, err := json.Marshal(&resp)
	if err != nil {
		log.Fatalf("Error when marshaling %v", err)
		return
	}
	w.Write(jsn)
}
