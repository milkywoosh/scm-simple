package api

import (
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
