package api

import (
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/google/uuid"
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

	ctx := req.Context()
	// nanti dapet dari token
	tokenInfo := struct {
		Bucket, Location string
	}{
		"lukefile", "PGC001",
	}

	err := req.ParseMultipartForm(5 << 20) // max 5Mb
	if err != nil {
		DataResp["message"] = err.Error()
		internals.WriteErrorResponse(w, http.StatusBadRequest, DataResp)
		return
	}

	file, header, err := req.FormFile("file-upload1")
	if err != nil {
		DataResp["message"] = err.Error()
		internals.WriteResponse(w, http.StatusBadRequest, DataResp)
		return
	}
	defer file.Close()

	fileName := header.Filename
	headerContentType := header.Header.Get("Content-Type")
	// size := header.Size
	// field := csvData.Columns
	// rows := csvData.Rows

	if headerContentType != "text/plain" && headerContentType != "text/csv" {
		DataResp["message"] = fmt.Errorf("error, tipe file bukan csv")
		internals.WriteResponse(w, http.StatusBadRequest, DataResp)
		return
	}

	// upload aja, nanti diproses di consumer
	// output,
	_, fileDir, err := s.s3.PutObject(ctx, tokenInfo.Bucket, tokenInfo.Location, fileName, file)
	if err != nil {
		DataResp["message"] = err.Error()
		internals.WriteResponse(w, http.StatusConflict, DataResp)
		return
	}

	// data csv upload ke S3 dulu, baru nanti diproses di rabbitMQ
	DataResp["message"] = "ok"
	DataResp["file_name"] = fileName
	DataResp["file_directory"] = fileDir

	jobId := uuid.NewString()
	jobType := "upload_items"

	// tipe payload harus sesuai antara publish dan consume
	payloadPublish := struct {
		Username  string
		Warehouse string
		FileDir string
	}{
		Username:  tokenInfo.Bucket,
		Warehouse: tokenInfo.Location,
		FileDir: fileDir,

	}

	err = s.publisherPool.PublishV2(ctx, "upload.exchange", "upload.create", jobId, jobType, payloadPublish)
	if err != nil {
		DataResp["message"] = err.Error()
		internals.WriteResponse(w, http.StatusConflict, DataResp)
		return
	}

	internals.WriteResponse(w, http.StatusAccepted, DataResp)
}

func (s *Server) UploadLogReadDirectly(w http.ResponseWriter, req *http.Request, pathParam httprouter.Params) {
	// read file dari formFile insomnia
	// read dengan csv.NewReader(file)

	DataResp := make(map[string]any)

	log.Printf("log UploadReadDirectly")

	file, _, err := req.FormFile("file-upload1")
	if err != nil {
		DataResp["message"] = err.Error()
		internals.WriteErrorResponse(w, http.StatusBadRequest, DataResp)
		return
	}

	defer file.Close()

	// read dengan csv.NewReader(file)
	// csv.NewReader() baca file 1 by 11
	reader := csv.NewReader(file)
	reader.Comma = ';'

	startRow := 0
	for {
		rec, err := reader.Read()
		log.Printf("each rec: %d", startRow)
		if startRow == 0 {
			startRow++
			continue
		}
		if err == io.EOF {
			break
		}
		log.Printf("each rec: %v", rec)

		if err != nil {
			DataResp["message"] = err.Error()
			internals.WriteErrorResponse(w, http.StatusBadRequest, DataResp)
			return
		}

		startRow++
	}

	// var rm json.RawMessage = json.RawMessage{}

	DataResp["message"] = "ok"
	internals.WriteResponse(w, http.StatusAccepted, DataResp)
}
