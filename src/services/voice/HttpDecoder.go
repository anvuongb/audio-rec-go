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

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		return nil, err
	}
	soundRate, err := strconv.Atoi(r.FormValue("sound_rate"))
	if err == nil {
		req.SoundRate = int32(soundRate)
	}
	return req, nil

}
