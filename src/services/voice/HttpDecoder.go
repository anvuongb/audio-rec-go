package voice

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
)

func DecodeGenericRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var req GenericRequest
	if r.Method == http.MethodGet {
		req.RequestId = r.URL.Query().Get("request_id")
		return req, nil
	}
	if r.Method == http.MethodPost {
		err := json.NewDecoder(r.Body).Decode(&req)
		if err != nil {
			return nil, err
		}
		return req, nil
	}
	return nil, fmt.Errorf("unauthorized http method %s", r.Method)
}

func DecodeSaveAudioRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var req VoiceFile

	file, _, err := r.FormFile("file")
	if err == nil {
		req.File = file
	}

	req.FileId = r.FormValue("file_id")
	req.RequestId = r.FormValue("request_id")
	req.GeneratedText = r.FormValue("generated_text")

	soundRate, err := strconv.Atoi(r.FormValue("sound_rate"))
	if err == nil {
		req.SoundRate = int32(soundRate)
	}

	masked, err := strconv.Atoi(r.FormValue("masked"))
	if err == nil {
		if masked == 0 {
			req.Masked = false
		} else {
			req.Masked = true
		}
	}
	return req, nil

}
