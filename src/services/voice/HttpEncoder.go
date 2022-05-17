package voice

import (
	"bytes"
	"context"
	"io"
	"net/http"
)

func EncodeFileAudioResponse(ctx context.Context, w http.ResponseWriter, response interface{}) error {
	w.Header().Set("Content-Type", "audio/wav")
	r := bytes.NewReader(response.([]byte))
	_, err := io.Copy(w, r)
	return err
}

func EncodeFileImageResponse(ctx context.Context, w http.ResponseWriter, response interface{}) error {
	w.Header().Set("Content-Type", "image/png")
	r := bytes.NewReader(response.([]byte))
	_, err := io.Copy(w, r)
	return err
}
