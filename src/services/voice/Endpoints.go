package voice

import (
	"context"

	"github.com/go-kit/kit/endpoint"
)

type Endpoints struct {
	InitMetadata endpoint.Endpoint
	SaveAudio    endpoint.Endpoint
}

func makeInitMetadataEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(GenericRequest)
		response, err := s.InitMetadata(ctx, req)
		return response, err
	}
}

func makeSaveAudioEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(VoiceFile)
		response, err := s.SaveAudio(ctx, req)
		return response, err
	}
}

func MakeEndpoints(s Service) Endpoints {
	return Endpoints{
		InitMetadata: makeInitMetadataEndpoint(s),
		SaveAudio:    makeSaveAudioEndpoint(s),
	}
}
