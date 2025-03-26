package handler

import (
	"encoding/json"
	"io"
	"net/http"
	"secretWall/internal/handler/dto"
	"secretWall/internal/service"
)

type authHandler struct {
	AuthService service.AuthService
}

type requestBody struct {
	IdentityToken string `json:"identity_token"`
}

func NewAuthHandler(authService service.AuthService, mux *http.ServeMux) {
	handler := &authHandler{
		AuthService: authService,
	}

	mux.HandleFunc("/auth/apple", handler.AppleLogin)
}

func (a *authHandler) AppleLogin(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		WriteHTTPError(w, err)
		return
	}

	var req requestBody
	if err := json.Unmarshal(body, &req); err != nil {
		WriteHTTPError(w, err)
		return
	}

	user, err := a.AuthService.AuthenticateWithApple(req.IdentityToken)

	if err != nil {
		WriteHTTPError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(dto.UserResponse{
		AccessToken: user.AccessToken,
		ID:          user.ID,
	})
	if err != nil {
		WriteHTTPError(w, err)
		return
	}
}
