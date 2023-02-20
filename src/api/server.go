package router

import (
	/*
	   "net/http"

	   "github.com/go-chi/chi/v5/middleware"
	*/
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	api_shorten_bulk "github.com/jei-el/vuo.be-backend/src/api/shorten-bulk"
)

func Serve() {
	r := chi.NewRouter()

	shortenBulkModule := api_shorten_bulk.NewShortenBulkModule()
	shortenBulkModule.Init(r)

	port := os.Getenv("SERVER_PORT")
	if port == "" {
		log.Fatalf("Server port not found")
	}
	port = ":" + port

	http.ListenAndServe(port, r)
}
