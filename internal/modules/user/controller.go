package user

import (
	"encoding/json"
	"github.com/go-chi/chi/v5"
	"net/http"
)

type Controller struct {
	userService Service
}

func NewUserController(userService Service) *Controller {
	return &Controller{userService: userService}
}

func (c *Controller) RegisterRoutes(r chi.Router) {
	r.Post("/public/auth/signup", c.SignUp)
}

func (c *Controller) SignUp(w http.ResponseWriter, r *http.Request) {
	var request SignUpRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if _, err := c.userService.SignUp(&request); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	err := json.NewEncoder(w).Encode(request)
	if err != nil {
		return
	}
}

func (c *Controller) SignIn(w http.ResponseWriter, r *http.Request) {
	var request SignInRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	response, err := c.userService.SignIn(&request)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func (c *Controller) GetByLogin(w http.ResponseWriter, r *http.Request) {
	login := chi.URLParam(r, "login")
	if login == "" {
		http.Error(w, "Login parameter is required", http.StatusBadRequest)
		return
	}

	response, err := c.userService.GetByLogin(login)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(response)
}
