package api

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
	"scm-simple-luke.com/dir/internals"
)

// UploadItems godoc
// @Summary      Upload item (multi) for warehouse
// @Description  Upload bulk items with csv file
// @Tags         warehouse_service item_transfer upload_item
// @Accept       json
// @Produce      json
// @Param		 request        body      UploadItemParams  true  "upload item bulk payload"
// @Success      200  {object}  map[string]any
// @Failure      400  {string}  ErrorResponse
// @Router       /api/v1/warehouse/upload/items [post]

func (s *Server) UploadItems(w http.ResponseWriter, req *http.Request, pathParam httprouter.Params) {

	DataResp := make(map[string]any)
	DataResp["message"] = "ok"
	internals.WriteResponse(w, http.StatusAccepted, DataResp)
}
