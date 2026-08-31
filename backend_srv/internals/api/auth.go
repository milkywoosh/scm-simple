package api

import (
	"encoding/json"
	"log"
	"net/http"

	// "github.com/hibiken/asynq"
	"github.com/julienschmidt/httprouter"
	// "github.com/milkywoosh/go_auth_pg_v1/worker"
	"scm-simple-luke.com/dir/internals"
	"scm-simple-luke.com/dir/internals/utils"
)

type reqLogin struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// note: oprek context
func (s *Server) Login(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {

	w.Header().Set("Content-Type", "application/json") // must be always first at handler
	var req reqLogin

	// read input from Body
	defer r.Body.Close()
	errDecode := json.NewDecoder(r.Body).Decode(&req)
	if errDecode != nil {
		errInfo := make(map[string]any)
		errInfo["message"] = errDecode.Error()
		internals.WriteErrorResponse(w, http.StatusConflict, errInfo)
		return
	}

	userInfo, err := s.service.GetInfoUser(r.Context(), req.Username)
	if err != nil {
		errInfo := make(map[string]any)
		errInfo["message"] = errDecode.Error()
		internals.WriteErrorResponse(w, http.StatusInternalServerError, errInfo)
		return
	}

	err = utils.ComparePassword(userInfo.HashedPassword.String, req.Password)
	if err != nil {
		errInfo := make(map[string]any)
		errInfo["message"] = err.Error()
		internals.WriteErrorResponse(w, http.StatusNotFound, errInfo)
		return
	}
	// w.WriteHeader(http.StatusAccepted)
	// json.NewEncoder(w).Encode("error")

	log.Printf("info username: %s", userInfo.Username.String)
	log.Printf("s.config.AccessTokenDuration: %s", s.config.AccessTokenDuration)

	accessToken, payload, err := s.token.CreateToken(userInfo.Username.String, "admin", s.config.AccessTokenDuration)
	if err != nil {
		errInfo := make(map[string]any)
		errInfo["message"] = err.Error()
		internals.WriteErrorResponse(w, http.StatusInternalServerError, errInfo)
		return
	}

	var responseData map[string]any = make(map[string]any)
	responseData["message"] = "sukses login"
	responseData["data"] = userInfo
	responseData["access_token"] = accessToken
	responseData["payload"] = payload
	internals.WriteResponse(w, http.StatusOK, responseData)

}

type reqRegistration struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Email    string `json:"email"`
	Fullname string `json:"full_name"`
}

// handler
func (s *Server) UserRegistration(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	w.Header().Set("Content-Type", "application/json") // must be always first at handler

	var registration reqRegistration

	defer r.Body.Close()

	err := json.NewDecoder(r.Body).Decode(&registration)
	if err != nil {
		errInfo := make(map[string]any)
		errInfo["message"] = err.Error()
		internals.WriteErrorResponse(w, http.StatusBadRequest, errInfo)
		return
	}

	hashedPW, err := utils.HashPassword(registration.Password)
	if err != nil {
		errInfo := make(map[string]any)
		errInfo["message"] = err.Error()
		internals.WriteErrorResponse(w, http.StatusInternalServerError, errInfo)
		return
	}

	// note: username param here, is actually same as "username string" AT FIRST ARGUMENT at this function UserRegistrationTx()
	// off
	/*
		afterCreate := func(username string) error {
			taskPayload := &worker.PayloadSendVerifyEmail{
				Username: username,
			}

			opts := []asynq.Option{
				asynq.MaxRetry(10),
				asynq.ProcessIn(10 * time.Second),
				asynq.Queue(worker.QueueCritical),
				asynq.Retention(10 * time.Minute),
			}

			return s.taskDistributor.DistributeTaskSendVerifyEmail(r.Context(), taskPayload, opts...)
		}
	*/

	// ctx context.Context, username string, password string, fullName string, email string, locationCode string
	err = s.service.UserRegistrationTx(r.Context(), registration.Username, hashedPW, registration.Fullname, registration.Email, nil)
	if err != nil {
		errInfo := make(map[string]any)
		errInfo["message"] = err.Error()
		internals.WriteErrorResponse(w, http.StatusExpectationFailed, errInfo)
		return
	}
	var responseData map[string]any = make(map[string]any)
	responseData["message"] = "sukses registrasi user baru"
	responseData["data"] = nil
	internals.WriteResponse(w, http.StatusOK, responseData)

}

func (s *Server) VerifyEmailSecretCode(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	w.Header().Set("Content-Type", "application/json") // must be always first at handler

	email := r.URL.Query().Get("email_id") // query param

	secret_code := r.URL.Query().Get("secret_code") // query param

	if email == "" {
		errInfo := make(map[string]any)
		errInfo["message"] = "param email tidak lengkap"
		internals.WriteErrorResponse(w, http.StatusBadRequest, errInfo)
		return
	}
	if secret_code == "" {
		errInfo := make(map[string]any)
		errInfo["message"] = "param secret_code tidak lengkap"
		internals.WriteErrorResponse(w, http.StatusBadRequest, errInfo)
		return
	}

	err := s.service.VerifyEmailSecretCodeTx(r.Context(), email, secret_code)
	if err != nil {
		log.Printf("err verify secret code: %v", err)
	} else {

		log.Printf("success verify secret code")
	}
}
