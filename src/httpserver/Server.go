package httpserver

import (
	"audio-rec-go/src/common"
	"audio-rec-go/src/config"
	"audio-rec-go/src/services/voice"
	"context"
	"fmt"
	"net/http"

	"github.com/go-chi/cors"
	"github.com/go-kit/kit/log"

	httptransport "github.com/go-kit/kit/transport/http"
	"github.com/gorilla/mux"
)

func version(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, config.GlobalConfig.Version)
}

func NewHTTPServer(ctx context.Context, logger log.Logger, useCORS bool, voiceEndpoints voice.Endpoints) http.Handler {
	r := mux.NewRouter()
	if useCORS {
		// cors := cors.New(cors.Options{
		// 	AllowedOrigins:   []string{"*"},
		// 	AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		// 	AllowedHeaders:   []string{"Accept", "X-Auth-Token", "Content-Type", "X-CSRF-Token"},
		// 	AllowCredentials: true,
		// })
		// r.Use(cors.Handler)
		// cors := cors.AllowAll()
		cors := cors.New(cors.Options{
			AllowedOrigins: []string{"*"},
			AllowedMethods: []string{
				http.MethodHead,
				http.MethodGet,
				http.MethodPost,
				http.MethodPut,
				http.MethodPatch,
				http.MethodDelete,
				http.MethodOptions,
			},
			AllowedHeaders:   []string{"*"},
			AllowCredentials: false,
		})
		r.Use(cors.Handler)
	}

	options := []httptransport.ServerOption{
		httptransport.ServerErrorLogger(logger),
		httptransport.ServerErrorEncoder(common.EncodeError),
	}

	// r.Use(commonMiddleware)
	r.Methods("GET", "OPTIONS").Path("/api/version").HandlerFunc(version)

	r.Methods("GET", "OPTIONS").Path("/api/getStats").Handler(httptransport.NewServer(
		voiceEndpoints.GetStats,
		voice.DecodeGenericRequest,
		common.EncodeResponse,
		options...,
	))

	r.Methods("GET", "OPTIONS").Path("/api/initMetadata").Handler(httptransport.NewServer(
		voiceEndpoints.InitMetadata,
		voice.DecodeGenericRequest,
		common.EncodeResponse,
		options...,
	))

	r.Methods("POST", "OPTIONS").Path("/api/saveAudio").Handler(httptransport.NewServer(
		voiceEndpoints.SaveAudio,
		voice.DecodeSaveAudioRequest,
		common.EncodeResponse,
		options...,
	))

	return r
}
