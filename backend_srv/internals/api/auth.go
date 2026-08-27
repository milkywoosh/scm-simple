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
		internals.WriteErrorResponse(w, http.StatusConflict, errDecode.Error())
		return
	}

	userInfo, err := s.service.GetInfoUser(r.Context(), req.Username)
	if err != nil {
		internals.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	err = utils.ComparePassword(userInfo.HashedPassword.String, req.Password)
	if err != nil {
		internals.WriteErrorResponse(w, http.StatusNotFound, err.Error())
		return
	}
	// w.WriteHeader(http.StatusAccepted)
	// json.NewEncoder(w).Encode("error")

	var responseData map[string]any = make(map[string]any)
	responseData["message"] = "sukses login"
	responseData["data"] = userInfo
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
		internals.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	hashedPW, err := utils.HashPassword(registration.Password)
	if err != nil {
		internals.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
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
		internals.WriteErrorResponse(w, http.StatusExpectationFailed, err.Error())
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
		internals.WriteErrorResponse(w, http.StatusBadRequest, "param email tidak lengkap")
		return
	}
	if secret_code == "" {
		internals.WriteErrorResponse(w, http.StatusBadRequest, "param secret_code tidak lengkap")
		return
	}

	err := s.service.VerifyEmailSecretCodeTx(r.Context(), email, secret_code)
	if err != nil {
		log.Printf("err verify secret code: %v", err)
	} else {

		log.Printf("success verify secret code")
	}
}
