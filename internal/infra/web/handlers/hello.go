package handlers

import (
	"net/http"

	"github.com/kameikay/rate-limiter/internal/infra/redis"
	"github.com/kameikay/rate-limiter/pkg/utils"
)

type HelloHandler struct {
	rdb *redis.Client
}

func NewHelloHandler(rdb *redis.Client) *HelloHandler {
	return &HelloHandler{
		rdb: rdb,
	}
}

func (h *HelloHandler) GetHelloHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		utils.JsonResponse(w, utils.ResponseDTO{
			StatusCode: http.StatusMethodNotAllowed,
			Message:    http.StatusText(http.StatusMethodNotAllowed),
			Success:    false,
		})
		return
	}

	utils.JsonResponse(w, utils.ResponseDTO{
		StatusCode: http.StatusOK,
		Message:    "success",
		Success:    true,
	})
}
