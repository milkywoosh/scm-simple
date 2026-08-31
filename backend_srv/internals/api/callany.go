package api

import (
	"log"
	"net/http"

	"github.com/julienschmidt/httprouter"
	"scm-simple-luke.com/dir/internals"
	"scm-simple-luke.com/dir/internals/token"
)

type CalTestInitResponse struct {
	Message string `json:"message"`
	Data    struct {
		TransactionId  string `json:"transaction_id"`
		StatusApproval string `json:"status_approval"`
		WarehouseCode  string `json:"warehouse_code"`
	} `json:"data"`
}

// CallTestInit godoc
// @Summary      Test init for swagger docs
// @Description  Call Test Init after construct code for swagger init
// @Tags         callany
// @Accept       json
// @Produce      json
// @Param        transaction_id   	path       int  	 true  		"Transaction Id"
// @Success      200  {object}  CalTestInitResponse
// @Failure      400  {string}  ErrorResponse
// @Failure      404  {string}  ErrorResponse
// @Router       /api/test-init/{transaction_id} [get]
func CallTestInit(w http.ResponseWriter, req *http.Request, pathParam httprouter.Params) {

	log.Printf("cek CallTestInit 11")

	var responseData map[string]any = make(map[string]any)

	transactionId := pathParam.ByName("transaction_id")

	// note : result of payload must be infered to expected type
	payload, ok := req.Context().Value(payloadKey).(*token.Payload)
	if !ok {
		errInfo := make(map[string]any)
		errInfo["message"] = "Error, gagal mendapat payload"
		internals.WriteErrorResponse(w, http.StatusForbidden, errInfo)
		return
	}

	responseData["message"] = "oke"
	responseData["payload"] = payload
	responseData["data"] = struct {
		TransactionId  string
		StatusApproval string
		WarehouseCode  string
	}{
		TransactionId:  transactionId,
		StatusApproval: "approved",
		WarehouseCode:  "L0012",
	}

	internals.WriteResponse(w, http.StatusAccepted, responseData)
}
