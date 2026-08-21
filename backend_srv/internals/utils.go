package internals

import (
	"bytes"
	"crypto/rand"
	"database/sql"
	"encoding/json"
	"net/http"
	"time"
)

func WriteErrorResponse(w http.ResponseWriter, statusHttp int, messageErr string) {

	var b bytes.Buffer
	if err := json.NewEncoder(&b).Encode(messageErr); err != nil {
		http.Error(w, "gagal encode info error", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusHttp)
	w.Write(b.Bytes())
}

func WriteResponse(w http.ResponseWriter, statusHttp int, data map[string]any) {

	if statusHttp != 200 && statusHttp != 202 {

		strData, ok := data["message"].(string)
		if !ok {
			WriteErrorResponse(w, http.StatusInternalServerError, "gagal parsing info error")
			return
		}
		WriteErrorResponse(w, statusHttp, strData)
		return
	}

	var b bytes.Buffer
	if err := json.NewEncoder(&b).Encode(data); err != nil {
		http.Error(w, "gagal encode info error", statusHttp)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusHttp)
	w.Write(b.Bytes())
}

type NullString struct {
	sql.NullString
}

func (ns NullString) MarshalJSON() ([]byte, error) {
	if !ns.Valid {
		return []byte("null"), nil
	}
	return json.Marshal(ns.String)
}

func (ns *NullString) UnmarshalJSON(data []byte) error {
	var s *string
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}
	if s == nil {
		ns.String, ns.Valid = "", false
		return nil
	}
	ns.String, ns.Valid = *s, true
	return nil
}

type NullTime struct {
	sql.NullTime
}

func (nt NullTime) MarshalJSON() ([]byte, error) {
	if !nt.Valid {
		return []byte("null"), nil
	}
	return json.Marshal(nt.Time)
}

func (nt *NullTime) UnmarshalJSON(data []byte) error {
	var t *time.Time
	if err := json.Unmarshal(data, &t); err != nil {
		return err
	}
	if t == nil {
		nt.Time, nt.Valid = time.Time{}, false
		return nil
	}
	nt.Time, nt.Valid = *t, true
	return nil
}

type NullBool struct{ sql.NullBool }

func (v NullBool) MarshalJSON() ([]byte, error) {
	if !v.Valid {
		return []byte("null"), nil
	}
	return json.Marshal(v.Bool)
}
func (v *NullBool) UnmarshalJSON(data []byte) error {
	var b *bool
	if err := json.Unmarshal(data, &b); err != nil {
		return err
	}
	if b == nil {
		v.Bool, v.Valid = false, false
		return nil
	}
	v.Bool, v.Valid = *b, true
	return nil
}

// ---- Int64 ----
type NullInt64 struct{ sql.NullInt64 }

func (v NullInt64) MarshalJSON() ([]byte, error) {
	if !v.Valid {
		return []byte("null"), nil
	}
	return json.Marshal(v.Int64)
}
func (v *NullInt64) UnmarshalJSON(data []byte) error {
	var n *int64
	if err := json.Unmarshal(data, &n); err != nil {
		return err
	}
	if n == nil {
		v.Int64, v.Valid = 0, false
		return nil
	}
	v.Int64, v.Valid = *n, true
	return nil
}

func RandomStringSuffix(n int) string {
	const chars = "ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	b := make([]byte, n)
	rand.Read(b)
	for i := range b {
		b[i] = chars[int(b[i])%len(chars)]
	}
	return string(b)
}
