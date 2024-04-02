package main

import (
	"fmt"
	"io"
	"net/http"
)

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("Starting API")
		io.WriteString(w, "Starting API")
	})

	http.ListenAndServe(":3000", nil)
}
