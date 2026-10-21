package services

import (
	"context"
	"encoding/csv"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"strings"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"scm-simple-luke.com/dir/internals/db"
	"scm-simple-luke.com/dir/internals/domain"
	msgqueue "scm-simple-luke.com/dir/internals/msg_queue"
	"scm-simple-luke.com/dir/internals/storage"
)

type UploadItemProcessor struct {
	itemRepo domain.ItemRepository
	storage  storage.ObjectS3
}

// klo mau ada concurrent utk proses data lgi maka buat struct baru
// implement `MessageProcessor`
func (u *UploadItemProcessor) ProcessData(ctx context.Context, envelopeMsg msgqueue.Envelope) error {

	// ** envelopeMsg msgqueue.Envelope ==> sudah di-Unmarshal

	// tipe payload harus sesuai antara publish dan consume
	type payloadPublish struct {
		Username  string
		Warehouse string
		FileDir   string
	}

	startTime := time.Now()

	var p payloadPublish
	err := json.Unmarshal(envelopeMsg.Payload, &p)
	if err != nil {
		return err
	}

	bucketKeyFilename := strings.Split(p.FileDir, "/")
	if len(bucketKeyFilename) != 3 {
		return fmt.Errorf("error info filedir storage, saat ini wajib 3 level: %s", bucketKeyFilename)
	}

	bucket := bucketKeyFilename[0]
	key := bucketKeyFilename[1]
	fileName := bucketKeyFilename[2]

	output, err := u.storage.GetObject(ctx, bucket, key, fileName)
	if err != nil {
		return err
	}

	defer output.Body.Close()

	r := csv.NewReader(output.Body)
	r.ReuseRecord = true

	// loop for first header
	_, err = r.Read()
	if err != nil {
		if errors.Is(err, io.EOF) {
			return fmt.Errorf("data kosong")
		}
		return err
	}

	// loop after header
	for {
		rec, err := r.Read()
		if err != nil {
			if errors.Is(err, io.EOF) {
				break
			}
		}
		log.Printf("each record csv: %v", rec)
		// log.Printf("cus ProcessData: %v", DataPayload)
		separated := strings.Split(rec[0], ";")
		eachRow := domain.ItemInfo{
			Serial_number: separated[1],
			Product_code:  separated[2],
		}
		err = u.itemRepo.UploadNewItem(ctx, eachRow)
		if err != nil {
			log.Printf("err UploadNewItem repo: %s", err.Error())
		}
	}

	durationLong := time.Since(startTime)
	log.Printf("duration after looping: %s", durationLong.String())

	return nil
}

func NewProcUpload(ctx context.Context, dbPool *pgxpool.Pool, storage storage.ObjectS3) (*UploadItemProcessor, error) {

	pgInstance := db.NewPgInstance(dbPool)

	// query, safety concurrency
	itemRepo, err := pgInstance.BeginUpload(ctx)
	if err != nil {
		return nil, err
	}

	return &UploadItemProcessor{
		itemRepo: itemRepo,
		storage:  storage,
	}, nil

}
