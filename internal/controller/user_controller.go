package controller

import (
	"encoding/json"
	"go-crud-oapi/internal/model"
	"go-crud-oapi/internal/repository"
	"go-crud-oapi/internal/service"
	"go-crud-oapi/pkg/logger"
	"go-crud-oapi/pkg/utils"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
)

type UserController struct {
	Repo repository.UserRepoInterface
	svc  service.UserServiceInterFace
}

func NewUserController(svc service.UserServiceInterFace) *UserController {
	return &UserController{svc: svc}
}

func (c *UserController) ListUsers(w http.ResponseWriter, r *http.Request) {
	log := logger.L(r.Context())
	log.Info("ListUsers handler invoked")

	users, err := c.svc.ListAllUsers(r.Context())
	if err != nil {
		log.Error("Failed to list users", zap.Error(err))
		utils.WriteJSONError(w, http.StatusBadRequest)
		return
	}

	log.Info("Successfully retrieved users", zap.Int("count", len(users)))
	json.NewEncoder(w).Encode(users)
}

func (c *UserController) CreateUser(w http.ResponseWriter, r *http.Request) {
	log := logger.L(r.Context())
	log.Info("CreateUser handler invoked")

	var user model.User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		log.Warn("Invalid request payload", zap.Error(err))
		utils.WriteJSONError(w, http.StatusBadRequest)
		return
	}

	// Check if email already exists
	existing, err := c.svc.GetUserByEmail(r.Context(), user.Email)
	if err != nil {
		log.Error("Failed to check email existence", zap.Error(err))
		utils.WriteJSONError(w, http.StatusInternalServerError)
		//http.Error(w, "Server error", http.StatusInternalServerError)
		return
	}
	if existing != nil {
		log.Warn("Email already in use", zap.String("email", user.Email))
		utils.WriteJSONError(w, http.StatusConflict)
		return
	}

	if err := c.svc.Create(r.Context(), &user); err != nil {
		log.Error("Failed to create user", zap.Error(err))
		utils.WriteJSONError(w, http.StatusInternalServerError)
		return
	}

	log.Info("User created successfully", zap.Uint("user_id", user.ID))
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(user)
}

func (c *UserController) GetUser(w http.ResponseWriter, r *http.Request) {
	log := logger.L(r.Context())
	idParam := chi.URLParam(r, "id")
	log.Info("GetUser handler invoked", zap.String("user_id_param", idParam))

	id, err := strconv.Atoi(idParam)
	if err != nil {
		log.Warn("Invalid user ID", zap.Error(err))
		utils.WriteJSONError(w, http.StatusBadRequest)
		//http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	user, err := c.svc.Get(r.Context(), uint(id))
	if err != nil {
		log.Error("User not found", zap.Int("user_id", id), zap.Error(err))
		utils.WriteJSONError(w, http.StatusNotFound)
		return
	}

	log.Info("User retrieved", zap.Uint("user_id", user.ID))
	json.NewEncoder(w).Encode(user)
}

func (c *UserController) UpdateUser(w http.ResponseWriter, r *http.Request) {
	log := logger.L(r.Context())
	idParam := chi.URLParam(r, "id")
	log.Info("UpdateUser handler invoked", zap.String("user_id_param", idParam))

	id, err := strconv.Atoi(idParam)
	if err != nil {
		log.Warn("Invalid user ID", zap.Error(err))
		utils.WriteJSONError(w, http.StatusBadRequest)
		return
	}

	var user model.User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		log.Warn("Invalid request payload", zap.Error(err))
		utils.WriteJSONError(w, http.StatusBadRequest)
		return
	}

	existing, err := c.svc.GetUserByEmail(r.Context(), user.Email)
	if err != nil {
		log.Error("Failed to check email existence", zap.Error(err))
		utils.WriteJSONError(w, http.StatusInternalServerError)
		return
	}
	if existing != nil {
		log.Warn("Email already in use", zap.String("email", user.Email))
		utils.WriteJSONError(w, http.StatusConflict)
		return
	}

	if err := c.svc.Update(r.Context(), uint(id), &user); err != nil {
		log.Error("Failed to update user", zap.Int("user_id", id), zap.Error(err))
		utils.WriteJSONError(w, http.StatusInternalServerError)
		return
	}

	log.Info("User updated successfully", zap.Uint("user_id", user.ID))
	json.NewEncoder(w).Encode(user)
}

func (c *UserController) DeleteUser(w http.ResponseWriter, r *http.Request) {
	log := logger.L(r.Context())
	idParam := chi.URLParam(r, "id")
	log.Info("DeleteUser handler invoked", zap.String("user_id_param", idParam))

	id, err := strconv.Atoi(idParam)
	if err != nil {
		log.Warn("Invalid user ID", zap.Error(err))
		utils.WriteJSONError(w, http.StatusBadRequest)
		return
	}

	if err := c.svc.Delete(r.Context(), uint(id)); err != nil {
		log.Error("Failed to delete user", zap.Int("user_id", id), zap.Error(err))
		utils.WriteJSONError(w, http.StatusInternalServerError)
		return
	}

	log.Info("User deleted successfully", zap.Int("user_id", id))
	w.WriteHeader(http.StatusNoContent)
}
