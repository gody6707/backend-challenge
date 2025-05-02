package handlers

import (
	"backend-challenge/models"
	"backend-challenge/repository"
	"backend-challenge/utils"
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
	"golang.org/x/crypto/bcrypt"
)

type UserHandler struct {
	repo      repository.UserRepository
	jwtSecret string
}

func NewUserHandler(repo repository.UserRepository, jwtSecret string) *UserHandler {
	return &UserHandler{repo: repo, jwtSecret: jwtSecret}
}

func (h *UserHandler) Register(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Name     string `json:"name"`
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	hashed, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		http.Error(w, "Could not hash password", http.StatusInternalServerError)
		return
	}
	user := models.User{Name: req.Name, Email: req.Email, Password: string(hashed), CreatedAt: time.Now()}
	if err := h.repo.Create(context.Background(), &user); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
}

func (h *UserHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	var users []models.User
	if err := h.repo.List(context.Background(), bson.M{"email": req.Email}, &users); err != nil || len(users) == 0 {
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}
	user := users[0]
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}
	token, err := utils.GenerateToken(user.ID.Hex(), h.jwtSecret)
	if err != nil {
		http.Error(w, "Token error", http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(map[string]string{"token": token})
}

func (h *UserHandler) List(w http.ResponseWriter, r *http.Request) {
	var users []models.User
	if err := h.repo.List(context.Background(), bson.D{}, &users); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(users)
}

func (h *UserHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	var user models.User
	if err := h.repo.GetByID(context.Background(), id, &user); err != nil {
		http.Error(w, "Not found", http.StatusNotFound)
		return
	}
	json.NewEncoder(w).Encode(user)
}

func (h *UserHandler) Update(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	var req struct {
		Name  string `json:"name,omitempty"`
		Email string `json:"email,omitempty"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	update := make(map[string]interface{})
	if req.Name != "" {
		update["name"] = req.Name
	}
	if req.Email != "" {
		update["email"] = req.Email
	}
	if err := h.repo.Update(context.Background(), id, update); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *UserHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	if err := h.repo.Delete(context.Background(), id); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
