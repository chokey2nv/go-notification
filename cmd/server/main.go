package server

import (
	"net/http"

	httpapi "github.com/chokey2nv/go-notification/api/http"
	"github.com/chokey2nv/go-notification/samples"
)

func main() {
	mux := http.NewServeMux()

	svc := samples.DefaultNotificationServer()
	handler := httpapi.NewHandler(svc)

	mux.HandleFunc("POST /v1/notifications", handler.Create)
	mux.HandleFunc("GET /v1/notifications/{id}", handler.GetByID)
	mux.HandleFunc("GET /v1/users/{user_id}/notifications", handler.GetByUser)

	http.ListenAndServe(":8080", mux)

}
