package api

import (
	"context"
	_ "context"
	"encoding/json"
	_ "encoding/json"
	"fmt"
	"net/http"

	"github.com/julienschmidt/httprouter"
	httpSwagger "github.com/swaggo/http-swagger"
	"scm-simple-luke.com/dir/internals"
	"scm-simple-luke.com/dir/internals/services"
	"scm-simple-luke.com/dir/internals/token"
	"scm-simple-luke.com/dir/internals/utils"
)

type Server struct {
	config          utils.Config
	srvInit         *http.Server
	route           *httprouter.Router
	service         *services.Services
	token           token.TokenMaker
	taskDistributor any //worker.TaskDistributor
}

type ErrorResponse string

func NewServer(cfg utils.Config, services *services.Services, token token.TokenMaker, taskDistributor any) (*Server, error) {

	// need token maker here JWT or Paseto

	newHttpRouter := httprouter.New()

	addr := fmt.Sprintf("%s:%s", cfg.Addr, cfg.Port)

	srv := &http.Server{
		Handler: newHttpRouter,
		Addr:    addr,
	}

	// note: watch out about http and route, both should be in sync
	s := &Server{
		config:          cfg,
		srvInit:         srv,
		route:           newHttpRouter, // httprouter.New()
		token:           token,
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
	s.userRoutes()

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
	s.route.GET("/api/v1/warehouse/info/:location_code", s.GetLocationByLocationCode)
	s.route.GET("/api/v1/warehouse/transfer-duration/:outbound_number/:inbound_number", s.DurationTransfer)

	s.route.POST("/api/v1/warehouse/outbound/create-draft", s.CreateDraftWarehouseOutbound)
	s.route.PUT("/api/v1/warehouse/outbound/set-status/:transaction_number/:status", s.SetStatusTransaction)
	s.route.POST("/api/v1/warehouse/outbound/input-item/:transaction_number", s.InputItemWarehouseTransfer)
	s.route.POST("/api/v1/warehouse/outbound/disallocate-item/:transaction_number", s.DisallocateItemWarehouseTransfer)
	s.route.GET("/api/v1/warehouse/outbound/list-items/:transaction_number", s.ListItemsSent)

	s.route.POST("/api/v1/warehouse/inbound/create-draft", s.CreateDraftWarehouseInbound)
	s.route.PUT("/api/v1/warehouse/inbound/set-status/:transaction_number/:status", s.SetStatusInboundTransaction)
	s.route.POST("/api/v1/warehouse/inbound/input-item/:transaction_number", s.InputItemWarehouseInbound)
	s.route.GET("/api/v1/warehouse/inbound/list-items/:transaction_number", s.ListItemsReceived)

}

func (s *Server) userRoutes() {
	s.route.POST("/api/v1/user/signin", s.Login)
	s.route.POST("/api/v1/user/registration", s.UserRegistration)
	// s.route.POST("/api/user/v1/verify-secret-code", s.VerifyEmailSecretCode)
}

func (s *Server) stockTransferRoutes() {}

func (s *Server) Start(addr string) error {
	return s.srvInit.ListenAndServe()
}

func (s *Server) Shutdown(ctx context.Context) error {
	return s.srvInit.Shutdown(ctx)
}
