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
	LocationOrigin      string `json:"location_origin"`
	LocationDestination string `json:"location_destination"`
	OutboundNumber      string `json:"outbound_number"`
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
// @Router       /api/v1/warehouse/create-draft-inbound [post]
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

	NewDraft, err := s.service.WarehouseService.CreateDraftInboundTx(ctx, ReqBody.LocationOrigin, ReqBody.LocationDestination, ReqBody.OutboundNumber)
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
