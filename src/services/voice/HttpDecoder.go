package voice

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/gorilla/mux"
)

func DecodeGenericRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var req GenericRequest
	if r.Method == http.MethodGet {
		req.RequestId = r.URL.Query().Get("request_id")
		req.FileId = r.URL.Query().Get("file_id")
		req.Masked = strings.ToLower(r.URL.Query().Get("masked")) == "true"

		pageNumber, err := strconv.Atoi(r.URL.Query().Get("page_number"))
		if err != nil {
			pageNumber = 0
		}
		req.PageNumber = pageNumber

		recordsPerPage, err := strconv.Atoi(r.URL.Query().Get("records_per_page"))
		if err != nil {
			recordsPerPage = 0
		}
		req.RecordsPerPage = recordsPerPage

		limit, err := strconv.Atoi(r.URL.Query().Get("limit"))
		if err != nil {
			limit = 0
		}
		req.Limit = limit
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

func DecodeGenericRequestAudio(_ context.Context, r *http.Request) (interface{}, error) {
	var req GenericRequest
	vars := mux.Vars(r)
	fileName, ok := vars["file"]
	if !ok {
		return nil, fmt.Errorf("file not")
	}
	if r.Method == http.MethodGet {
		req.RequestId = r.URL.Query().Get("request_id")
		req.FileId = strings.Split(fileName, "_")[0]
		req.Masked = strings.Split(strings.Split(fileName, "_")[1], ".")[0] == "masked"

		pageNumber, err := strconv.Atoi(r.URL.Query().Get("page_number"))
		if err != nil {
			pageNumber = 0
		}
		req.PageNumber = pageNumber

		recordsPerPage, err := strconv.Atoi(r.URL.Query().Get("records_per_page"))
		if err != nil {
			recordsPerPage = 0
		}
		req.RecordsPerPage = recordsPerPage

		limit, err := strconv.Atoi(r.URL.Query().Get("limit"))
		if err != nil {
			limit = 0
		}
		req.Limit = limit
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

	if r.FormValue("gender") != "" && r.FormValue("gender") != "undefined" {
		req.Gender = r.FormValue("gender")
	} else {
		req.Gender = "N/A"
	}
	if r.FormValue("country") != "" && r.FormValue("country") != "undefined" {
		req.Country = r.FormValue("country")
	} else {
		req.Country = "N/A"
	}
	if r.FormValue("mask_type") != "" && r.FormValue("mask_type") != "undefined" {
		req.MaskType = r.FormValue("mask_type")
	} else {
		req.MaskType = "N/A"
	}
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
