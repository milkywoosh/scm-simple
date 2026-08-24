package api

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/julienschmidt/httprouter"
	"scm-simple-luke.com/dir/internals"
)

func (s *Server) GetListAvailableOutbound(w http.ResponseWriter, req *http.Request, pathParam httprouter.Params) {
	// get all nomor transfer outbound yg menuju warehouse inbound
}

type CreateDraftWarehouseInboundParams struct {
	OutboundNumber string `json:"outbound_number"`
}

// CreateDraftWarehouseInbound godoc
// @Summary      Create draft for warehouse inbound
// @Description  Warehouse inbound receiving all incoming items from warehouse origin
// @Tags         warehouse_service item_transfer
// @Accept       json
// @Produce      json
// @Param		 request        body      CreateDraftWarehouseInboundParams  true  "draft transaction transfer payload"
// @Success      200  {object}  map[string]any
// @Failure      400  {string}  ErrorResponse
// @Router       /api/v1/warehouse/inbound/create-draft [post]
func (s *Server) CreateDraftWarehouseInbound(w http.ResponseWriter, req *http.Request, pathParam httprouter.Params) {
	ctx := req.Context()

	var ReqBody CreateDraftWarehouseInboundParams
	DataResp := make(map[string]any)

	// read input from Body
	errDecode := json.NewDecoder(req.Body).Decode(&ReqBody)
	if errDecode != nil {
		DataResp["message"] = "error parsing request body"
		internals.WriteResponse(w, http.StatusBadRequest, DataResp)
		return
	}

	NewDraft, err := s.service.WarehouseService.CreateDraftInboundTx(ctx, ReqBody.OutboundNumber)
	if err != nil {
		DataResp["message"] = err.Error()
		internals.WriteResponse(w, http.StatusConflict, DataResp)
		return
	}
	DataResp["message"] = "create draft inbound berhasil"
	DataResp["data"] = NewDraft
	internals.WriteResponse(w, http.StatusAccepted, DataResp)
	return

}

// kondisi nomor transfer outbound sudah diterima di gudang tujuan

// DurationTransfer godoc
// @Summary      Check duration transfer between warehouse
// @Description  Subtract time between arrived_at and delivery_at timestamptz
// @Tags         warehouse_service item_transfer
// @Accept       json
// @Produce      json
// @Param		 request        path   string  true  "outbound number transaction"
// @Param		 request        path   string  true  "inbound number transaction"
// @Success      200  {object}  map[string]any
// @Failure      400  {string}  ErrorResponse
// @Router       /api/v1/warehouse/transfer-duration/:outbound_number/:inbound_number [get]
func (s *Server) DurationTransfer(w http.ResponseWriter, req *http.Request, pathParam httprouter.Params) {

	ctx := req.Context()
	outbound_number := pathParam.ByName("outbound_number")
	inbound_number := pathParam.ByName("inbound_number")

	DurationInfo, err := s.service.InfoTransferDuration(ctx, outbound_number, inbound_number)
	if err != nil {
		log.Println(err)

		internals.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	var resp map[string]any = make(map[string]any)
	resp["message"] = "informasi durasi berhasil didapat"
	resp["info_durasi"] = DurationInfo

	internals.WriteResponse(w, http.StatusAccepted, resp)
}

type InputItemWarehouseInboundParams struct {
	Identifier string `json:"identifier"`
}

// InputItemWarehouseInbound godoc
// @Summary      Input each item one by one(serial_number) in Inbound Transaction
// @Description  User scan identifier for each item that will be delivered as AVAILABLE
// @Tags         warehouse_service item_transfer
// @Accept       json
// @Produce      json
// @Param		 request        body      InputItemWarehouseInboundParams  true  "receiving item transfer payload"
// @Param		 request        path      string  true  "transaction number inbound identifier"
// @Success      200  {object}  map[string]any
// @Failure      400  {string}  ErrorResponse
// @Router       /api/v1/warehouse/inbound/input-item/:transaction_number [post]
func (s *Server) InputItemWarehouseInbound(w http.ResponseWriter, req *http.Request, pathParam httprouter.Params) {

	ctx := req.Context()

	var ReqBody InputItemWarehouseInboundParams
	TransactionNumber := pathParam.ByName("transaction_number")

	DataResp := make(map[string]any)

	// read input from Body
	errDecode := json.NewDecoder(req.Body).Decode(&ReqBody)
	if errDecode != nil {
		DataResp["message"] = "error parsing request body"
		internals.WriteResponse(w, http.StatusBadRequest, DataResp)
		return
	}

	err := s.service.WarehouseService.AllocateInboundItem(ctx, TransactionNumber, ReqBody.Identifier)
	if err != nil {
		internals.WriteErrorResponse(w, http.StatusConflict, err.Error())
		return
	}

	DataResp["message"] = "return okay"
	internals.WriteResponse(w, http.StatusAccepted, DataResp)
	return
}

// ListItemsReceived godoc
// @Summary      get data transaction number and list of inbound items from warehouse
// @Description  Authorized User can receive this info, limit by authorization token.
// @Tags         warehouse_service item_transfer
// @Accept       json
// @Produce      json
// @Param		 request        path      string  true  "transaction number identifier"
// @Success      200  {object}  map[string]any
// @Failure      400  {string}  ErrorResponse
// @Router       /api/v1/warehouse/inbound/list-items/:transaction_number [get]
func (s *Server) ListItemsReceived(w http.ResponseWriter, req *http.Request, pathParam httprouter.Params) {

	ctx := req.Context()
	TransactionNumber := pathParam.ByName("transaction_number")

	listItems, err := s.service.WarehouseService.ListItemsReceived(ctx, TransactionNumber)
	if err != nil {
		internals.WriteErrorResponse(w, http.StatusConflict, err.Error())
		return
	}

	RespData := make(map[string]any)
	RespData["data"] = listItems
	RespData["message"] = "list of warehouse inbound"
	internals.WriteResponse(w, http.StatusAccepted, RespData)
}

// SetStatusInboundTransaction godoc
// @Summary      set update status Transaction inbound
// @Description  Authorized User can update transaction process ("submit", "reject", "cancel", "approve")
// @Tags         warehouse_service item_transfer
// @Accept       json
// @Produce      json
// @Param		 request        path      string  true  "transaction number identifier"
// @Param		 request        path      string  true  "status to update"
// @Success      200  {object}  map[string]any
// @Failure      400  {string}  ErrorResponse
// @Router       /api/v1/warehouse/inbound/set-status/:transaction_number/:status [put]
func (s *Server) SetStatusInboundTransaction(w http.ResponseWriter, req *http.Request, pathParam httprouter.Params) {

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

		err := s.service.WarehouseService.SetCancelInbound(ctx, TransactionNumber)
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
		err := s.service.WarehouseService.SetApproveInbound(ctx, TransactionNumber)
		if err != nil {
			internals.WriteErrorResponse(w, http.StatusConflict, err.Error())
			return
		}

		RespData := make(map[string]any)
		RespData["message"] = "update status warehouse inbound berhasil"
		RespData["transaction_number"] = TransactionNumber

		internals.WriteResponse(w, http.StatusAccepted, RespData)
		return
	} else {
		internals.WriteErrorResponse(w, http.StatusBadRequest, "Proses ini tidak dapat dilakukan. Request tersedia hanya ada (submit, cancel, reject, approve)")
		return

	}

}
