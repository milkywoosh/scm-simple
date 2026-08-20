package api

import (
	"context"
	_ "context"
	"encoding/json"
	_ "encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/julienschmidt/httprouter"
	httpSwagger "github.com/swaggo/http-swagger"
	"scm-simple-luke.com/dir/internals"
	"scm-simple-luke.com/dir/internals/services"
)

type Server struct {
	srvInit         *http.Server
	route           *httprouter.Router
	service         *services.Services
	token           string
	taskDistributor any //worker.TaskDistributor
}

type ErrorResponse string

func NewServer(services *services.Services, taskDistributor any) (*Server, error) {

	// need token maker here JWT or Paseto

	newHttpRouter := httprouter.New()

	addr := fmt.Sprintf("%s:%s", os.Getenv("ADDR"), os.Getenv("PORT"))

	srv := &http.Server{
		Handler: newHttpRouter,
		Addr:    addr,
	}

	// note: watch out about http and route, both should be in sync
	s := &Server{
		srvInit:         srv,
		route:           newHttpRouter, // httprouter.New()
		service:         services,
		taskDistributor: taskDistributor,
	}

	s.setupRoutes()
	return s, nil
}

type TestReqBody struct {
	Param1 string `json:"param1"`
	Param2 string `json:"param2"`
}

func (s *Server) setupRoutes() {
	// define each routes here, must satisfy interface httprouter.Handle
	s.route.GET("/health-check", func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {

		// how to parse body request from http.Request?
		/*
			says i have
			{
				"username": string
				"password": string
			}
		*/

		var ReqBody TestReqBody
		data := make(map[string]any)

		// read input from Body
		errDecode := json.NewDecoder(r.Body).Decode(&ReqBody)
		if errDecode != nil {
			data["message"] = "error parsing request body"
			internals.WriteResponse(w, http.StatusBadRequest, data)
			return
		}

		data["success"] = true
		data["message"] = "health check okay"
		data["data"] = ReqBody

		internals.WriteResponse(w, http.StatusOK, data)
	})

	s.swaggerRoute()
	s.callInit()
	s.warehouseRoutes()

	// s.userRoutes()
	// s.stockTransferRoutes() // each of stock_transfer route
}

func (s *Server) swaggerRoute() {
	s.route.GET("/swagger/*any", func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		httpSwagger.WrapHandler(w, r)
	})
}

func (s *Server) callInit() {
	s.route.GET("/api/v1/test-init/:transaction_id", CallTestInit)
}

func (s *Server) warehouseRoutes() {
	s.route.GET("/api/v1/warehouse/:location_code", s.GetLocationByLocationCode)
	s.route.POST("/api/v1/warehouse/create-draft-transfer", s.CreateDraftWarehouseTransfer)
}

func (s *Server) userRoutes() {

	s.route.POST("/api/user/v1/registration", nil)
	// s.route.POST("/api/user/v1/registration", s.UserRegistration)
	// s.route.GET("/api/user/v1/login", s.Login)
	// s.route.POST("/api/user/v1/verify-secret-code", s.VerifyEmailSecretCode)
}

func (s *Server) stockTransferRoutes() {
	s.route.POST("/stock-transfer/v1/create", nil)
	s.route.GET("/stock-transfer/v1/header", nil)
	s.route.PUT("/stock-transfer/v1/submit", nil)
	s.route.PUT("/stock-transfer/v1/reject", nil)
	s.route.PUT("/stock-transfer/v1/approve", nil)
	s.route.PUT("/stock-transfer/v1/cancel", nil)
}

func (s *Server) Start(addr string) error {
	return s.srvInit.ListenAndServe()
}

func (s *Server) Shutdown(ctx context.Context) error {
	return s.srvInit.Shutdown(ctx)
}
