package api_shorten_bulk

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
)

type ShortenBulkController struct {
	service *ShortenBulkService
}

func (c *ShortenBulkController) Post(w http.ResponseWriter, r *http.Request) {
	url := chi.URLParam(r, "url")
	mp, statusCode := c.service.Post(url)

	w.WriteHeader(statusCode)
	w.Header().Set("Content-Type", "application/json")
	res, err := json.Marshal(mp)
	if err != nil {
		log.Fatalf("Error happened in JSON marshal. Err: %s", err.Error())
	}

	w.Write(res)

}

func (c *ShortenBulkController) Get(w http.ResponseWriter, r *http.Request) {
	hash := chi.URLParam(r, "hash")
	mp, statusCode := c.service.Get(hash)

	w.WriteHeader(statusCode)
	w.Header().Set("Content-Type", "application/json")
	res, err := json.Marshal(mp)
	if err != nil {
		log.Fatalf("Error happened in JSON marshal. Err: %s", err.Error())
	}

	w.Write(res)
}

func (c *ShortenBulkController) Route(r *chi.Mux) {

	r.Get("/", c.Get)
	r.Post("/{hash}", c.Post)
}

func NewShortenBulkController(service *ShortenBulkService) *ShortenBulkController {
	return &ShortenBulkController{
		service,
	}
}
