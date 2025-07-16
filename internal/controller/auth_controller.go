package controller

import (
	"encoding/json"
	"go-crud-oapi/internal/repository"
	"go-crud-oapi/pkg/auth"
	"log"
	"net/http"

	"golang.org/x/crypto/bcrypt"
)

type AuthController struct {
	Repo repository.UserRepoInterface
}

func NewAuthController(repo repository.UserRepoInterface) *AuthController {
	return &AuthController{Repo: repo}
}

func (a *AuthController) Login(w http.ResponseWriter, r *http.Request) {
	var creds struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	if err := json.NewDecoder(r.Body).Decode(&creds); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	user, err := a.Repo.FindByEmail(r.Context(), creds.Email)
	if err != nil || user == nil {
		http.Error(w, "User not found", http.StatusUnauthorized)
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(creds.Password)); err != nil {
		http.Error(w, "Invalid password", http.StatusUnauthorized)
		return
	}

	if user.Role != "admin" {
		http.Error(w, "Unauthorized - admin access only", http.StatusForbidden)
		return
	}

	token, err := auth.GenerateToken(user.Email, user.Role)

	if err != nil {
		log.Printf("JWT Signing failed: %v", err)
		http.Error(w, "Failed to generate token", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(map[string]string{"token": token})

}
