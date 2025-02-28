package handlers

import (
	"net/http"
)

// HelloHandler handles the /api/hello endpoint
func HelloHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	w.Write([]byte("<div class='alert alert-success'>Hello from the server!</div>"))
}
