package handlers

import (
	"log"
	"net/http"

	"github.com/tm-acme-shop/acme-shop-gateway/internal/config"
	"github.com/tm-acme-shop/acme-shop-gateway/internal/proxy"
)

type UsersHandler struct {
	config *config.Config
	client *proxy.Client
}

func NewUsersHandler(cfg *config.Config, client *proxy.Client) *UsersHandler {
	return &UsersHandler{
		config: cfg,
		client: client,
	}
}

func (h *UsersHandler) GetUserV1(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	log.Printf("GetUserV1: id=%s", id)
	h.client.ProxyRequest(w, r, h.config.UsersServiceURL+"/users/"+id)
}

func (h *UsersHandler) CreateUserV1(w http.ResponseWriter, r *http.Request) {
	log.Printf("CreateUserV1")
	h.client.ProxyRequest(w, r, h.config.UsersServiceURL+"/users")
}
