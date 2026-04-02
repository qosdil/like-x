package http

import (
	user "likexuser/model"
	"likexuser/service"
	"log"
	"net/http"

	"github.com/gofiber/fiber/v3"
	httphandler "github.com/qosdil/like-x/backend/common/http/handler"
)

// NewHandler creates an HTTP handler with the given service dependency.
func NewHandler(svc *service.Service) *handler { return &handler{svc: svc} }

// HandleSignUp handles sign-up requests, validates input, and returns HTTP response.
func (h *handler) HandleSignUp(c fiber.Ctx) error {
	req := new(user.CreateInput)
	if err := c.Bind().Body(req); err != nil {
		log.Printf("error binding Body: %v", err)
		return c.SendStatus(http.StatusBadRequest)
	}

	signUp, err := h.svc.SignUp(c, *req)
	return httphandler.ObjResp(c, signUpResp{ID: signUp.PublicID}, err)
}

// handler handles HTTP requests and maps them to service operations.
type handler struct {
	svc *service.Service
}

// signUpResp is the HTTP response model for successful registration.
type signUpResp struct {
	ID user.PublicID `json:"id"` // We don't output internal "id"
}
