package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"scm-simple-luke.com/dir/internals"
)

type renewAccessTokenRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}

type renewAccessTokenResponse struct {
	AccessToken          string    `json:"access_token"`
	AccessTokenExpiresAt time.Time `json:"access_token_expires_at"`
}

func (server *Server) renewAccessToken(w http.ResponseWriter, req *http.Request) {
	var ReqBody renewAccessTokenRequest
	ctx := req.Context()
	var DataResp map[string]any = make(map[string]any)

	errDecode := json.NewDecoder(req.Body).Decode(&ReqBody)
	if errDecode != nil {
		DataResp["message"] = "error parsing request body"
		internals.WriteResponse(w, http.StatusBadRequest, DataResp)
		return
	}

	refreshPayload, err := server.token.VerifyToken(ReqBody.RefreshToken)
	if err != nil {
		DataResp["message"] = err.Error()
		internals.WriteResponse(w, http.StatusUnauthorized, DataResp)
		return
	}

	// Payload.ID for Session ID table
	session, err := server.service.GetSession(ctx, refreshPayload.ID)
	if err != nil {
		DataResp["message"] = err.Error()
		internals.WriteResponse(w, http.StatusUnauthorized, DataResp)
		return
	}

	// kapan ini diblock??
	if session.IsBlocked {
		err := fmt.Errorf("blocked session")
		DataResp["message"] = err
		internals.WriteResponse(w, http.StatusUnauthorized, DataResp)
		return
	}

	if session.Username != refreshPayload.Username {
		err := fmt.Errorf("incorrect session user")
		DataResp["message"] = err
		internals.WriteResponse(w, http.StatusUnauthorized, DataResp)
		return
	}

	if session.RefreshToken != ReqBody.RefreshToken {
		err := fmt.Errorf("mismatched session token")
		DataResp["message"] = err
		internals.WriteResponse(w, http.StatusUnauthorized, DataResp)
		return
	}

	// token refresh expired ga? if yes, wajib relogin!
	if time.Now().After(session.ExpiresAt) {
		err := fmt.Errorf("expired session")
		DataResp["message"] = err
		internals.WriteResponse(w, http.StatusUnauthorized, DataResp)
		return
	}

	accessToken, accessPayload, err := server.token.CreateToken(
		refreshPayload.Username,
		refreshPayload.Role,
		server.config.AccessTokenDuration,
	)
	if err != nil {
		DataResp["message"] = err
		internals.WriteResponse(w, http.StatusUnauthorized, DataResp)
		return
	}

	// stored at Web Storage (sessionStorage, localStorage)
	rsp := renewAccessTokenResponse{
		AccessToken:          accessToken,
		AccessTokenExpiresAt: accessPayload.ExpiredAt,
	}
	DataResp["message"] = "renewal token berhasil"
	DataResp["data"] = rsp
	internals.WriteResponse(w, http.StatusAccepted, DataResp)
	
}
