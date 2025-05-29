package api

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)

type UserHandler struct{}

func NewUserHandler() *UserHandler {
	return &UserHandler{}
}

func (user *UserHandler) GetUser(resWritter http.ResponseWriter, req *http.Request) {
	paramsUserId := chi.URLParam(req, "id")

	if paramsUserId == "" {
		http.Error(resWritter, "User ID is required", http.StatusBadRequest)
		return
	}

	userId, err := strconv.ParseInt(paramsUserId, 10, 64)
	if err != nil {
		http.Error(resWritter, "Invalid user ID", http.StatusBadRequest)
		return
	}

	resWritter.Header().Set("Content-Type", "application/json")
	response := struct {
		UserId int64 `json:"userId"`
	}{
		UserId: userId,
	}

	json.NewEncoder(resWritter).Encode(response)
}
