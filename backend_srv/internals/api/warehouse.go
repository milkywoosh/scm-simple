package api

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/julienschmidt/httprouter"
	"scm-simple-luke.com/dir/internals"
	"scm-simple-luke.com/dir/internals/domain"
)

type LocationSearchResponse struct {
	Message string               `json:"message"`
	Data    []domain.LocationRow `json:"data"`
}

// GetLocationByLocationCode godoc
// @Summary      Search location by location code
// @Description  Search location by location code, customer, technician, warehouse
// @Tags         location warehouse customer technician
// @Accept       json
// @Produce      json
// @Param        q              query  string  false  "name search by q"  Format(email)
// @Param        location_code  path  string  false  "Location code"
// @Success      200  {object}  map[string]any
// @Failure      400  {string}  ErrorResponse
// @Router       /api/warehouse/:location_code [get]
func (s *Server) GetLocationByLocationCode(w http.ResponseWriter, req *http.Request, pathParam httprouter.Params) {

	locationCode := pathParam.ByName("location_code")

	whInfo, err := s.service.WarehouseInfo(req.Context(), locationCode)
	if err != nil {
		log.Println(err)

		internals.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	var resp map[string]any = make(map[string]any)
	resp["Message"] = "lokasi ditemukan"
	resp["Data"] = whInfo

	internals.WriteResponse(w, http.StatusAccepted, resp)
}

type CreateDraftWarehouseTransferParams struct {
	LocationOrigin      string `json:"location_origin"`
	LocationDestination string `json:"location_destination"`
}

// CreateDraftWarehouseTransfer godoc
// @Summary      Create New Draft Transaction
// @Description  Create New Draft Transaction Warehouse to Warehouse
// @Tags         warehouse_service item_transfer
// @Accept       json
// @Produce      json
// @Param		 request        body      CreateDraftWarehouseTransferParams  true  "draft transaction transfer payload"
// @Success      200  {object}  map[string]any
// @Failure      400  {string}  ErrorResponse
// @Router       /api/v1/warehouse/create-draft-transfer [post]
func (s *Server) CreateDraftWarehouseTransfer(w http.ResponseWriter, req *http.Request, pathParam httprouter.Params) {

	ctx := req.Context()

	var ReqBody CreateDraftWarehouseTransferParams
	DataResp := make(map[string]any)

	// read input from Body
	errDecode := json.NewDecoder(req.Body).Decode(&ReqBody)
	if errDecode != nil {
		DataResp["message"] = "error parsing request body"
		internals.WriteResponse(w, http.StatusBadRequest, DataResp)
		return
	}

	NewDraft, err := s.service.WarehouseService.CreateDraftTx(ctx, ReqBody.LocationOrigin, ReqBody.LocationDestination)
	if err != nil {
		DataResp["message"] = err.Error()
		internals.WriteResponse(w, http.StatusConflict, DataResp)
		return
	}
	DataResp["message"] = "create draft berhasil"
	DataResp["data"] = NewDraft
	internals.WriteResponse(w, http.StatusAccepted, DataResp)
	return

}

type InputItemWarehouseTransferParams struct {
	Identifier string `json:"identifier"`
}

// InputItemWarehouseTransfer godoc
// @Summary      Input each item one by one(serial_number) in Transaction
// @Description  User scan identifier for each item that will be delivered as ALLOCATED
// @Tags         warehouse_service item_transfer
// @Accept       json
// @Produce      json
// @Param		 request        body      InputItemWarehouseTransferParams  true  "item transfer payload"
// @Param		 request        path      string  true  "transaction number identifier"
// @Success      200  {object}  map[string]any
// @Failure      400  {string}  ErrorResponse
// @Router       /api/v1/warehouse/input-item/:transaction_number [post]
func (s *Server) InputItemWarehouseTransfer(w http.ResponseWriter, req *http.Request, pathParam httprouter.Params) {

	ctx := req.Context()

	var ReqBody InputItemWarehouseTransferParams
	TransactionNumber := pathParam.ByName("transaction_number")

	DataResp := make(map[string]any)

	// read input from Body
	errDecode := json.NewDecoder(req.Body).Decode(&ReqBody)
	if errDecode != nil {
		DataResp["message"] = "error parsing request body"
		internals.WriteResponse(w, http.StatusBadRequest, DataResp)
		return
	}

	err := s.service.WarehouseService.AllocateItem(ctx, TransactionNumber, ReqBody.Identifier)
	if err != nil {
		internals.WriteErrorResponse(w, http.StatusConflict, err.Error())
		return
	}

	DataResp["message"] = "return okay"
	internals.WriteResponse(w, http.StatusAccepted, DataResp)
	return
}

type DisallocateItemWarehouseTransferParams struct {
	Identifier string `json:"identifier"`
}

// DisallocateItemWarehouseTransfer godoc
// @Summary      Takeout each item one by one(serial_number) in Transaction
// @Description  User scan identifier for each item that will be takenout as set to be
// @Tags         warehouse_service item_transfer
// @Accept       json
// @Produce      json
// @Param		 request        body      DisallocateItemWarehouseTransferParams  true  "item transfer payload"
// @Param		 request        path      string  true  "transaction number identifier"
// @Success      200  {object}  map[string]any
// @Failure      400  {string}  ErrorResponse
// @Router       /api/v1/warehouse/disallocate-item/:transaction_number [post]
func (s *Server) DisallocateItemWarehouseTransfer(w http.ResponseWriter, req *http.Request, pathParam httprouter.Params) {

	ctx := req.Context()

	var ReqBody DisallocateItemWarehouseTransferParams
	TransactionNumber := pathParam.ByName("transaction_number")

	DataResp := make(map[string]any)

	// read input from Body
	errDecode := json.NewDecoder(req.Body).Decode(&ReqBody)
	if errDecode != nil {
		DataResp["message"] = "error parsing request body"
		internals.WriteResponse(w, http.StatusBadRequest, DataResp)
		return
	}

	err := s.service.WarehouseService.DisAllocateItem(ctx, TransactionNumber, ReqBody.Identifier)
	if err != nil {
		internals.WriteErrorResponse(w, http.StatusConflict, err.Error())
		return
	}

	DataResp["message"] = "return okay"
	internals.WriteResponse(w, http.StatusAccepted, DataResp)
	return
}

// SetStatusTransaction godoc
// @Summary      set update status Transaction
// @Description  Authorized User can update transaction process ("submit", "reject", "cancel", "approve")
// @Tags         warehouse_service item_transfer
// @Accept       json
// @Produce      json
// @Param		 request        path      string  true  "transaction number identifier"
// @Param		 request        path      string  true  "status to update"
// @Success      200  {object}  map[string]any
// @Failure      400  {string}  ErrorResponse
// @Router       /api/v1/warehouse/set-status/:transaction_number/:status [put]
func (s *Server) SetStatusTransaction(w http.ResponseWriter, req *http.Request, pathParam httprouter.Params) {

	ctx := req.Context()
	TransactionNumber := pathParam.ByName("transaction_number")
	SetStatus := pathParam.ByName("status")

	if SetStatus == "submit" {
		err := s.service.WarehouseService.SetSubmit(ctx, TransactionNumber)
		if err != nil {
			internals.WriteErrorResponse(w, http.StatusConflict, err.Error())
			return
		}

		RespData := make(map[string]any)
		RespData["message"] = "update status berhasil"
		RespData["transaction_number"] = TransactionNumber

		internals.WriteResponse(w, http.StatusAccepted, RespData)
		return
	} else if SetStatus == "cancel" {

		err := s.service.WarehouseService.SetCancel(ctx, TransactionNumber)
		if err != nil {
			internals.WriteErrorResponse(w, http.StatusConflict, err.Error())
			return
		}

		RespData := make(map[string]any)
		RespData["message"] = "update status berhasil"
		RespData["transaction_number"] = TransactionNumber

		internals.WriteResponse(w, http.StatusAccepted, RespData)
		return
	} else if SetStatus == "reject" {
		err := s.service.WarehouseService.SetReject(ctx, TransactionNumber)
		if err != nil {
			internals.WriteErrorResponse(w, http.StatusConflict, err.Error())
			return
		}

		RespData := make(map[string]any)
		RespData["message"] = "update status berhasil"
		RespData["transaction_number"] = TransactionNumber

		internals.WriteResponse(w, http.StatusAccepted, RespData)
		return

	} else if SetStatus == "approve" {
		err := s.service.WarehouseService.SetApprove(ctx, TransactionNumber)
		if err != nil {
			internals.WriteErrorResponse(w, http.StatusConflict, err.Error())
			return
		}

		RespData := make(map[string]any)
		RespData["message"] = "update status berhasil"
		RespData["transaction_number"] = TransactionNumber

		internals.WriteResponse(w, http.StatusAccepted, RespData)
		return
	} else {
		internals.WriteErrorResponse(w, http.StatusBadRequest, "Proses ini tidak dapat dilakukan. Request tersedia hanya ada (submit, cancel, reject, approve)")
		return

	}

}

// ListItemsSent godoc
// @Summary      get data transaction number and list of outbound items from warehouse
// @Description  Authorized User can receive this info, limit by authorization token.
// @Tags         warehouse_service item_transfer
// @Accept       json
// @Produce      json
// @Param		 request        path      string  true  "transaction number identifier"
// @Success      200  {object}  map[string]any
// @Failure      400  {string}  ErrorResponse
// @Router       /api/v1/warehouse/list-items/:transaction_number [get]
func (s *Server) ListItemsSent(w http.ResponseWriter, req *http.Request, pathParam httprouter.Params) {

	ctx := req.Context()
	TransactionNumber := pathParam.ByName("transaction_number")

	listItems, err := s.service.WarehouseService.ListItemsSent(ctx, TransactionNumber)
	if err != nil {
		internals.WriteErrorResponse(w, http.StatusConflict, err.Error())
		return
	}

	RespData := make(map[string]any)
	RespData["data"] = listItems
	RespData["message"] = "list of warehouse outbound"
	internals.WriteResponse(w, http.StatusAccepted, RespData)
}
