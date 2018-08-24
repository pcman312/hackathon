package services

import (
	"fmt"
	"net/http"
)

func HelloWorld(rw http.ResponseWriter, req *http.Request) {
	rw.WriteHeader(http.StatusOK)
	fmt.Fprintln(rw, "Hello, LogRhythm Hackathon! This is a new feature!")
}
