package api

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"

	jwtmiddleware "github.com/auth0/go-jwt-middleware/v2"
	"github.com/auth0/go-jwt-middleware/v2/validator"
	"github.com/leonardosantos2/GoUserServer/internal/middlewares"
)

type UserHandler struct{}

type RolesAndMetadata struct {
	Roles     []string `json:"/roles"`
	Merchants []string `json:"/merchants"`
}

type UserInfo struct {
	RolesAndMetadata
	Email         string `json:"email"`
	EmailVerified bool   `json:"email_verified"`
	Name          string `json:"name"`
	Nickname      string `json:"nickname"`
	Picture       string `json:"picture"`
	Sub           string `json:"sub"`
}

func NewUserHandler() *UserHandler {
	return &UserHandler{}
}

func (user *UserHandler) GetUserRolesAndAppMetadata(resWritter http.ResponseWriter, req *http.Request) {
	// Get the claims from the request context
	claims, ok := req.Context().Value(jwtmiddleware.ContextKey{}).(*validator.ValidatedClaims)
	if !ok {
		http.Error(resWritter, "Invalid token claims", http.StatusUnauthorized)
		return
	}

	// Cast the custom claims to our CustomClaims type
	customClaims, ok := claims.CustomClaims.(*middleware.CustomClaims)
	if !ok {
		http.Error(resWritter, "Invalid custom claims", http.StatusUnauthorized)
		return
	}

	// Create RolesAndMetadata from the claims
	rolesAndMetadata := RolesAndMetadata{
		Roles:     customClaims.Roles,
		Merchants: customClaims.Merchants,
	}

	resWritter.Header().Set("Content-Type", "application/json")
	json.NewEncoder(resWritter).Encode(rolesAndMetadata)
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
