package common

import (
	"bytes"
	"context"
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"io"
	"net/http"

	kithttp "github.com/go-kit/kit/transport/http"
)

func GetMD5Hash(text string) string {
	hash := md5.Sum([]byte(text))
	return hex.EncodeToString(hash[:])
}

// encodeResponse is the common method to encode all response types to the client.
func EncodeResponse(ctx context.Context, w http.ResponseWriter, response interface{}) error {
	return kithttp.EncodeJSONResponse(ctx, w, response)
}

func EncodeFileResponse(ctx context.Context, w http.ResponseWriter, response interface{}) error {
	w.Header().Set("Content-Type", "audio/wav")
	r := bytes.NewReader(response.([]byte))
	_, err := io.Copy(w, r)
	return err
}

func EncodeError(ctx context.Context, err error, w http.ResponseWriter) {
	// maybe we can be smart here by returning text/json error based on request's
	// content-type header
	EncodeJSONError(ctx, err, w)
}

func EncodeJSONError(_ context.Context, err error, w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	// we can have custom response headers by implementing kithttp.Headerer in
	// our response struct
	if headerer, ok := err.(kithttp.Headerer); ok {
		for k, values := range headerer.Headers() {
			for _, v := range values {
				w.Header().Add(k, v)
			}
		}
	}
	code := http.StatusInternalServerError
	// and custom status code
	if sc, ok := err.(kithttp.StatusCoder); ok {
		code = sc.StatusCode()
	}
	w.WriteHeader(code)
	// enforce json err response
	json.NewEncoder(w).Encode(map[string]interface{}{
		"error": err.Error(),
	})
}
