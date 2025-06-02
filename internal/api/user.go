package api

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
)

type UserHandler struct{}

type UserInfo struct {
	Sub           string   `json:"sub"`
	Email         string   `json:"email"`
	EmailVerified bool     `json:"email_verified"`
	Name          string   `json:"name"`
	Nickname      string   `json:"nickname"`
	Picture       string   `json:"picture"`
	Roles         []string `json:"/roles"`
}

func NewUserHandler() *UserHandler {
	return &UserHandler{}
}

func (user *UserHandler) GetUser(resWritter http.ResponseWriter, req *http.Request) {
	// 1. Get the bearer token from the authorization header
	authHeader := req.Header.Get("Authorization")
	if authHeader == "" {
		http.Error(resWritter, "Authorization header is required", http.StatusUnauthorized)
		return
	}

	// 2. Make request to Auth0 userinfo endpoint
	auth0Domain := os.Getenv("AUTH0_DOMAIN")
	if auth0Domain == "" {
		http.Error(resWritter, "Auth0 domain not configured", http.StatusInternalServerError)
		return
	}

	userInfoURL := fmt.Sprintf("https://%s/userinfo", auth0Domain)
	userInfoReq, err := http.NewRequest("GET", userInfoURL, nil)
	if err != nil {
		http.Error(resWritter, "Error creating request", http.StatusInternalServerError)
		return
	}

	userInfoReq.Header.Set("Authorization", authHeader)

	client := &http.Client{}
	resp, err := client.Do(userInfoReq)
	if err != nil {
		http.Error(resWritter, "Error fetching user info", http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		http.Error(resWritter, fmt.Sprintf("Error from Auth0: %s", string(body)), resp.StatusCode)
		return
	}

	// 3. Parse and return user info
	var userInfo UserInfo
	if err := json.NewDecoder(resp.Body).Decode(&userInfo); err != nil {
		http.Error(resWritter, "Error parsing user info", http.StatusInternalServerError)
		return
	}

	resWritter.Header().Set("Content-Type", "application/json")

	// Return user info as JSON
	json.NewEncoder(resWritter).Encode(userInfo)
}
